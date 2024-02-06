package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mobrick/name-shortener/config"
	"github.com/Mobrick/name-shortener/database"
	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/handler"
	"github.com/Mobrick/name-shortener/internal/compression"
	"github.com/Mobrick/name-shortener/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	cfg := config.MakeConfig()

	file, err := filestorage.MakeFile(cfg.FlagFileStoragePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	env := &handler.HandlerEnv{
		ConfigStruct: cfg,
		DatabaseData: database.NewDBFromFile(file),
	}

	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(compression.DecompressMiddleware)
	r.Use(logger.LoggingMiddleware)

	r.Get(`/{shortURL}`, env.ShortenedURLHandle)

	r.Post(`/`, env.LongURLHandle)
	r.Post(`/api/shorten`, env.LongURLFromJSONHandle)

	sugar.Infow(
		"Starting server",
		"addr", cfg.FlagShortURLBaseAddr,
	)

	server := &http.Server{
		Addr:    cfg.FlagRunAddr,
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	file.Close()
	sugar.Infow("Server stopped")
}
