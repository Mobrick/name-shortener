package userauth

import (
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
