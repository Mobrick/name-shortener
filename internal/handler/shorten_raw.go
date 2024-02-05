package handler

import (
	"io"
	"net/http"

	"github.com/Mobrick/name-shortener/urltf"
)

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
	db := env.DatabaseData
	shortAddress, shortURL := urltf.MakeShortAddressAndURL(env.ConfigStruct.FlagShortURLBaseAddr, db, urlToShorten, req, ShortURLLength)
	db.Add(shortURL, string(urlToShorten))

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(shortAddress))
}
