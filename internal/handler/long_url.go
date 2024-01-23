package handler

import (
	"io"
	"net/http"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/url_transform"
)

type HandlerEnv config.Env

func (env HandlerEnv) LongURLHandle(res http.ResponseWriter, req *http.Request) {
	urlToShorten, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	if len(urlToShorten) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	shortAddress := url_transform.MakeShortUrl(env.ConfigStruct.FlagShortURLBaseAddr, env.DatabaseMap, urlToShorten, req, config.ShortURLLength)

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortAddress))
}
