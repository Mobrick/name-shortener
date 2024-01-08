package main

import (
	"io"
	"math/rand"
	"net/http"
)

var dbMap map[string]string

func urlHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		urlToShorten, err := io.ReadAll(io.Reader(req.Body))
		if err != nil {
			res.Write([]byte(err.Error()))
			return
		}
		shortUrl := req.Host + req.URL.Path + makeShortUrl(urlToShorten)

		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(shortUrl))
	}
	if req.Method == http.MethodGet {
		shortUrl := req.URL.Path
		location := dbMap[string(shortUrl[1:])]
		res.Header().Set("Location", location)
		http.Redirect(res, req, dbMap[string(shortUrl[1:])], http.StatusTemporaryRedirect)
	}
}

func main() {
	dbMap = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, urlHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

func makeShortUrl(longUrl []byte) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortUrlLength := 8
	var newUrl []byte
	for {
		newUrl = make([]byte, shortUrlLength)
		for i := 0; i < shortUrlLength; i++ {
			newUrl[i] = letters[rand.Intn(len(letters))]
		}
		if _, ok := dbMap[string(newUrl)]; !ok {
			break
		}
	}
	dbMap[string(newUrl)] = string(longUrl)
	return string(newUrl)
}
