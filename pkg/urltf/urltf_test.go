package urltf

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeURL(t *testing.T) {
	tests := []struct {
		name            string
		longURL         []byte
		resultURLLength int
		wantURLLength   int
		wantHasError    bool
	}{
		{
			name:            "positive encode URL test #1",
			longURL:         []byte("https://music.yandex.ru/home"),
			resultURLLength: 8,
			wantURLLength:   8,
			wantHasError:    false,
		},
		{
			name:            "positive encode URL test #2",
			longURL:         []byte("https://music.yandex.ru/home"),
			resultURLLength: 16,
			wantURLLength:   16,
			wantHasError:    false,
		},
		{
			name:            "positive encode URL test #3",
			longURL:         []byte("https://music.yandex.ru/home"),
			resultURLLength: 100,
			wantURLLength:   100,
			wantHasError:    false,
		},
		{
			name:            "positive encode URL test #4",
			longURL:         nil,
			resultURLLength: 8,
			wantURLLength:   8,
			wantHasError:    false,
		},
		{
			name:            "negative encode URL test #1",
			longURL:         nil,
			resultURLLength: -8,
			wantURLLength:   0,
			wantHasError:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newURL, err := EncodeURL(tt.longURL, tt.resultURLLength)
			if err != nil {
				log.Printf("encode error %v", err)
			}
			assert.Equal(t, tt.wantHasError, err != nil)
			assert.Equal(t, tt.wantURLLength, len(newURL))
		})
	}
}

func BenchmarkEncodeURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeURL([]byte("https://music.yandex.ru/home"), 8)
	}
}

func TestMakeResultShortenedURL(t *testing.T) {
	tests := []struct {
		name         string
		shortAddress string
		shortURL     string
		request      *http.Request
		wantURL      string
	}{
		{
			name:         "positive make result test #1",
			shortAddress: "http://localhost:8080/",
			shortURL:     "xxxxdddd",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			wantURL:      "http://localhost:8080/xxxxdddd",
		},
		{
			name:         "positive make result test #2",
			shortAddress: "http://shortener/",
			shortURL:     "xxxxDDDD",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.go.com/")),
			wantURL:      "http://shortener/xxxxDDDD",
		},
		{
			name:         "positive make result test #3",
			shortAddress: "",
			shortURL:     "xxxxDDDD",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.go.com/")),
			wantURL:      "/xxxxDDDD",
		},
		{
			name:         "positive make result test #4",
			shortAddress: "http://shortener/",
			shortURL:     "xxxxDDDD",
			request:      nil,
			wantURL:      "http://shortener/xxxxDDDD",
		},
		{
			name:         "positive make result test #5",
			shortAddress: "",
			shortURL:     "",
			request:      nil,
			wantURL:      "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := test.request
			resultURL := MakeResultShortenedURL(test.shortAddress, test.shortURL, req)
			assert.Equal(t, test.wantURL, resultURL)
		})
	}
}
