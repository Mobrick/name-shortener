package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mobrick/name-shortener/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLongURLHandle(t *testing.T) {
	env := &HandlerEnv{
		DatabaseMap:  config.NewDBMap(),
		ConfigStruct: config.MakeConfig(),
	}
	shortURLLength := config.ShortURLLength
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	type want struct {
		code        int
		responseLen int
		contentType string
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "positive POST test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			want: want{
				code:        201,
				responseLen: len(env.ConfigStruct.FlagShortURLBaseAddr) + shortURLLength,
				contentType: "text/plain",
			},
		},
		{
			name:    "empty POST test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")),
			want: want{
				code:        400,
				responseLen: 0,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := test.request
			w := httptest.NewRecorder()
			env.LongURLHandle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, test.want.responseLen, len(string(resBody)))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
