package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Mobrick/name-shortener/database"
)

func MakeResultShortenedURL(shortAddress string, shortURL string, req *http.Request) string {
	if !strings.HasSuffix(shortAddress, "/") {
		shortAddress += "/"
	}

	if len(shortAddress) != 0 {
		shortAddress += shortURL
	} else {
		shortAddress = req.Host + req.URL.Path + shortURL
		if len(req.URL.Scheme) == 0 {
			shortAddress = "http://" + shortAddress
		}
	}
	return shortAddress
}

func EncodeURL(longURL []byte, db database.DatabaseData, shortURLLength int) string {
	var newURL string

	for {
		hash := make([]byte, shortURLLength)
		_, err := rand.Read(hash)
		if err != nil {
			panic(err)
		}

		encodedHash := base64.URLEncoding.EncodeToString(hash)
		newURL = encodedHash[:shortURLLength]

		if !db.Contains(newURL) {
			break
		}
	}
	return newURL
}
