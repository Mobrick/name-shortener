package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerEnv_BatchHandler(t *testing.T) {
	env := &HandlerEnv{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()
	type want struct {
		code        int
		responseLen int
		contentType string
	}
	tests := []struct {
		name string
		body []models.BatchRequestURL
		want want
	}{
		{
			name: "positive shorten test #1",
			body: []models.BatchRequestURL{{
				CorrelationID: "1234",
				OriginalURL:   "https://www.google.com/",
			}, {
				CorrelationID: "1235",
				OriginalURL:   "https://www.go.com/",
			}},
			want: want{
				code:        201,
				responseLen: len(env.ConfigStruct.FlagShortURLBaseAddr) + ShortURLLength,
				contentType: "text/plain",
			},
		},
		{
			name: "empty shorten test #1",
			body: []models.BatchRequestURL{},
			want: want{
				code:        400,
				responseLen: 0,
				contentType: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.body)
			if err != nil {
				assert.Error(t, err, err.Error())
			}

			request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
			w := httptest.NewRecorder()
			env.BatchHandler(w, request)

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
