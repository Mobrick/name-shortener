package handler

import (
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/userauth"
)

// Env - структура окружения для хендлеров, в которой хранятся данные о хранилище и ссылка на конфигурацию
type Env struct {
	ConfigStruct *config.Config
	Storage      database.Storage
}

// ShortURLLength константа отражающая количество символов до которого происходит сокращение адреса
const (
	ShortURLLength = 8
)

// GetUserIDFromRequest возвращает id пользователя, который вызвал обработчик, либо ничего
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
