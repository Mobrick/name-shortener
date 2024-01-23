package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/config"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHandlerEnv_ShortenedURLHandle(t *testing.T) {
	env := &HandlerEnv{
		DatabaseMap: config.NewDBMap(),
	}
	type want struct {
		code     int
		location string
	}
	tests := []struct {
		name    string
		request string
		db      map[string]string
		want    want
	}{
		{
			name:    "positive GET test #1",
			request: "DDAAssaa",
			db: map[string]string{
				"DDAAssaa": "http://example.com/",
			},
			want: want{
				code:     307,
				location: "http://example.com/",
			},
		},
		{
			name:    "9 letters request GET test #1",
			request: "DDAAssaaD",
			db: map[string]string{
				"DDAAssaa": "http://example.com/",
			},
			want: want{
				code:     400,
				location: "",
			},
		},
		{
			name:    "empty request GET test #1",
			request: "",
			db: map[string]string{
				"DDAAssaa": "http://example.com/",
			},
			want: want{
				code:     400,
				location: "",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			env.DatabaseMap = test.db
			request := httptest.NewRequest(http.MethodGet, "/{shortURL}", nil)
			requestContext := chi.NewRouteContext()
			requestContext.URLParams.Add("shortURL", test.request)

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestContext))

			w := httptest.NewRecorder()

			env.ShortenedURLHandle(w, request)

			res := w.Result()
			defer res.Body.Close()
			resLocation := res.Header.Get("Location")

			assert.Equal(t, test.want.location, resLocation)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}
