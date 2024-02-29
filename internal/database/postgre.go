package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/Mobrick/name-shortener/urltf"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Реализация для постгре

type PostgreDB struct {
	DatabaseConnection *sql.DB
	DatabaseMap        map[string]string
}

func (dbData PostgreDB) PingDB() error {
	err := dbData.DatabaseConnection.Ping()
	return err
}

func (dbData PostgreDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL, userId string) error {
	var sliceOfRecords []models.URLRecord
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID, userId)
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

func (dbData PostgreDB) Add(ctx context.Context, shortURL string, originalURL string, userId string) (string, error) {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id, userId)

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
		} else {
			log.Printf("Failed to insert a record: " + originalURL)
			return "", err
		}
	}

	return "", nil
}

func (dbData PostgreDB) Get(ctx context.Context, shortURL string) (string, bool, bool, error) {
	var location string
	var isDeleted bool

	// TODO: проверить, флаг is_deleted, если true, ответить 410 Gone
	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM url_records WHERE short_url = $1", shortURL)
	err := row.Scan(&location, &isDeleted)
	if err == sql.ErrNoRows {
		return location, false, isDeleted, nil
	} else if err != nil {
		return location, false, isDeleted, err
	}
	return location, true, isDeleted, nil
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

func (dbData PostgreDB) Close() {
	dbData.DatabaseConnection.Close()
}

func (dbData PostgreDB) GetUrlsByUserId(ctx context.Context, userId string, hostAndPathPart string, req *http.Request) ([]models.SimpleURLRecord, error) {
	var usersUrls []models.SimpleURLRecord
	stmt := "SELECT short_url, original_url FROM url_records WHERE user_id = $1"
	rows, err := dbData.DatabaseConnection.QueryContext(ctx, stmt, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var shortURL, originalURL string
		err := rows.Scan(&shortURL, &originalURL)
		if err != nil {
			return nil, err
		}

		usersUrl := models.SimpleURLRecord{
			ShortURL:    urltf.MakeResultShortenedURL(hostAndPathPart, shortURL, req),
			OriginalURL: originalURL,
		}
		usersUrls = append(usersUrls, usersUrl)
	}

	defer rows.Close()
	return usersUrls, nil
}

func (dbData PostgreDB) Delete(ctx context.Context, urlsToDelete []string, userID string) error {
	var err error
	go func() {
		stmt := "UPDATE url_records SET is_deleted = true WHERE short_url = ANY($1) AND user_id = $2"
		_, err = dbData.DatabaseConnection.ExecContext(ctx, stmt, urlsToDelete, userID)
	}()

	if err != nil {
		return err
	}
	return nil
}
