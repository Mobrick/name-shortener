package database

import (
	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
)

type DatabaseData struct {
	URLRecords  []models.URLRecord
	DatabaseMap map[string]string
}

func NewDBFromFile(fileStoragePath string) DatabaseData {
	if len(fileStoragePath) == 0 {
		return DatabaseData{
			URLRecords:  make([]models.URLRecord, 0),
			DatabaseMap: make(map[string]string),
		}
	}
	urlRecords, err := filestorage.LoadURLRecords(fileStoragePath)
	if err != nil {
		panic(err)
	}

	dbMap, urlRecords := dbMapFromURLRecords(urlRecords)
	databaseData := DatabaseData{
		URLRecords:  urlRecords,
		DatabaseMap: dbMap,
	}

	return databaseData
}

func dbMapFromURLRecords(urlRecords []models.URLRecord) (map[string]string, []models.URLRecord) {
	dbMap := make(map[string]string)
	for _, urlRecord := range urlRecords {
		dbMap[urlRecord.ShortURL] = urlRecord.OriginalURL
	}
	return dbMap, urlRecords
}

func AddNewRecordToDatabase(dbData DatabaseData, shortURL string, originalURL string) {
	newRecord := models.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	newRecord.UUID = generateNewUUID(dbData)

	dbData.DatabaseMap[shortURL] = originalURL
}

func generateNewUUID(dbData DatabaseData) {

}
