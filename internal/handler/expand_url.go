package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (env HandlerEnv) ShortenedURLHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	shortURL := chi.URLParam((req), "shortURL")
	if len(shortURL) != ShortURLLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	location, ok, err := env.Storage.Get(ctx, string(shortURL))
	if err != nil {
		log.Printf("could not copmplete original address request")		
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Printf("no matching data to %s found", shortURL)	
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", location)
	http.Redirect(res, req, location, http.StatusTemporaryRedirect)
}
