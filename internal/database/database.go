package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
)

type DatabaseData struct {
	URLRecords         []models.URLRecord
	DatabaseMap        map[string]string
	FileStorage        *os.File
	DatabaseConnection *sql.DB
}

func (dbData DatabaseData) Get(shortURL string) (string, bool) {
	location, ok := dbData.DatabaseMap[shortURL]
	return location, ok
}

func NewDB(fileStorage *os.File) DatabaseData {
	dbData := NewDBFromFile(fileStorage)
	dbData.DatabaseConnection = NewDBConnection()
	return dbData
}

func NewDBFromFile(fileStorage *os.File) DatabaseData {
	if fileStorage == nil {
		return DatabaseData{
			URLRecords:  make([]models.URLRecord, 0),
			DatabaseMap: make(map[string]string),
		}
	}
	urlRecords, err := filestorage.LoadURLRecords(fileStorage)
	if err != nil {
		panic(err)
	}

	dbMap, urlRecords := dbMapFromURLRecords(urlRecords)
	databaseData := DatabaseData{
		URLRecords:  urlRecords,
		DatabaseMap: dbMap,
		FileStorage: fileStorage,
	}

	return databaseData
}

func NewDBConnection() *sql.DB {
	data, err := os.ReadFile("connection.txt")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	connectionString := string(data)

	ps := fmt.Sprintf(connectionString)

	db, err := sql.Open("pgx", ps)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func dbMapFromURLRecords(urlRecords []models.URLRecord) (map[string]string, []models.URLRecord) {
	dbMap := make(map[string]string)
	for _, urlRecord := range urlRecords {
		dbMap[urlRecord.ShortURL] = urlRecord.OriginalURL
	}
	return dbMap, urlRecords
}

func (dbData DatabaseData) Add(shortURL string, originalURL string) {
	newRecord := models.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}
	newRecord.UUID = uuid.New().String()

	dbData.URLRecords = append(dbData.URLRecords, newRecord)
	dbData.DatabaseMap[shortURL] = originalURL

	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
}

func (dbData DatabaseData) Contains(shortUrl string) bool {
	dbMap := dbData.DatabaseMap

	if _, ok := dbMap[string(shortUrl)]; !ok {
		return false
	} else {
		return true
	}
}

func (dbData DatabaseData) PingDB() error {
	err := dbData.DatabaseConnection.Ping()
	return err
}
