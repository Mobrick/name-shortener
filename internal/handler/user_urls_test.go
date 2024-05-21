package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/internal/auth"
	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandlerEnv_UserUrlsHandler(t *testing.T) {
	env := &Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	type want struct {
		code  int
		count int
	}
	tests := []struct {
		name string
		id   string
		want want
	}{
		{
			name: "2 urls ok test #1",
			id:   "1a91a181-80ec-45cb-a576-14db11505700",
			want: want{
				code:  200,
				count: 2,
			},
		}, {
			name: "0 url ok test #2",
			id:   "1954c654-dee9-44c7-81d1-6da6cfe918b2",
			want: want{
				code:  409,
				count: 0,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
			requestContext := chi.NewRouteContext()

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestContext))

			w := httptest.NewRecorder()

			cookie, err := auth.CreateNewCookie(test.id)
			if err != nil {
				assert.Error(t, err, err.Error())
				return
			}
			request.AddCookie(&cookie)

			env.UserUrlsHandler(w, request)

			res := w.Result()
			defer res.Body.Close()

			var urls []model.SimpleURLRecord

			var buf bytes.Buffer
			_, err = buf.ReadFrom(res.Body)
			if err != nil {
				assert.Error(t, err, err.Error())
				return
			}
			if err = json.Unmarshal(buf.Bytes(), &urls); err != nil {
				assert.Error(t, err, err.Error())
				return
			}

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.Equal(t, test.want.count, len(urls))
		})
	}
}

func BenchmarkUserUrlsHandler(b *testing.B) {
	env := &Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	request := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
	w := httptest.NewRecorder()
	cookie, err := auth.CreateNewCookie("1a91a181-80ec-45cb-a576-14db11505700")
	if err != nil {
		return
	}
	request.AddCookie(&cookie)
	for i := 0; i < b.N; i++ {
		env.UserUrlsHandler(w, request)
	}
}
