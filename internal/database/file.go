package database

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"

	"github.com/Mobrick/name-shortener/internal/filestorage"
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

func (dbData FileDB) Get(ctx context.Context, shortURL string) (string, bool, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok, false, nil
}

func (dbData *FileDB) Add(ctx context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)

	dbData.URLRecords = append(dbData.URLRecords, newRecord)
	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)

	return "", nil
}

func (dbData *FileDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL, userID string) error {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userID)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
		filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
	}
	return nil
}

func (dbData FileDB) Close() {
	dbData.FileStorage.Close()
}

func (dbData FileDB) GetUrlsByUserID(ctx context.Context, userID string, hostAndPathPart string, req *http.Request) ([]models.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userID, hostAndPathPart, req)
	return usersUrls, nil
}

func (dbData *FileDB) Delete(ctx context.Context, urlsToDelete []string, userID string) error {
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
