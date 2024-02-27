package database

import (
	"context"
	"errors"

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

func (dbData InMemoryDB) Get(ctx context.Context, shortURL string) (string, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok, nil
}

func (dbData *InMemoryDB) Add(ctx context.Context, shortURL string, originalURL string, userId string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userId)
	dbData.URLRecords = append(dbData.URLRecords, newRecord)

	return "", nil
}

func (dbData *InMemoryDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL, userId string) error {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userId)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
	}
	return nil
}

func (dbData InMemoryDB) Close() {
	return
}

func (dbData InMemoryDB) GetUrlsByUserId(ctx context.Context, userId string) ([]models.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userId)
	return usersUrls, nil
}
