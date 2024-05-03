package urltf

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_encodeURL(t *testing.T) {
	tests := []struct {
		name            string
		longURL         []byte
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
			newURL, err := EncodeURL(tt.longURL, tt.wantURLLength)
			if err != nil {
				log.Printf("encode error %v", err)
			}
			assert.Equal(t, tt.wantURLLength, len(newURL))
		})
	}
}

func BenchmarkEncodeURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		EncodeURL([]byte("https://music.yandex.ru/home"), 8)
	}
}
