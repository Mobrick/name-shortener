package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mobrick/name-shortener/internal/userauth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUserIDFromRequest(t *testing.T) {
	type want struct {
		userID string
		ok     bool
	}
	tests := []struct {
		name         string
		request      *http.Request
		createCookie bool
		want         want
	}{
		{
			name:         "positive id test #1",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			createCookie: true,
			want: want{
				userID: "1a91a181-80ec-45cb-a576-14db11505700",
				ok:     true,
			},
		},
		{
			name:         "negative id test #1",
			request:      httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			createCookie: false,
			want: want{
				userID: "",
				ok:     false,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := test.request
			if test.createCookie {
				cookie, err := userauth.CreateNewCookie(test.want.userID)
				if err != nil {
					assert.Error(t, err, err.Error())
					return
				}
				request.AddCookie(&cookie)
				require.NoError(t, err)
			}
			resuldID, resultOK := GetUserIDFromRequest(request)

			assert.Equal(t, test.want.userID, resuldID)
			assert.Equal(t, test.want.ok, resultOK)
		})
	}
}
