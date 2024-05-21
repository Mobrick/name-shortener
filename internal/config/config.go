package config

import (
	"flag"
	"os"
)

// Config хранит данные по флагам.
type Config struct {
	FlagRunAddr             string // адрес на котором запущен сервер
	FlagShortURLBaseAddr    string // базовый адрес сокращенного URL
	FlagLogLevel            string // уровень логировани
	FlagFileStoragePath     string // путь к файлу с сохраненными URL
	FlagDBConnectionAddress string // строка подключения к БД
	FlagEnableHTTPS         bool   // использовать ли HTTPS

	CertFilepath string // Путь к сертификату
	KeyFilepath  string // Путь к ключу
}

// MakeConfig формирует конфигурацию по флагам, либо если есть, по переменным окружения.
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
	if flag.Lookup("s") == nil {
		flag.BoolVar(&config.FlagEnableHTTPS, "s", false, "database connection address")
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

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		config.FlagEnableHTTPS = true
	}

	config.CertFilepath = "tls/cert.cer"
	config.KeyFilepath = "tls/key.cer"

	return config
}
