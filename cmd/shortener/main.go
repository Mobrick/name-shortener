package main

import (
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/database"
	"github.com/Mobrick/name-shortener/handler"
	"github.com/go-chi/chi/v5"
)

func main() {
	env := &handler.HandlerEnv{
		DatabaseMap:  database.NewDBMap(),
		ConfigStruct: config.MakeConfig(),
	}
	r := chi.NewRouter()

	r.Post(`/`, env.LongURLHandle)
	r.Get(`/{shortURL}`, env.ShortenedURLHandle)

	log.Fatal(http.ListenAndServe(env.ConfigStruct.FlagRunAddr, r))
}
