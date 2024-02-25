package handler

import (
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
	location, ok := env.Storage.Get(ctx, string(shortURL))
	if !ok {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	res.Header().Set("Location", location)
	http.Redirect(res, req, location, http.StatusTemporaryRedirect)
}
