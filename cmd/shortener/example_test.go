package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mobrick/name-shortener/internal/compression"
	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/handler"
	"github.com/Mobrick/name-shortener/internal/logger"
	"github.com/Mobrick/name-shortener/internal/mocks"
	"github.com/Mobrick/name-shortener/internal/userauth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func ExampleEnv_ShortenedURLHandle() {
	env := &handler.Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()

	r := chi.NewRouter()

	r.Get(`/{shortURL}`, env.ShortenedURLHandle)

	server := &http.Server{
		Addr:    env.ConfigStruct.FlagRunAddr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Пример с мидлварями
func ExampleEnv_LongURLFromJSONHandle() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	sugar := *zapLogger.Sugar()
	logger.Sugar = sugar

	cfg := config.MakeConfig()

	sugar.Info(cfg.FlagDBConnectionAddress + " " + cfg.FlagFileStoragePath)

	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(compression.DecompressMiddleware)
	r.Use(logger.LoggingMiddleware)
	r.Use(userauth.CookieMiddleware)
	env := &handler.Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()

	r.Post(`/`, env.LongURLHandle)

	server := &http.Server{
		Addr:    env.ConfigStruct.FlagRunAddr,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

// Пример с красивой остановкой сервера
func ExampleEnv_DeleteUserUsrlsHandler() {
	env := &handler.Env{
		Storage:      mocks.NewMockDB(),
		ConfigStruct: config.MakeConfig(),
	}
	defer env.Storage.Close()

	r := chi.NewRouter()

	r.Delete(`/api/user/urls`, env.DeleteUserUsrlsHandler)

	server := &http.Server{
		Addr:    env.ConfigStruct.FlagRunAddr,
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

	env.Storage.Close()
}
