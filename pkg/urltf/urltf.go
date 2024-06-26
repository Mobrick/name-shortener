package urltf

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

// MakeResultShortenedURL создает адрес для ответа пользователю по полученному сокращенному URL.
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

// EncodeURL сокращает URL с помощью рандома и хэша.
func EncodeURL(longURL []byte, shortURLLength int) (string, error) {
	var newURL string

	if shortURLLength <= 0 {
		return "", errors.New("expected length is not valid")
	}

	hash := make([]byte, shortURLLength)
	_, err := rand.Read(hash)
	if err != nil {
		return newURL, err
	}

	encodedHash := base64.URLEncoding.EncodeToString(hash)
	newURL = encodedHash[:shortURLLength]

	return newURL, nil
}
