package config

import (
	"flag"
	"os"
)

type Config struct {
	FlagRunAddr             string
	FlagShortURLBaseAddr    string
	FlagLogLevel            string
	FlagFileStoragePath     string
	FlagDBConnectionAddress string
}

func MakeConfig() *Config {
	config := &Config{}

	flag.StringVar(&config.FlagRunAddr, "a", ":8080", "address to run server")
	flag.StringVar(&config.FlagShortURLBaseAddr, "b", "http://localhost:8080/", "base address of shortened URL")
	flag.StringVar(&config.FlagLogLevel, "l", "info", "log level")
	flag.StringVar(&config.FlagFileStoragePath, "f", "tmp/short-url-db.json", "path of file with saved URLs")
	flag.StringVar(&config.FlagDBConnectionAddress, "d", "host=127.0.0.1 port=5432 user=postgres password=vbhjyjdrf dbname=URLRecords sslmode=disable", "database connection address")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		config.FlagRunAddr = envRunAddr
	}

	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		config.FlagShortURLBaseAddr = envBaseAddr
	}

	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		config.FlagLogLevel = envLogLevel
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		config.FlagFileStoragePath = envFileStoragePath
	}

	if envDBConnectionAddress := os.Getenv("DATABASE_DSN"); envDBConnectionAddress != "" {
		config.FlagDBConnectionAddress = envDBConnectionAddress
	}

	return config
}
