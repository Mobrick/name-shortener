package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetURLHandler(t *testing.T) {
	dbMap = make(map[string]string)
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/"+test.request, nil)
			w := httptest.NewRecorder()
			urlHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resLocation, err := res.Location()

			require.NoError(t, err)
			assert.Equal(t, test.want.location, resLocation)
			assert.Equal(t, test.want.code, res.StatusCode)
		})
	}
}

func TestPostURLHandler(t *testing.T) {
	dbMap = make(map[string]string)
	shortURLLength := 8
	type want struct {
		code        int
		responseLen int
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive POST test #1",
			want: want{
				code:        201,
				responseLen: len("http://example.com/") + shortURLLength,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/"))
			w := httptest.NewRecorder()
			urlHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			log.Print("response body: " + string(resBody))
			assert.Equal(t, test.want.responseLen, len(string(resBody)))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
