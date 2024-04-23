package database

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
)

// FileDB предствялет собой базу данных которая хранится в памяти.
type InMemoryDB struct {
	URLRecords  []models.URLRecord
	DatabaseMap map[string]string
}

// PingDB возвращает ошибку так как хранилище локальное.
func (dbData InMemoryDB) PingDB() error {
	return errors.New("missing database connection")
}

// Get возвращает оригинальный URL.
func (dbData InMemoryDB) Get(ctx context.Context, shortURL string) (string, bool, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok, false, nil
}

// Add добавляет данные о сокращенном URL в хранилище.
func (dbData *InMemoryDB) Add(ctx context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)
	dbData.URLRecords = append(dbData.URLRecords, newRecord)

	return "", nil
}

// AddMany добавляет множество данных о сокращенных URL в хранилище.
func (dbData *InMemoryDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL, userID string) error {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userID)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
	}
	return nil
}

// Close ничего не делает в случае хранилища хранящегося в памяти.
func (dbData InMemoryDB) Close() {
}

// GetUrlsByUserID возвращает записи созданные пользователем.
func (dbData InMemoryDB) GetUrlsByUserID(ctx context.Context, userID string, hostAndPathPart string, req *http.Request) ([]models.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userID, hostAndPathPart, req)
	return usersUrls, nil
}

// Delete удаляет данные о сокращенном URL из хранилища.
func (dbData *InMemoryDB) Delete(ctx context.Context, urlsToDelete []string, userID string) error {
	var result []models.URLRecord
	for _, urlRecord := range dbData.URLRecords {
		if urlRecord.UserID != userID {
			continue
		}
		if !slices.Contains(urlsToDelete, urlRecord.ShortURL) {
			continue
		}
		result = append(result, urlRecord)
	}
	dbData.URLRecords = result
	return nil
}
