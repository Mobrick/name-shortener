package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/userauth"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestEnv_ShortenedURLHandle(t *testing.T) {

	env := &Env{
		Storage: database.NewDB("tmp/short-url-db-test.json", ""),
	}
	defer env.Storage.Close()
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
			name:    "positive expand test #1",
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
			name:    "9 letters request expand test #1",
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
			name:    "empty request expand test #1",
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
			for k, v := range test.db {
				env.Storage.Add(context.Background(), k, v, "")
			}
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

func BenchmarkShortenedURLHandle(b *testing.B) {
	env := &Env{
		Storage:      database.NewDB("tmp/short-url-db-test.json", ""),
		ConfigStruct: config.MakeConfig(),
	}
	db := map[string]string{
		"DDAAssaa": "http://example.com/",
	}
	for k, v := range db {
		env.Storage.Add(context.Background(), k, v, "")
	}
	request := httptest.NewRequest(http.MethodGet, "/{shortURL}", nil)
	requestContext := chi.NewRouteContext()
	requestContext.URLParams.Add("shortURL", "DDAAssaa")

	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, requestContext))

	w := httptest.NewRecorder()
	cookie, err := userauth.CreateNewCookie("1a91a181-80ec-45cb-a576-14db11505700")
	if err != nil {
		log.Fatal(err.Error())
	}
	request.AddCookie(&cookie)
	for i := 0; i < b.N; i++ {
		env.UserUrlsHandler(w, request)
	}
}
