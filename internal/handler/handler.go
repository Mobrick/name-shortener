package handler

import (
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/userauth"
)

type HandlerEnv struct {
	//DatabaseData database.DatabaseData
	ConfigStruct *config.Config
	Storage      database.Storage
}

const (
	ShortURLLength = 8
)

func GetUserIDFromRequest(req *http.Request) (string, bool) {
	cookie, err := req.Cookie("auth_token")
	if err != nil {
		log.Printf("no cookie found. " + err.Error())
		return "", false
	}

	token := cookie.Value
	userID, ok := userauth.GetUserID(token)
	if !ok {
		log.Printf("invalid token")
		return "", false
	}
	return userID, true
}
