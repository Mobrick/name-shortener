package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/logger"
)

// UserUrlsHandler возвращает адреса созданные пользователем.
func (env Env) UserUrlsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, ok := GetUserIDFromRequest(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr

	urls, err := env.Storage.GetUrlsByUserID(ctx, userID, hostAndPathPart, req)
	if err != nil {
		logger.Log.Debug("could not get urls by user id")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	resp, err := json.Marshal(urls)
	if err != nil {
		logger.Log.Debug("could not marshal response")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(resp))
}
