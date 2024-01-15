package main

import (
	"flag"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/go-chi/chi/v5"
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

	shortURL := flagShortURLBaseAddr

	// TODO возможно эти условия понадобятся только для тестирования, так что переделать тесты на использование args
	if len(flagShortURLBaseAddr) != 0 {
		shortURL += makeShortURL(urlToShorten)
	} else {
		shortURL = req.Host + req.URL.Path + makeShortURL(urlToShorten)
		if len(req.URL.Scheme) == 0 {
			shortURL = "http://" + shortURL
		}
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortURL))
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

func shortenedURLHandle(res http.ResponseWriter, req *http.Request) {
	shortURL := "/" + chi.URLParam((req), "shortURL")
	log.Print("Короткий урл " + shortURL)
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

// TODO перенести в другое место эту штуку

var flagRunAddr string
var flagShortURLBaseAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address to run server")
	flag.StringVar(&flagShortURLBaseAddr, "b", "http://localhost:8080/", "base address of shortened URL")

	flag.Parse()
}

func main() {
	parseFlags()
	dbMap = make(map[string]string)
	r := chi.NewRouter()

	r.Post(`/`, longURLHandle)
	r.Get(`/{shortURL}`, shortenedURLHandle)

	log.Fatal(http.ListenAndServe(flagRunAddr, r))
}
