package handler

import (
	"io"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/logger"
	"github.com/Mobrick/name-shortener/pkg/urltf"
	"go.uber.org/zap"
)

// LongURLHandle возвращает сокращенный адрес.
func (env Env) LongURLHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, _ := GetUserIDFromRequest(req)

	urlToShorten, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		res.Write([]byte(err.Error()))
		return
	}
	if len(urlToShorten) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	storage := env.Storage

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr
	encodedURL, err := urltf.EncodeURL(urlToShorten, ShortURLLength)
	if err != nil {
		logger.Log.Debug("could not copmplete url encoding", zap.String("URL to encode", string(urlToShorten)))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	existingShortURL, err := storage.Add(ctx, encodedURL, string(urlToShorten), userID)
	if err != nil {
		logger.Log.Debug("could not copmplete url storaging", zap.String("URL to shorten", string(urlToShorten)))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
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
