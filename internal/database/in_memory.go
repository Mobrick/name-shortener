package database

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/google/uuid"
)

// FileDB предствялет собой базу данных которая хранится в памяти.
type InMemoryDB struct {
	URLRecords  []model.URLRecord
	DatabaseMap map[string]string
}

// PingDB возвращает ошибку так как хранилище локальное.
func (dbData InMemoryDB) PingDB() error {
	return errors.New("missing database connection")
}

// Get возвращает оригинальный URL.
func (dbData InMemoryDB) Get(_ context.Context, shortURL string) (string, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	if !ok {
		return "", false, nil
	}
	return location, false, nil
}

// Add добавляет данные о сокращенном URL в хранилище.
func (dbData *InMemoryDB) Add(_ context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)
	dbData.URLRecords = append(dbData.URLRecords, newRecord)

	return "", nil
}

// AddMany добавляет множество данных о сокращенных URL в хранилище.
func (dbData *InMemoryDB) AddMany(_ context.Context, shortURLRequestMap map[string]model.BatchRequestURL, userID string) error {
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
func (dbData InMemoryDB) GetUrlsByUserID(_ context.Context, userID string, hostAndPathPart string, req *http.Request) ([]model.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userID, hostAndPathPart, req)
	return usersUrls, nil
}

// Delete удаляет данные о сокращенном URL из хранилища.
func (dbData *InMemoryDB) Delete(_ context.Context, urlsToDelete []string, userID string) error {
	var result []model.URLRecord
	for _, urlRecord := range dbData.URLRecords {
		if urlRecord.UserID != userID {
			result = append(result, urlRecord)
			continue
		}
		if !slices.Contains(urlsToDelete, urlRecord.ShortURL) {
			result = append(result, urlRecord)
			continue
		}
	}

	dbData.URLRecords = result
	return nil
}
