package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ShortenedURLHandle перенаправляет по адресу сохраненному в хранилищу в соотвествии с полученным сокращенным адресом.
func (env Env) ShortenedURLHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	shortURL := chi.URLParam((req), "shortURL")
	if len(shortURL) != ShortURLLength {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	location, ok, isDeleted, err := env.Storage.Get(ctx, string(shortURL))
	if err != nil {
		log.Printf("could not complete original address request")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Printf("no matching data to %s found", shortURL)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if isDeleted {
		log.Printf("this URL is deleted: %s", shortURL)
		res.WriteHeader(http.StatusGone)
		return
	}
	res.Header().Set("Location", location)
	http.Redirect(res, req, location, http.StatusTemporaryRedirect)
}
