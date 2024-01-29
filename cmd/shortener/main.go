package main

import (
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/database"
	"github.com/Mobrick/name-shortener/handler"
	"github.com/Mobrick/name-shortener/internal/compression"
	"github.com/Mobrick/name-shortener/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	sugar := *zapLogger.Sugar()

	logger.Sugar = sugar

	env := &handler.HandlerEnv{
		DatabaseMap:  database.NewDBMap(),
		ConfigStruct: config.MakeConfig(),
	}
	r := chi.NewRouter()

	r.Use(compression.GzipMiddleware)
	r.Use(logger.LoggingMiddleware)

	r.Get(`/{shortURL}`, env.ShortenedURLHandle)

	r.Post(`/`, env.LongURLHandle)
	r.Post(`/api/shorten`, env.LongURLFromJSONHandle)

	sugar.Infow(
		"Starting server",
		"addr", env.ConfigStruct.FlagShortURLBaseAddr,
	)

	log.Fatal(http.ListenAndServe(env.ConfigStruct.FlagRunAddr, r))
}
