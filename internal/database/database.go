package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
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
	StorageType        StorageType
}

func (dbData DatabaseData) Get(shortURL string) (string, bool) {
	switch dbData.StorageType {
	case SQLDB:
		location, ok := dbData.GetLocationFromSQLDB(shortURL)
		return location, ok
	default:
		location, ok := dbData.DatabaseMap[shortURL]
		return location, ok
	}

}

func (dbData DatabaseData) GetLocationFromSQLDB(shortURL string) (string, bool) {
	var location string

	row := dbData.DatabaseConnection.QueryRow("SELECT original_url FROM url_records WHERE short_url = $1", shortURL)
	err := row.Scan(&location)
	if err == sql.ErrNoRows {
		return location, false
	}
	return location, true
}

func NewDB(fileName string, connectionString string) DatabaseData {
	var dbData DatabaseData

	if len(connectionString) != 0 {
		dbData = DatabaseData{
			StorageType:        SQLDB,
			URLRecords:         make([]models.URLRecord, 0),
			DatabaseMap:        make(map[string]string),
			DatabaseConnection: NewDBConnection(connectionString),
		}
	} else if len(fileName) != 0 {
		dbData = NewDBFromFile(fileName)
	} else {
		dbData = DatabaseData{
			StorageType: Local,
			URLRecords:  make([]models.URLRecord, 0),
			DatabaseMap: make(map[string]string),
		}
	}

	return dbData
}

func NewDBFromFile(fileName string) DatabaseData {
	file, err := filestorage.MakeFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	urlRecords, err := filestorage.LoadURLRecords(file)
	if err != nil {
		panic(err)
	}

	dbMap, urlRecords := dbMapFromURLRecords(urlRecords)
	databaseData := DatabaseData{
		StorageType: File,
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

func (dbData DatabaseData) Add(shortURL string, originalURL string) {
	id := uuid.New().String()
	newRecord := dbData.createRecordAndUpdateDBMap(originalURL, shortURL, id)

	switch dbData.StorageType {
	case SQLDB:
		dbData.SQLDBAdd(newRecord)
	case File:
		dbData.FileAdd(newRecord)
	default:
		dbData.LocalAdd(newRecord)
	}
}

func (dbData DatabaseData) SQLDBAdd(urlRecord models.URLRecord) {
	dbData.CreateURLRecordsTableIfNotExisits()

	_, err := dbData.DatabaseConnection.Exec("INSERT INTO url_records (uuid, short_url, original_url)"+
		" VALUES ($1, $2, $3)", urlRecord.UUID, urlRecord.ShortURL, urlRecord.OriginalURL)
	if err != nil {
		log.Fatal("Failed to insert a record: "+urlRecord.OriginalURL, err)
		return
	}
}

func (dbData DatabaseData) CreateURLRecordsTableIfNotExisits() {
	_, err := dbData.DatabaseConnection.Exec("CREATE TABLE IF NOT EXISTS " + urlRecordsTableName + " (uuid VARCHAR(255) PRIMARY KEY, short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) NOT NULL)")

	if err != nil {
		log.Fatal(err)
	}
}

func (dbData DatabaseData) FileAdd(newRecord models.URLRecord) {
	dbData.LocalAdd(newRecord)
	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
}

func (dbData DatabaseData) LocalAdd(newRecord models.URLRecord) {
	dbData.URLRecords = append(dbData.URLRecords, newRecord)
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

func (dbData DatabaseData) AddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	switch dbData.StorageType {
	case SQLDB:
		dbData.SQLDBAddMany(shortURLRequestMap)
	case File:
		dbData.FileAddMany(shortURLRequestMap)
	default:
		dbData.LocalAddMany(shortURLRequestMap)
	}
}

func (dbData DatabaseData) LocalAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	for shortURL, record := range shortURLRequestMap {
		newRecord := dbData.createRecordAndUpdateDBMap(record.OriginalURL, shortURL, record.CorrelationID)
		dbData.LocalAdd(newRecord)
	}
}

func (dbData DatabaseData) FileAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	for shortURL, record := range shortURLRequestMap {
		newRecord := dbData.createRecordAndUpdateDBMap(record.OriginalURL, shortURL, record.CorrelationID)
		dbData.LocalAdd(newRecord)
		filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
	}
}

func (dbData DatabaseData) SQLDBAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	// Создание списка всех записей
	var sliceOfRecords []models.URLRecord
	for shortURL, record := range shortURLRequestMap {
		newRecord := dbData.createRecordAndUpdateDBMap(record.OriginalURL, shortURL, record.CorrelationID)
		sliceOfRecords = append(sliceOfRecords, newRecord)
	}

	tx, err := dbData.DatabaseConnection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO url_records (uuid, short_url, original_url)" +
		" VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, record := range sliceOfRecords {
		_, err := stmt.Exec(record.UUID, record.ShortURL, record.OriginalURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (dbData DatabaseData) createRecordAndUpdateDBMap(originalURL string, shortURL string, id string) models.URLRecord {
	newRecord := models.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		UUID:        id,
	}

	dbData.DatabaseMap[shortURL] = originalURL
	return newRecord
}
