package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Mobrick/name-shortener/logger"
)

func (env HandlerEnv) UserUrlsHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userId, ok := GetUserIdFromRequest(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	hostAndPathPart := env.ConfigStruct.FlagShortURLBaseAddr

	urls, err := env.Storage.GetUrlsByUserId(ctx, userId, hostAndPathPart, req)
	if err != nil {
		logger.Log.Debug("could not get urls by user id")
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(urls) == 0 {
		res.WriteHeader(http.StatusNoContent)
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
