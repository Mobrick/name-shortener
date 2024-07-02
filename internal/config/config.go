package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/Mobrick/name-shortener/internal/model"
)

// Config хранит данные по флагам.
type Config struct {
	FlagConfigFile          string // имя файла конфигурации
	FlagRunAddr             string // адрес на котором запущен сервер
	FlagShortURLBaseAddr    string // базовый адрес сокращенного URL
	FlagLogLevel            string // уровень логировани
	FlagFileStoragePath     string // путь к файлу с сохраненными URL
	FlagDBConnectionAddress string // строка подключения к БД
	FlagEnableHTTPS         bool   // использовать ли HTTPS
	FlagTrustedSubnet       string // строковое представление беклассовой адресации

	CertFilepath string // Путь к сертификату
	KeyFilepath  string // Путь к ключу
}

// MakeConfig формирует конфигурацию по флагам, либо если есть, по переменным окружения.
func MakeConfig() *Config {
	config := &Config{}

	if flag.Lookup("c") == nil {
		flag.StringVar(&config.FlagConfigFile, "c", "", "configuration file name")
	}
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
	if flag.Lookup("t") == nil {
		flag.StringVar(&config.FlagTrustedSubnet, "t", "", "database connection address")
	}
	if flag.Lookup("s") == nil {
		flag.BoolVar(&config.FlagEnableHTTPS, "s", false, "database connection address")
	}

	flag.Parse()

	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		config.FlagConfigFile = envConfig
	}

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

	if envTrustedSubnet := os.Getenv("TRUSTED_SUBNET"); envTrustedSubnet != "" {
		config.FlagTrustedSubnet = envTrustedSubnet
	}

	if envEnableHTTPS := os.Getenv("ENABLE_HTTPS"); envEnableHTTPS != "" {
		config.FlagEnableHTTPS = true
	}

	config.CertFilepath = "tls/cert.cer"
	config.KeyFilepath = "tls/key.cer"

	if config.FlagConfigFile != "" {
		configFromFile, err := getConfigFromFile(config.FlagConfigFile)
		if err != nil {
			log.Fatal(err)
		}

		if config.FlagRunAddr == "" {
			config.FlagRunAddr = configFromFile.ServerAddress
		}

		if config.FlagShortURLBaseAddr == "" {
			config.FlagShortURLBaseAddr = configFromFile.BaseURL
		}

		if config.FlagFileStoragePath == "" {
			config.FlagFileStoragePath = configFromFile.FileStoragePath
		}

		if config.FlagDBConnectionAddress == "" {
			config.FlagDBConnectionAddress = configFromFile.DatabaseDsn
		}

		if config.FlagTrustedSubnet == "" {
			config.FlagTrustedSubnet = configFromFile.TrustedSubnet
		}

		if !config.FlagEnableHTTPS {
			config.FlagEnableHTTPS = configFromFile.EnableHTTPS
		}
	}

	return config
}

func getConfigFromFile(filename string) (model.ConfigFromFile, error) {
	var config model.ConfigFromFile
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
