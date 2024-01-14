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
	if req.Method == http.MethodGet {
		shortURL := req.URL.Path
		if len(shortURL) != 9 {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		location, ok := dbMap[string(shortURL[1:])]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Header().Set("Location", location)
		http.Redirect(res, req, dbMap[string(shortURL[1:])], http.StatusTemporaryRedirect)
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
