package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/Mobrick/name-shortener/pkg/urltf"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/sync/errgroup"
)

// PostgreDB реализация для постгре.
type PostgreDB struct {
	DatabaseConnection *sql.DB
	DatabaseMap        map[string]string
}

// PingDB пингует подключение к бд.
func (dbData PostgreDB) PingDB() error {
	err := dbData.DatabaseConnection.Ping()
	return err
}

// AddMany добавляет множество данных о сокращенных URL в хранилище.
func (dbData PostgreDB) AddMany(ctx context.Context, shortURLRequestMap map[string]model.BatchRequestURL, userID string) error {
	var sliceOfRecords []model.URLRecord
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userID)
		sliceOfRecords = append(sliceOfRecords, newRecord)
	}

	tx, err := dbData.DatabaseConnection.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt := "INSERT INTO url_records (uuid, short_url, original_url, user_id)" +
		" VALUES ($1, $2, $3, $4)"

	for _, record := range sliceOfRecords {
		_, err := dbData.DatabaseConnection.ExecContext(ctx, stmt, record.UUID, record.ShortURL, record.OriginalURL, record.UserID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Add добавляет данные о сокращенном URL в хранилище.
func (dbData PostgreDB) Add(ctx context.Context, shortURL string, originalURL string, userID string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userID)

	err := dbData.createURLRecordsTableIfNotExists(ctx)
	if err != nil {
		return "", err
	}

	insertStmt := "INSERT INTO url_records (uuid, short_url, original_url, user_id)" +
		" VALUES ($1, $2, $3, $4)"

	_, err = dbData.DatabaseConnection.ExecContext(ctx, insertStmt, newRecord.UUID, newRecord.ShortURL, originalURL, newRecord.UserID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("url %s already in database", originalURL)
			exisitingURL, findErr := dbData.findExisitingShortURL(ctx, originalURL)
			if findErr != nil {
				return "", findErr
			}
			return exisitingURL, nil
		}
		log.Printf("Failed to insert a record: " + originalURL)
		return "", err
	}

	return "", nil
}

// Get возвращает оригинальный URL, либо сообщает об отсуствии соответсвующего URL, также возвращает пометку об удалении.
func (dbData PostgreDB) Get(ctx context.Context, shortURL string) (string, bool, error) {
	var location string
	var isDeleted bool

	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM url_records WHERE short_url = $1", shortURL)
	err := row.Scan(&location, &isDeleted)
	if err == sql.ErrNoRows {
		return "", isDeleted, nil
	} else if err != nil {
		return "", isDeleted, err
	}
	return location, isDeleted, nil
}

func (dbData PostgreDB) createURLRecordsTableIfNotExists(ctx context.Context) error {
	_, err := dbData.DatabaseConnection.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS "+urlRecordsTableName+
			` (uuid TEXT PRIMARY KEY, 
			short_url TEXT NOT NULL, 
			original_url TEXT NOT NULL UNIQUE,
			user_id TEXT NOT NULL, 
			is_deleted BOOLEAN DEFAULT FALSE)`)

	if err != nil {
		return err
	}
	return nil
}

func (dbData PostgreDB) findExisitingShortURL(ctx context.Context, originalURL string) (string, error) {
	stmt := "SELECT short_url FROM url_records WHERE original_url = $1"
	var shortURL string
	err := dbData.DatabaseConnection.QueryRowContext(ctx, stmt, originalURL).Scan(&shortURL)
	if err != nil {
		log.Printf("Failed to find existing shotened url for this: " + originalURL)
		return "", err
	}
	return shortURL, nil
}

// Close закрывает подключение к хранилищу.
func (dbData PostgreDB) Close() {
	dbData.DatabaseConnection.Close()
}

// GetUrlsByUserID возвращает записи созданные пользователем.
func (dbData PostgreDB) GetUrlsByUserID(ctx context.Context, userID string, hostAndPathPart string, req *http.Request) ([]model.SimpleURLRecord, error) {
	var usersUrls []model.SimpleURLRecord
	stmt := "SELECT short_url, original_url FROM url_records WHERE user_id = $1 AND is_deleted = false"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var shortURL, originalURL string
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}

		usersURL := model.SimpleURLRecord{
			ShortURL:    urltf.MakeResultShortenedURL(hostAndPathPart, shortURL, req),
			OriginalURL: originalURL,
		}
		usersUrls = append(usersUrls, usersURL)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	defer rows.Close()
	return usersUrls, nil
}

// Delete удаляет данные о сокращенном URL из хранилища. При запросе пользователя удаление происходит не сразу, и просто ставится пометка об удалении.
func (dbData PostgreDB) Delete(ctx context.Context, urlsToDelete []string, userID string) error {
	g := new(errgroup.Group)
	g.Go(func() error {
		err := deletionRecepient(ctx, dbData.DatabaseConnection, urlsToDelete, userID)
		if err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func deletionRecepient(ctx context.Context, dbConnection *sql.DB, urlsToDelete []string, userID string) error {
	stmt := "UPDATE url_records SET is_deleted = true WHERE short_url = ANY($1) AND user_id = $2"
	_, err := dbConnection.ExecContext(ctx, stmt, urlsToDelete, userID)

	if err != nil {
		return err
	}
	return nil
}

// GetStats возвращает число юзеров и урлов
func (dbData PostgreDB) GetStats(ctx context.Context) (model.Stats, error) {
	db := dbData.DatabaseConnection

	var rowCount int
	var uniqueUserCount int

	err := db.QueryRow("SELECT COUNT(*) FROM url_records").Scan(&rowCount)
	if err != nil {
		return model.Stats{}, err
	}

	err = db.QueryRow("SELECT COUNT(DISTINCT user_id) FROM url_records").Scan(&uniqueUserCount)
	if err != nil {
		return model.Stats{}, err
	}
	
	return model.Stats{Urls: rowCount, Users: uniqueUserCount}, nil
}
