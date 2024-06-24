package compression

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecompressMiddleware(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive decompress middleware #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
			})

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			DecompressMiddleware(handler).ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "OK", w.Body.String())
		})
	}
}

func Test_gzipReader_Close(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive gzip reader close test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var b bytes.Buffer
			w := gzip.NewWriter(&b)
			w.Write([]byte("https://www.google.com/"))
			request := httptest.NewRequest(http.MethodPost, "/", &b)
			w.Close()
			request.Header.Set("Content-Encoding", "gzip")
			newReader, err := newGzipReader(request.Body)
			require.NoError(t, err)
			require.NoError(t, newReader.Close())
		})
	}
}

func Test_gzipReader_Read(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive read gzip test #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				var b bytes.Buffer
				w := gzip.NewWriter(&b)
				w.Write([]byte("https://www.google.com/"))
				request := httptest.NewRequest(http.MethodPost, "/", &b)
				w.Close()
				request.Header.Set("Content-Encoding", "gzip")
				newReader, err := newGzipReader(request.Body)
				require.NoError(t, err)
				urlToShorten, err := io.ReadAll(io.Reader(request.Body))
				defer request.Body.Close()
				require.NoError(t, err)
				_, err = newReader.Read(urlToShorten)
				require.NoError(t, err)
				require.NoError(t, newReader.Close())
			})
		})
	}
}
