package main

import (
    "flag"
    "os"
)

var flagRunAddr string
var flagShortURLBaseAddr string



func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address to run server")
	flag.StringVar(&flagShortURLBaseAddr, "b", "http://localhost:8080/", "base address of shortened URL")

	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
        flagRunAddr = envRunAddr
    }

	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
        flagShortURLBaseAddr = envBaseAddr
    }
}