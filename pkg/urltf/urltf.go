package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

func MakeShortAddressAndURL(shortAddress string, dbMap map[string]string, urlToShorten []byte, req *http.Request, shortURLLength int) (string, string) {
	if !strings.HasSuffix(shortAddress, "/") {
		shortAddress += "/"
	}

	shortURL := encodeURL(urlToShorten, dbMap, shortURLLength)

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

func encodeURL(longURL []byte, dbMap map[string]string, shortURLLength int) string {
	var newURL string

	for {
		hash := make([]byte, shortURLLength)
		_, err := rand.Read(hash)
		if err != nil {
			panic(err)
		}

		encodedHash := base64.URLEncoding.EncodeToString(hash)
		newURL = encodedHash[:shortURLLength]

		if _, ok := dbMap[string(newURL)]; !ok {
			break
		}
	}
	return newURL
}
