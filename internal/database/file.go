package database

import (
	"context"
	"errors"
	"os"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
)

type FileDB struct {
	DatabaseMap map[string]string
	URLRecords  []models.URLRecord
	FileStorage *os.File
}

// для хранения в файле
func (dbData FileDB) PingDB() error {
	return errors.New("missing database connection")
}

func (dbData FileDB) Get(ctx context.Context, shortURL string) (string, bool) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok
}

func (dbData *FileDB) Add(ctx context.Context, shortURL string, originalURL string) string {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id)

	dbData.URLRecords = append(dbData.URLRecords, newRecord)
	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)

	return ""
}

func (dbData *FileDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL) {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
		filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
	}
}

func (dbData FileDB) Close() {
	dbData.FileStorage.Close()
}