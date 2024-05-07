package mocks

import (
	"context"
	"errors"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/models"
)

// MockDB - структура для мока
type MockDB struct {
}

// NewMockDB создает мок хранилища
func NewMockDB() database.Storage {
	return MockDB{}
}

// Add добавляет данные о сокращенном URL в хранилище.
func (dbData MockDB) Add(_ context.Context, _ string, encodedURL string, _ string) (string, error) {
	if encodedURL == "https://www.go.com/" {
		return "dupeDupe", nil
	}
	return "", nil
}

// AddMany добавляет множество данных о сокращенных URL в хранилище.
func (dbData MockDB) AddMany(context.Context, map[string]models.BatchRequestURL, string) error {
	return nil
}

// Close закрывает подключение к хранилищу.
func (dbData MockDB) Close() {
}

// Delete удаляет данные о сокращенном URL из хранилища.
func (dbData MockDB) Delete(context.Context, []string, string) error {
	return nil
}

// Get возвращает оригинальный URL, либо сообщает об отсуствии соответсвующего URL, также возвращает пометку об удалении.
func (dbData MockDB) Get(context.Context, string) (string, bool, error) {
	return "", false, nil
}

// GetUrlsByUserID возвращает записи созданные пользователем.
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

// PingDB пингует подключение к бд.
func (dbData MockDB) PingDB() error {
	return errors.New("no connection to database")
}
