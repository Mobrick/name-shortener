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
	var request models.Request
	var buf bytes.Buffer

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
	db := env.DatabaseData
	urlToShorten := []byte(request.URL)

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr

	encodedURL := urltf.EncodeURL(urlToShorten, db, ShortURLLength)
	
	existingShortURL := db.Add(encodedURL, string(urlToShorten))

	var response models.Response
	var status int
	if (len(existingShortURL) != 0)	{
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
