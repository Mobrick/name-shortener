package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/logger"
	"go.uber.org/zap"
)

// DeleteUserUsrlsHandler удаляет URL из хранилища переданные пользователем, которые были созданы им.
func (env Env) DeleteUserUsrlsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, _ := GetUserIDFromRequest(req)

	urlsToDelete, err := parseRequestBody(req)
	if err != nil {
		logger.Log.Debug("Error parsing request body", zap.String("error message: ", err.Error()))
		return
	}

	if len(urlsToDelete) == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	storage := env.Storage

	err = storage.Delete(ctx, urlsToDelete, userID)
	if err != nil {
		logger.Log.Debug("could not delete urls")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func parseRequestBody(req *http.Request) ([]string, error) {
	urlsFromBody, err := io.ReadAll(io.Reader(req.Body))
	if err != nil {
		return nil, err
	}

	var data []string
	err = json.Unmarshal(urlsFromBody, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
