package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/logger"
	"github.com/Mobrick/name-shortener/urltf"
)

func (env HandlerEnv) BatchHandler(res http.ResponseWriter, req *http.Request) {
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

	responseRecords := processMultipleURLRecords(env, urls, req)

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

func processMultipleURLRecords(env HandlerEnv, urlsToShorten []models.BatchRequestURL, req *http.Request) []models.BatchResponseURL {
	var responseRecords []models.BatchResponseURL
	db := env.DatabaseData
	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr
	shortURLRequestMap := make(map[string]models.BatchRequestURL)

	// Creating shorten urls for each record in request
	for _, originalURLRecord := range urlsToShorten {
		shortAddress, shortURL := urltf.MakeShortAddressAndURL(hostAndPathPart, db, []byte(originalURLRecord.OriginalURL), req, ShortURLLength)
		responseRecord := models.BatchResponseURL{
			CorrelationID: originalURLRecord.CorrelationID,
			ShortURL:      shortAddress,
		}
		shortURLRequestMap[shortURL] = originalURLRecord
		responseRecords = append(responseRecords, responseRecord)
	}

	db.AddMany(shortURLRequestMap)

	return responseRecords
}
