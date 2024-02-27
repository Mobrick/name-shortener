package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
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

func EncodeURL(longURL []byte, shortURLLength int) (string, error) {
	var newURL string

	hash := make([]byte, shortURLLength)
	_, err := rand.Read(hash)
	if err != nil {
		return newURL, err
	}

	encodedHash := base64.URLEncoding.EncodeToString(hash)
	newURL = encodedHash[:shortURLLength]
	
	return newURL, nil
}
