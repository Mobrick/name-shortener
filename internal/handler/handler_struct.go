package handler

import (
	"github.com/Mobrick/name-shortener/config"
)

type HandlerEnv struct {
	DatabaseMap  map[string]string
	ConfigStruct *config.Config
}

const (
	ShortURLLength = 8
)
