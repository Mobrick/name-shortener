package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/logger"
	"github.com/Mobrick/name-shortener/urltf"
)

func (env HandlerEnv) BatchHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	var buf bytes.Buffer

	// читаем тело запроса
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		logger.Log.Debug("could not unmarshal request")
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	var urls []models.BatchRequestURL
	if err = json.Unmarshal(buf.Bytes(), &urls); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	responseRecords, err := processMultipleURLRecords(ctx, env, urls, req)
	if err != nil {
		logger.Log.Debug("could not copmplete url storaging")		
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(responseRecords)
	if err != nil {
		logger.Log.Debug("could not marshal response")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resp))
}

func processMultipleURLRecords(ctx context.Context, env HandlerEnv, urlsToShorten []models.BatchRequestURL, req *http.Request) ([]models.BatchResponseURL, error) {
	var responseRecords []models.BatchResponseURL
	storage := env.Storage
	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr
	shortURLRequestMap := make(map[string]models.BatchRequestURL)

	// Creating shorten urls for each record in request
	for _, originalURLRecord := range urlsToShorten {
		encodedURL, err := urltf.EncodeURL([]byte(originalURLRecord.OriginalURL), ShortURLLength)
		if err != nil {
			return nil, err
		}
		shortAddress := urltf.MakeResultShortenedURL(hostAndPathPart, encodedURL, req)

		responseRecord := models.BatchResponseURL{
			CorrelationID: originalURLRecord.CorrelationID,
			ShortURL:      shortAddress,
		}
		shortURLRequestMap[encodedURL] = originalURLRecord
		responseRecords = append(responseRecords, responseRecord)
	}

	err := storage.AddMany(ctx, shortURLRequestMap)
	if err != nil {
		return nil, err
	}

	return responseRecords, err
}
