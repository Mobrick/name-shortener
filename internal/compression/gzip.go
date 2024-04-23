package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read считывает для дальнейшего распаковывания слайса байт.
func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close закрывает ридер декомпрессора.
func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// DecompressMiddleware распактовывает данные из тела запроса
func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// проверка сжаты ли данные в запросе клиента
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			gzipReader, err := newGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = gzipReader
			defer gzipReader.Close()
		}

		next.ServeHTTP(w, r)
	})
}
