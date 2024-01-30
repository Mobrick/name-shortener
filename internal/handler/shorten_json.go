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
	dbMap := env.DatabaseData.DatabaseMap
	urlToShorten := []byte(request.URL)
	shortAddress, shortURL := urltf.MakeShortAddressAndURL(env.ConfigStruct.FlagShortURLBaseAddr, dbMap, urlToShorten, req, ShortURLLength)
	env.DatabaseData.AddNewRecordToDatabase(shortURL, string(urlToShorten), env.ConfigStruct.FlagFileStoragePath)
	response := models.Response{
		Result: shortAddress,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		logger.Log.Debug("could not marshal response", zap.String("Response result", response.Result))
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(resp))
}