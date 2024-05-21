package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewCookie(t *testing.T) {
	tests := []struct {
		name           string
		newID          string
		wantCookieName string
		wantHasErr     bool
	}{
		{
			name:           "positive new cookie test #1",
			newID:          uuid.New().String(),
			wantCookieName: "auth_token",
			wantHasErr:     false,
		},
		{
			name:           "negative new cookie test #1",
			newID:          "",
			wantCookieName: "",
			wantHasErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie, err := CreateNewCookie(tt.newID)
			assert.Equal(t, tt.wantHasErr, err != nil)
			assert.Equal(t, tt.wantCookieName, cookie.Name)
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name         string
		newID        string
		scuffedToken bool
		wantID       string
		wantIsValid  bool
	}{
		{
			name:         "positive get user id test #1",
			newID:        "1a91a181-80ec-45cb-a576-14db11505700",
			scuffedToken: false,
			wantID:       "1a91a181-80ec-45cb-a576-14db11505700",
			wantIsValid:  true,
		},
		{
			name:         "negative get user id test #1",
			newID:        "1a91a181-80ec-45cb-a576-14db11505700",
			scuffedToken: true,
			wantID:       "",
			wantIsValid:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie, err := CreateNewCookie(tt.newID)
			require.NoError(t, err)
			if tt.scuffedToken {
				cookie.Value = cookie.Value[:len(cookie.Value)-1]
			}
			id, ok := GetUserID(cookie.Value)
			assert.Equal(t, tt.wantIsValid, ok)
			assert.Equal(t, tt.wantID, id)
		})
	}
}

func Test_buildJWTString(t *testing.T) {
	tests := []struct {
		name         string
		newID        string
		wantNotEmpty bool
		wantErr      bool
	}{
		{
			name:         "positive build JWT test #1",
			newID:        uuid.New().String(),
			wantNotEmpty: true,
			wantErr:      false,
		},
		{
			name:         "positive build JWT test #2",
			newID:        "",
			wantNotEmpty: true,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildJWTString(tt.newID)
			assert.Equal(t, tt.wantNotEmpty, len(got) > 0)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func Test_cookieIsValid(t *testing.T) {
	tests := []struct {
		name         string
		r            *http.Request
		scuffedToken bool
		want         bool
	}{
		{
			name:         "pisitive valid cookie test #1",
			r:            httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			scuffedToken: false,
			want:         true,
		},
		{
			name:         "pisitive valid cookie test #2",
			r:            httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.google.com/")),
			scuffedToken: true,
			want:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie, err := CreateNewCookie(uuid.New().String())
			require.NoError(t, err)
			if tt.scuffedToken {
				cookie.Value = cookie.Value[:len(cookie.Value)-1]
			}
			tt.r.AddCookie(&cookie)
			got := cookieIsValid(tt.r)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCookieMiddleware(t *testing.T) {
	type args struct {
		h http.Handler
	}
	tests := []struct {
		name            string
		args            args
		wantCookieName  string
		wantCookieCount int
	}{
		{
			name: "positive cookie muddleware test #1",
			args: args{
				h: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
			wantCookieName:  "auth_token",
			wantCookieCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookieMiddleware := CookieMiddleware(tt.args.h)
			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()
			cookieMiddleware.ServeHTTP(rr, req)
			result := rr.Result()
			cookies := result.Cookies()
			defer result.Body.Close()
			assert.Equal(t, tt.wantCookieCount, len(cookies))
			assert.Equal(t, tt.wantCookieName, cookies[0].Name)
		})
	}
}
