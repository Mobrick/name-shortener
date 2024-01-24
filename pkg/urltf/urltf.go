package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

func MakeShortUrl(shortAddress string, dbMap map[string]string, urlToShorten []byte, req *http.Request) string {
	if !strings.HasSuffix(shortAddress, "/") {
		shortAddress += "/"
	}


	shortURL := encodeURL(urlToShorten, dbMap)
	dbMap[shortURL] = string(urlToShorten)

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

func encodeURL(longURL []byte, dbMap map[string]string) string {
	var newURL string

	for {
		hash := make([]byte, 8)
		_, err := rand.Read(hash)
		if err != nil {
			panic(err)
		}

		encodedHash := base64.URLEncoding.EncodeToString(hash)
		newURL = encodedHash[:8]

		if _, ok := dbMap[string(newURL)]; !ok {
			break
		}
	}
	return newURL
}
