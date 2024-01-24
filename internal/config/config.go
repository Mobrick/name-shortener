package config

import (
    "flag"
    "os"
)

type Config struct {
	FlagRunAddr string
    FlagShortURLBaseAddr string
}

func MakeConfig() *Config {
    config := &Config{}

	flag.StringVar(&config.FlagRunAddr, "a", ":8080", "address to run server")
	flag.StringVar(&config.FlagShortURLBaseAddr, "b", "http://localhost:8080/", "base address of shortened URL")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
        config.FlagRunAddr = envRunAddr
    }

	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
        config.FlagShortURLBaseAddr = envBaseAddr
    }

    return config
}