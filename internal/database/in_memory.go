package database

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
)

// для хранения в памяти

type InMemoryDB struct {
	URLRecords  []models.URLRecord
	DatabaseMap map[string]string
}

func (dbData InMemoryDB) PingDB() error {
	return errors.New("missing database connection")
}

func (dbData InMemoryDB) Get(ctx context.Context, shortURL string) (string, bool, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok, false, nil
}

func (dbData *InMemoryDB) Add(ctx context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)
	dbData.URLRecords = append(dbData.URLRecords, newRecord)

	return "", nil
}

func (dbData *InMemoryDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL, userID string) error {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userID)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
	}
	return nil
}

func (dbData InMemoryDB) Close() {
}

func (dbData InMemoryDB) GetUrlsByUserID(ctx context.Context, userID string, hostAndPathPart string, req *http.Request) ([]models.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userID, hostAndPathPart, req)
	return usersUrls, nil
}

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
