package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi"
)

var dbMap map[string]string

func longURLHandle(res http.ResponseWriter, req *http.Request) {
	urlToShorten, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	if len(urlToShorten) == 0 {			
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := req.Host + req.URL.Path + makeShortURL(urlToShorten)
	if len(req.URL.Scheme) == 0 {
		shortURL = "http://" + shortURL
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
}

func shortenedURLHandle(res http.ResponseWriter, req *http.Request) {
	shortURL := chi.URLParam(req, "shortURL")
	if len(shortURL) != 8 {
		res.WriteHeader(http.StatusBadRequest)
		return			
	}
	location, ok := dbMap[string(shortURL)]
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", location)
	http.Redirect(res, req, location, http.StatusTemporaryRedirect)
}

func main() {
	dbMap = make(map[string]string)
	r := chi.NewRouter()

	r.Post(`/`, longURLHandle)
	r.Get(`/{shortURL}`, shortenedURLHandle)

	log.Fatal(http.ListenAndServe(`:8080`, r))
}

func makeShortURL(longURL []byte) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortURLLength := 8
	var newURL []byte
	for {
		newURL = make([]byte, shortURLLength)
		for i := 0; i < shortURLLength; i++ {
			newURL[i] = letters[rand.Intn(len(letters))]
		}
		if _, ok := dbMap[string(newURL)]; !ok {
			break
		}
	}
	dbMap[string(newURL)] = string(longURL)
	return string(newURL)
}
