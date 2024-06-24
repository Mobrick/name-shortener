package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mobrick/name-shortener/internal/auth"
	"github.com/Mobrick/name-shortener/internal/compression"
	"github.com/Mobrick/name-shortener/internal/config"
	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/handler"
	"github.com/Mobrick/name-shortener/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version: %s", buildVersion)
	fmt.Printf("Build date: %s", buildDate)
	fmt.Printf("Build commit: %s", buildCommit)
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()

	sugar := *zapLogger.Sugar()
	logger.Sugar = sugar

	cfg := config.MakeConfig()

	sugar.Info(cfg.FlagDBConnectionAddress + " " + cfg.FlagFileStoragePath)
	// Определение типа стораджа и создание соотвествующего объекта чтобы потом положить в хендлер

	env := &handler.Env{
		ConfigStruct: cfg,
		Storage:      database.NewDB(cfg.FlagFileStoragePath, cfg.FlagDBConnectionAddress),
	}
	// Добавить Close в интерфейс и закрвать через интерфейс
	defer env.Storage.Close()

	r := chi.NewRouter()

	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(compression.DecompressMiddleware)
	r.Use(logger.LoggingMiddleware)
	r.Use(auth.CookieMiddleware)

	r.Get(`/{shortURL}`, env.ShortenedURLHandle)
	r.Get(`/ping`, env.PingDBHandle)
	r.Get(`/api/user/urls`, env.UserUrlsHandler)
	r.Get(`/api/internal/stats`, env.StatsHandle)

	r.Post(`/`, env.LongURLHandle)
	r.Post(`/api/shorten`, env.LongURLFromJSONHandle)
	r.Post(`/api/shorten/batch`, env.BatchHandler)

	r.Delete(`/api/user/urls`, env.DeleteUserUsrlsHandler)

	r.Mount("/debug", middleware.Profiler())

	sugar.Infow(
		"Starting server",
		"addr", cfg.FlagShortURLBaseAddr,
	)

	server := &http.Server{
		Addr:    cfg.FlagRunAddr,
		Handler: r,
	}

	go func() {
		if cfg.FlagEnableHTTPS {
			if err := server.ListenAndServeTLS(cfg.CertFilepath, cfg.KeyFilepath); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-stop

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	env.Storage.Close()
	sugar.Infow("Server stopped")
}
