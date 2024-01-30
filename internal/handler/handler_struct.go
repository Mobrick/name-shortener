package handler

import (
	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/database"
)

type HandlerEnv struct {
	DatabaseData  database.DatabaseData
	ConfigStruct *config.Config
}

const (
	ShortURLLength = 8
)
