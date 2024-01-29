package compression

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Mobrick/name-shortener/logger"
)

var acceptableContentTypes = []string{"application/json", "text/html"}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

type gzipReader struct {
	r     io.ReadCloser
	zr *gzip.Reader
}

func newGzipReader(r io.ReadCloser) (*gzipReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &gzipReader{
		r:     r,
		zr: zr,
	}, nil
}

func (c gzipReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *gzipReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
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

		var acceptContentType bool

		// проверка допустимости типа контента
		headerContentType := r.Header.Get("Content-Type")

		for _, contentType := range acceptableContentTypes {
			if headerContentType == contentType {
				acceptContentType = true
				break
			}
		}

		if !acceptContentType {
			next.ServeHTTP(w, r)
			return
		}

		// проверка поддержки gzip
		encodeHeaderValues := r.Header.Values("Accept-Encoding")

		var acceptGzip bool
	out:
		for _, value := range encodeHeaderValues {
			encodings := strings.Split(value, ",")

			for _, encoding := range encodings {
				if strings.TrimSpace(encoding) == "gzip" {
					acceptGzip = true
					break out
				}
			}
		}
		logger.Sugar.Infoln("Header values has gzip", acceptGzip)

		if !acceptGzip {
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer zw.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: zw}, r)
	})
}
