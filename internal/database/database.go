package database

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/urltf"
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
	Add(context.Context, string, string, string) (string, error)
	AddMany(context.Context, map[string]models.BatchRequestURL, string) error
	PingDB() error
	Get(context.Context, string) (string, bool, error)
	GetUrlsByUserId(context.Context, string, string, *http.Request) ([]models.SimpleURLRecord, error)
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

func CreateRecordAndUpdateDBMap(dbMap map[string]string, originalURL string, shortURL string, id string, userId string) models.URLRecord {
	newRecord := models.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		UUID:        id,
		UserID:      userId,
	}

	dbMap[shortURL] = originalURL
	return newRecord
}

func GetUrlsCreatedByUser(urlRecords []models.URLRecord, userId string, hostAndPathPart string, req *http.Request) []models.SimpleURLRecord {
	var usersUrls []models.SimpleURLRecord
	for _, urlRecord := range urlRecords {
		if urlRecord.UserID == userId {
			usersUrl := models.SimpleURLRecord{
				ShortURL:    urltf.MakeResultShortenedURL(hostAndPathPart, urlRecord.ShortURL, req),
				OriginalURL: urlRecord.OriginalURL,
			}
			usersUrls = append(usersUrls, usersUrl)
		}
	}
	return usersUrls
}
