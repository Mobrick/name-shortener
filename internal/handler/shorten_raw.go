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

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr
	encodedURL := urltf.EncodeURL(urlToShorten, db, ShortURLLength)

	existingShortURL := db.Add(encodedURL, string(urlToShorten))

	var resultAddress string
	var status int
	if len(existingShortURL) != 0 {
		resultAddress = urltf.MakeResultShortenedURL(hostAndPathPart, existingShortURL, req)
		status = http.StatusConflict
	} else {
		resultAddress = urltf.MakeResultShortenedURL(hostAndPathPart, encodedURL, req)
		status = http.StatusCreated
	}

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(status)
	res.Write([]byte(resultAddress))
}
