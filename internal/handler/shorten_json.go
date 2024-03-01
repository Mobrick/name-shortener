package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/logger"
	"github.com/Mobrick/name-shortener/urltf"
	"go.uber.org/zap"
)

func (env HandlerEnv) LongURLFromJSONHandle(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var request models.Request
	var buf bytes.Buffer

	userId, _ := GetUserIdFromRequest(req)

	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &request); err != nil {
		logger.Log.Debug("could not unmarshal request", zap.String("Requset URL", request.URL))
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	storage := env.Storage

	urlToShorten := []byte(request.URL)

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr

	encodedURL, err := urltf.EncodeURL(urlToShorten, ShortURLLength)
	if err != nil {
		logger.Log.Debug("could not copmplete url encoding", zap.String("URL to encode", string(urlToShorten)))		
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	existingShortURL, err := storage.Add(ctx, encodedURL, string(urlToShorten), userId)
	if err != nil {
		logger.Log.Debug("could not complete url storaging", zap.String("URL to shorten", string(urlToShorten)))		
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	var response models.Response
	var status int
	if len(existingShortURL) != 0 {
		response = models.Response{
			Result: urltf.MakeResultShortenedURL(hostAndPathPart, existingShortURL, req),
		}
		status = http.StatusConflict
	} else {
		response = models.Response{
			Result: urltf.MakeResultShortenedURL(hostAndPathPart, encodedURL, req),
		}
		status = http.StatusCreated
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger.Log.Debug("could not marshal response", zap.String("Response result", response.Result))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	res.Write([]byte(resp))
}
