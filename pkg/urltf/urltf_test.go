package urltf

import (
	"log"
	"testing"

	"github.com/Mobrick/name-shortener/database"
	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/stretchr/testify/assert"
)

func Test_encodeURL(t *testing.T) {
	file, err := filestorage.MakeFile("tmp/test.json")	
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	db := database.NewDBFromFile(file)
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
			newURL := encodeURL(tt.longURL, db, tt.wantURLLength)
			assert.Equal(t, tt.wantURLLength, len(newURL))
		})
	}
}
