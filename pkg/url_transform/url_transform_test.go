package url_transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeURL(t *testing.T) {
	tests := []struct {
		name            string
		longURL         []byte
		keys            []string
		resultURLLength int
		wantURLLength   int
	}{
		{
			name:          "positive encode URL test #1",
			longURL:       []byte("https://music.yandex.ru/home"),
			wantURLLength: 8,
		},
		{
			name:          "positive encode URL test #2",
			longURL:       []byte("https://music.yandex.ru/home"),
			wantURLLength: 16,
		},
		{
			name:          "positive encode URL test #1",
			longURL:       []byte("https://music.yandex.ru/home"),
			wantURLLength: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newURL := encodeURL(tt.longURL, tt.keys, tt.wantURLLength)
			assert.Equal(t, tt.wantURLLength, len(newURL))
		})
	}
}
