package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/internal/userauth"
	"github.com/google/uuid"
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
				contentType: "application/json",
			},
		},
		{
			name: "empty shorten test #1",
			body: []models.BatchRequestURL{},
			want: want{
				code:        400,
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

			cookie, err := userauth.CreateNewCookie(uuid.New().String())
			if err != nil {
				assert.Error(t, err, err.Error())
				return
			}
			request.AddCookie(&cookie)

			env.BatchHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func BenchmarkBatchHandler (b *testing.B) {
	env := &HandlerEnv{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}

	bodySlice := []models.BatchRequestURL{{
		CorrelationID: "1234",
		OriginalURL:   "https://www.google.com/",
	}, {
		CorrelationID: "1235",
		OriginalURL:   "https://www.go.com/",
	}}

	body, err := json.Marshal(bodySlice)
	if err != nil {
		log.Fatal(err.Error())
	}

	request := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewReader(body))
	w := httptest.NewRecorder()
	cookie, err := userauth.CreateNewCookie(uuid.New().String())
	if err != nil {		
		return
	}
	request.AddCookie(&cookie)
	for i := 0; i < b.N; i++ {
		env.BatchHandler(w, request)
	}
}