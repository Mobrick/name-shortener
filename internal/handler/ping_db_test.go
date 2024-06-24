package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestEnv_PingDBHandle(t *testing.T) {
	env := &Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()
	type want struct {
		code int
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative ping test #1",
			want: want{
				code: 500,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, "/api/shorten", nil)
			w := httptest.NewRecorder()
			env.PingDBHandle(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
