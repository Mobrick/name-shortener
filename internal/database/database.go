package database

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/Mobrick/name-shortener/filestorage"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

func (dbData DatabaseData) Add(shortURL string, originalURL string) string {
	id := uuid.New().String()
	newRecord := dbData.createRecordAndUpdateDBMap(originalURL, shortURL, id)

	switch dbData.StorageType {
	case SQLDB:
		return dbData.sqldbAdd(newRecord)
	case File:
		dbData.fileAdd(newRecord)
	default:
		dbData.localAdd(newRecord)
	}
	return ""
}

func (dbData DatabaseData) sqldbAdd(urlRecord models.URLRecord) string {
	dbData.createURLRecordsTableIfNotExists()

	originalURL := urlRecord.OriginalURL

	insertStmt, err := dbData.DatabaseConnection.Prepare("INSERT INTO url_records (uuid, short_url, original_url)" +
		" VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal("Failed to prepare the SQL statement of: "+originalURL, err)
	}

	_, err = insertStmt.Exec(urlRecord.UUID, urlRecord.ShortURL, originalURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("url %s already in database", originalURL)
			return dbData.findExisitingShortURL(originalURL)
		} else {
			log.Fatal("Failed to insert a record: "+originalURL, err)
		}
	}

	return ""
}

func (dbData DatabaseData) findExisitingShortURL(originalURL string) string {
	stmt, err := dbData.DatabaseConnection.Prepare("SELECT short_url FROM url_records WHERE original_url = $1")
	if err != nil {
		log.Fatal("Failed to prepare the statement of getting existing shortURL:  "+originalURL, err)
	}
	var shortURL string
	err = stmt.QueryRow(originalURL).Scan(&shortURL)
	return shortURL
}

func (dbData DatabaseData) createURLRecordsTableIfNotExists() {
	_, err := dbData.DatabaseConnection.Exec("CREATE TABLE IF NOT EXISTS " + urlRecordsTableName + " (uuid VARCHAR(255) PRIMARY KEY, short_url VARCHAR(255) NOT NULL, original_url VARCHAR(255) NOT NULL UNIQUE)")

	if err != nil {
		log.Fatal(err)
	}
}

func (dbData DatabaseData) fileAdd(newRecord models.URLRecord) {
	dbData.localAdd(newRecord)
	filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
}

func (dbData DatabaseData) localAdd(newRecord models.URLRecord) {
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
		dbData.sqldbAddMany(shortURLRequestMap)
	case File:
		dbData.fileAddMany(shortURLRequestMap)
	default:
		dbData.localAddMany(shortURLRequestMap)
	}
}

func (dbData DatabaseData) localAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	for shortURL, record := range shortURLRequestMap {
		newRecord := dbData.createRecordAndUpdateDBMap(record.OriginalURL, shortURL, record.CorrelationID)
		dbData.localAdd(newRecord)
	}
}

func (dbData DatabaseData) fileAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
	for shortURL, record := range shortURLRequestMap {
		newRecord := dbData.createRecordAndUpdateDBMap(record.OriginalURL, shortURL, record.CorrelationID)
		dbData.localAdd(newRecord)
		filestorage.UploadNewURLRecord(newRecord, dbData.FileStorage)
	}
}

func (dbData DatabaseData) sqldbAddMany(shortURLRequestMap map[string]models.BatchRequestURL) {
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
