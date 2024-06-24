package database

import (
	"context"
	"errors"
	"net/http"
	"os"
	"slices"

	"github.com/Mobrick/name-shortener/internal/filestorage"
	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/google/uuid"
)

// FileDB предствялет собой базу данных построенную по файлу при включении сервера.
type FileDB struct {
	DatabaseMap map[string]string
	URLRecords  []model.URLRecord
	FileStorage *os.File
}

// PingDB возвращает ошибку так как хранилище локальное.
func (dbData FileDB) PingDB() error {
	return errors.New("missing database connection")
}

// Get возвращает оригинальный URL.
func (dbData FileDB) Get(_ context.Context, shortURL string) (string, bool, error) {
	location, ok := dbData.DatabaseMap[shortURL]
	if !ok {
		return "", false, nil
	}
	return location, false, nil
}

// Add добавляет данные о сокращенном URL в хранилище.
func (dbData *FileDB) Add(_ context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)

	dbData.URLRecords = append(dbData.URLRecords, newRecord)
	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)

	return "", nil
}

// AddMany добавляет множество данных о сокращенных URL в хранилище.
func (dbData *FileDB) AddMany(_ context.Context, shortURLRequestMap map[string]model.BatchRequestURL, userID string) error {
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userID)
		dbData.URLRecords = append(dbData.URLRecords, newRecord)
		filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
	}
	return nil
}

// Close закрывает подключение к файловому хранилищу.
func (dbData FileDB) Close() {
	dbData.FileStorage.Close()
}

// GetUrlsByUserID возвращает записи созданные пользователем.
func (dbData FileDB) GetUrlsByUserID(_ context.Context, userID string, hostAndPathPart string, req *http.Request) ([]model.SimpleURLRecord, error) {
	urlRecords := dbData.URLRecords
	usersUrls := GetUrlsCreatedByUser(urlRecords, userID, hostAndPathPart, req)
	return usersUrls, nil
}

// Delete удаляет данные о сокращенном URL из хранилища.
func (dbData *FileDB) Delete(_ context.Context, urlsToDelete []string, userID string) error {
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
