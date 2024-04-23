package mocks

import (
	"context"
	"net/http"

	"github.com/Mobrick/name-shortener/database"
	"github.com/Mobrick/name-shortener/internal/models"
)

type MockDB struct {
}

func NewMockDB() database.Storage {
	return MockDB{}
}

func (dbData MockDB) Add(_ context.Context, _ string, encodedURL string, _ string) (string, error) {
	if encodedURL == "https://www.go.com/" {
		return "dupeDupe", nil
	}
	return "", nil
}

func (dbData MockDB) AddMany(context.Context, map[string]models.BatchRequestURL, string) error {
	return nil
}

func (dbData MockDB) Close() {
}

func (dbData MockDB) Delete(context.Context, []string, string) error {
	return nil
}

func (dbData MockDB) Get(context.Context, string) (string, bool, bool, error) {
	return "", false, false, nil
}

func (dbData MockDB) GetUrlsByUserID(_ context.Context, userID string, _ string, _ *http.Request) ([]models.SimpleURLRecord, error) {
	if userID == "1a91a181-80ec-45cb-a576-14db11505700" {
		urls := []models.SimpleURLRecord{
			{
				ShortURL:    "DDDDdddd",
				OriginalURL: "https://www.google.com/",
			},
			{
				ShortURL:    "vvvv4444",
				OriginalURL: "https://www.go.com/",
			},
		}
		return urls, nil
	} else {
		urls := []models.SimpleURLRecord{}
		return urls, nil
	}
}

func (dbData MockDB) PingDB() error {
	return nil
}
