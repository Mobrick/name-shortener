package handler

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mobrick/name-shortener/internal/auth"
	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnv_LongURLHandle(t *testing.T) {
	env := &Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "positive shorten test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			want: want{
				code:        201,
				contentType: "text/plain",
			},
		},
		{
			name:    "empty shorten test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("")),
			want: want{
				code:        400,
				contentType: "",
			},
		},
		{
			name:    "conflict shorten test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.go.com/")),
			want: want{
				code:        409,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := test.request
			w := httptest.NewRecorder()

			cookie, err := auth.CreateNewCookie(uuid.New().String())
			if err != nil {
				assert.Error(t, err, err.Error())
				return
			}
			request.AddCookie(&cookie)

			env.LongURLHandle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			log.Print("Result short address " + string(resBody))
			log.Print("Config base address " + env.ConfigStruct.FlagShortURLBaseAddr)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func BenchmarkLongURLHandle(b *testing.B) {
	env := &Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/"))
	w := httptest.NewRecorder()
	cookie, err := auth.CreateNewCookie(uuid.New().String())
	if err != nil {
		return
	}
	request.AddCookie(&cookie)
	for i := 0; i < b.N; i++ {
		env.LongURLHandle(w, request)
	}
}
