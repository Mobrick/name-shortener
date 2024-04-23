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

	if flag.Lookup("a") == nil {
		flag.StringVar(&config.FlagRunAddr, "a", ":8080", "address to run server")
	}
	if flag.Lookup("b") == nil {
		flag.StringVar(&config.FlagShortURLBaseAddr, "b", "http://localhost:8080/", "base address of shortened URL")
	}
	if flag.Lookup("l") == nil {
		flag.StringVar(&config.FlagLogLevel, "l", "info", "log level")
	}
	if flag.Lookup("f") == nil {
		flag.StringVar(&config.FlagFileStoragePath, "f", "", "path of file with saved URLs")
	}
	if flag.Lookup("d") == nil {
		flag.StringVar(&config.FlagDBConnectionAddress, "d", "", "database connection address")
	}

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
