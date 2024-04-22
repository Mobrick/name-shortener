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

func (dbData MockDB) GetUrlsByUserId(context.Context, string, string, *http.Request) ([]models.SimpleURLRecord, error) {
	return nil, nil
}

func (dbData MockDB) PingDB() error {
	return nil
}
