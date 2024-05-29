package logger

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoggingMiddleware(t *testing.T) {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	Sugar = *zapLogger.Sugar()

	tests := []struct {
		name       string
		status     int
		wantStatus int
		body       string
		wantBody   string
	}{
		{
			name:       "positive logging middleware #1",
			status:     http.StatusOK,
			body:       "OK",
			wantStatus: http.StatusOK,
			wantBody:   "OK",
		},
		{
			name:       "positive logging middleware #2",
			status:     http.StatusAccepted,
			body:       "ook",
			wantStatus: http.StatusAccepted,
			wantBody:   "ook",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
				w.Write([]byte(tt.body))
			})

			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			LoggingMiddleware(handler).ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
