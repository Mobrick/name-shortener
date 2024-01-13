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
			dbMap = test.db
			request := httptest.NewRequest(http.MethodGet, "/"+test.request, nil)
			w := httptest.NewRecorder()
			urlHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resLocation := res.Header.Get("Location")

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
		request *http.Request
		want want
	}{
		{
			name: "positive POST test #1",
			request: httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			want: want{
				code:        201,
				responseLen: len("http://example.com/") + shortURLLength,
				contentType: "text/plain",
			},
		},		
		{
			name: "empty POST test #1",
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
