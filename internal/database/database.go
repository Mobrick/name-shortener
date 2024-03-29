package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const urlRecordsTableName = "url_records"

type StorageType int

const (
	SQLDB StorageType = iota
	File
	Local
)

type DatabaseData struct {
	URLRecords         []models.URLRecord
	DatabaseMap        map[string]string
	FileStorage        *os.File
	DatabaseConnection *sql.DB
}

type Storage interface {
	Add(context.Context, string, string) (string, error)
	AddMany(context.Context, map[string]models.BatchRequestURL) error
	PingDB() error
	Get(context.Context, string) (string, bool, error)
	Close()
}

func NewDB(fileName string, connectionString string) Storage {
	var dbData Storage

	if len(connectionString) != 0 {
		dbData = PostgreDB{
			DatabaseMap:        make(map[string]string),
			DatabaseConnection: NewDBConnection(connectionString),
		}
	} else if len(fileName) != 0 {
		dbData = NewDBFromFile(fileName)
	} else {
		dbData = &InMemoryDB{
			URLRecords:  make([]models.URLRecord, 0),
			DatabaseMap: make(map[string]string),
		}
	}

	return dbData
}

func NewDBFromFile(fileName string) Storage {
	file, err := filestorage.MakeFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	urlRecords, err := filestorage.LoadURLRecords(file)
	if err != nil {
		panic(err)
	}

	dbMap, urlRecords := dbMapFromURLRecords(urlRecords)
	databaseData := &FileDB{
		URLRecords:  urlRecords,
		DatabaseMap: dbMap,
		FileStorage: file,
	}

	return databaseData
}

func NewDBConnection(connectionString string) *sql.DB {
	// Закрывается в основном потоке
	db, err := sql.Open("pgx", connectionString)
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

func CreateRecordAndUpdateDBMap(dbMap map[string]string, originalURL string, shortURL string, id string) models.URLRecord {
	newRecord := models.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		UUID:        id,
	}

	dbMap[shortURL] = originalURL
	return newRecord
}
