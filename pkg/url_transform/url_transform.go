package url_transform

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"slices"
	"strings"

	"github.com/Mobrick/name-shortener/utils"
)

func MakeShortUrl(shortAddress string, dbMap map[string]string, urlToShorten []byte, req *http.Request) string {
	if !strings.HasSuffix(shortAddress, "/") {
		shortAddress += "/"
	}

	keys := utils.GetKeys(dbMap)

	shortURL := encodeURL(urlToShorten, keys)
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

func encodeURL(longURL []byte, keys []string) string {
	var newURL string

	for {
		hash := make([]byte, 8)
		_, err := rand.Read(hash)
		if err != nil {
			panic(err)
		}

		encodedHash := base64.URLEncoding.EncodeToString(hash)
		newURL = encodedHash[:8]

		if !slices.Contains(keys, newURL) {
			break
		}
	}
	return newURL
}
