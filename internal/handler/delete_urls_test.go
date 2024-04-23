package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerEnv_DeleteUserUsrlsHandler(t *testing.T) {
	env := &HandlerEnv{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()
	type want struct {
		code        int
		contentType string
	}
	tests := []struct {
		name string
		body []string
		want want
	}{
		{
			name: "positive deletion test #1",
			body: []string{
				"6qxTVvsy", "RTfd56hn", "Jlfd67ds",
			},
			want: want{
				code: 202,
			},
		},
		{
			name: "empty deletion test #1",
			body: []string{},
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.body)
			if err != nil {
				assert.Error(t, err, err.Error())
			}

			request := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(body))
			w := httptest.NewRecorder()
			env.DeleteUserUsrlsHandler(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()

			require.NoError(t, err)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
