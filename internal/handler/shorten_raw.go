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
	existingShortURL := db.Add(shortURL, string(urlToShorten))

	var resultAddress string
	var status int
	if len(existingShortURL) != 0 {
		resultAddress = existingShortURL
		status = http.StatusConflict
	} else {
		resultAddress = shortAddress
		status = http.StatusCreated
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(status)
	res.Write([]byte(resultAddress))
}
