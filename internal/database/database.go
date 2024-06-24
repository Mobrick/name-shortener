package database

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/filestorage"
	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/Mobrick/name-shortener/pkg/urltf"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const urlRecordsTableName = "url_records"

// Storage описывает поведение хранилища.
type Storage interface {
	// Add добавляет данные о сокращенном URL в хранилище.
	Add(context.Context, string, string, string) (string, error)

	// AddMany добавляет множество данных о сокращенных URL в хранилище.
	AddMany(context.Context, map[string]model.BatchRequestURL, string) error

	// Close закрывает подключение к хранилищу.
	Close()

	// Delete удаляет данные о сокращенном URL из хранилища.
	Delete(context.Context, []string, string) error

	// Get возвращает оригинальный URL, либо сообщает об отсуствии соответсвующего URL, также возвращает пометку об удалении.
	Get(context.Context, string) (string, bool, error)

	// GetUrlsByUserID возвращает записи созданные пользователем.
	GetUrlsByUserID(context.Context, string, string, *http.Request) ([]model.SimpleURLRecord, error)

	// PingDB пингует подключение к бд.
	PingDB() error
}

// NewDB создает объект хранилища в зависимости от того, чем заполнены флаги.
func NewDB(fileName string, connectionString string) Storage {
	var dbData Storage

	switch {
	case len(connectionString) != 0:
		dbData = PostgreDB{
			DatabaseMap:        make(map[string]string),
			DatabaseConnection: newDBConnection(connectionString),
		}
	case len(fileName) != 0:
		dbData = newDBFromFile(fileName)
	default:
		dbData = &InMemoryDB{
			URLRecords:  make([]model.URLRecord, 0),
			DatabaseMap: make(map[string]string),
		}
	}

	return dbData
}

// newDBFromFile формирует БД в памяти по содержимогу файла, путь к которому указан в флаге.
func newDBFromFile(fileName string) Storage {
	file, err := filestorage.MakeFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	urlRecords, err := filestorage.LoadURLRecords(file)
	if err != nil {
		panic(err)
	}

	dbMap := dbMapFromURLRecords(urlRecords)
	databaseData := &FileDB{
		URLRecords:  urlRecords,
		DatabaseMap: dbMap,
		FileStorage: file,
	}

	return databaseData
}

// newDBConnection создает подключение к базе данные Postgre.
func newDBConnection(connectionString string) *sql.DB {
	// Закрывается в основном потоке
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

func dbMapFromURLRecords(urlRecords []model.URLRecord) map[string]string {
	dbMap := make(map[string]string)
	for _, urlRecord := range urlRecords {
		dbMap[urlRecord.ShortURL] = urlRecord.OriginalURL
	}
	return dbMap
}

// CreateRecordAndUpdateDBMap создает запись в памяти и вписывает её в мапу.
func CreateRecordAndUpdateDBMap(dbMap map[string]string, originalURL string, shortURL string, id string, userID string) model.URLRecord {
	newRecord := model.URLRecord{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		UUID:        id,
		UserID:      userID,
	}

	dbMap[shortURL] = originalURL
	return newRecord
}

// GetUrlsCreatedByUser возвращает URL которые создал текущий пользователь.
func GetUrlsCreatedByUser(urlRecords []model.URLRecord, userID string, hostAndPathPart string, req *http.Request) []model.SimpleURLRecord {
	var usersUrls []model.SimpleURLRecord
	for _, urlRecord := range urlRecords {
		if urlRecord.UserID == userID && !urlRecord.DeletedFlag {
			usersURL := model.SimpleURLRecord{
				ShortURL:    urltf.MakeResultShortenedURL(hostAndPathPart, urlRecord.ShortURL, req),
				OriginalURL: urlRecord.OriginalURL,
			}
			usersUrls = append(usersUrls, usersURL)
		}
	}
	return usersUrls
}
