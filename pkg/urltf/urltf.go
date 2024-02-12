package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Mobrick/name-shortener/database"
)

func MakeShortAddressAndURL(shortAddress string, db database.DatabaseData, urlToShorten []byte, req *http.Request, shortURLLength int) (string, string) {
	if !strings.HasSuffix(shortAddress, "/") {
		shortAddress += "/"
	}

	shortURL := encodeURL(urlToShorten, db, shortURLLength)

	if len(shortAddress) != 0 {
		shortAddress += shortURL
	} else {
		shortAddress = req.Host + req.URL.Path + shortURL
		if len(req.URL.Scheme) == 0 {
			shortAddress = "http://" + shortAddress
		}
	}
	return shortAddress, shortURL
}

func encodeURL(longURL []byte, db database.DatabaseData, shortURLLength int) string {
	var newURL string

	for {
		hash := make([]byte, shortURLLength)
		_, err := rand.Read(hash)
		if err != nil {
			panic(err)
		}

		encodedHash := base64.URLEncoding.EncodeToString(hash)
		newURL = encodedHash[:shortURLLength]

		// TODO: сделать проверку есть ли такой адрес уже в бд а не в мапе, если работа идет с бд
		if !db.Contains(newURL) {
			break
		}
	}
	return newURL
}
