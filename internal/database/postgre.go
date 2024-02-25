package database

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Mobrick/name-shortener/internal/models"
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

func (dbData PostgreDB) AddMany(ctx context.Context, shortURLRequestMap map[string]models.BatchRequestURL) {
	var sliceOfRecords []models.URLRecord
	for shortURL, record := range shortURLRequestMap {
		newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, record.OriginalURL, shortURL, record.CorrelationID)
		sliceOfRecords = append(sliceOfRecords, newRecord)
	}

	tx, err := dbData.DatabaseConnection.Begin()
	if err != nil {
		log.Fatal(err)
	}

	defer tx.Rollback()

	stmt := "INSERT INTO url_records (uuid, short_url, original_url)" +
		" VALUES ($1, $2, $3)"

	for _, record := range sliceOfRecords {
		_, err := dbData.DatabaseConnection.ExecContext(ctx, stmt, record.UUID, record.ShortURL, record.OriginalURL)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func (dbData PostgreDB) Add(ctx context.Context, shortURL string, originalURL string) string {
	id := uuid.New().String()
	newRecord := CreateRecordAndUpdateDBMap(dbData.DatabaseMap, originalURL, shortURL, id)

	dbData.createURLRecordsTableIfNotExists(ctx)

	insertStmt := "INSERT INTO url_records (uuid, short_url, original_url)" +
		" VALUES ($1, $2, $3)"

	_, err := dbData.DatabaseConnection.ExecContext(ctx, insertStmt, newRecord.UUID, newRecord.ShortURL, originalURL)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			log.Printf("url %s already in database", originalURL)
			return dbData.findExisitingShortURL(ctx, originalURL)
		} else {
			log.Fatal("Failed to insert a record: "+originalURL, err)
		}
	}

	return ""
}

func (dbData PostgreDB) Get(ctx context.Context, shortURL string) (string, bool) {
	var location string

	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT original_url FROM url_records WHERE short_url = $1", shortURL)
	err := row.Scan(&location)
	if err == sql.ErrNoRows {
		return location, false
	}
	return location, true
}

func (dbData PostgreDB) createURLRecordsTableIfNotExists(ctx context.Context) {
	_, err := dbData.DatabaseConnection.ExecContext(ctx,
		"CREATE TABLE IF NOT EXISTS "+urlRecordsTableName+
			` (uuid TEXT PRIMARY KEY, 
			short_url TEXT NOT NULL, 
			original_url TEXT NOT NULL UNIQUE)`)

	if err != nil {
		log.Fatal(err)
	}
}

func (dbData PostgreDB) findExisitingShortURL(ctx context.Context, originalURL string) string {
	stmt := "SELECT short_url FROM url_records WHERE original_url = $1"
	var shortURL string
	err := dbData.DatabaseConnection.QueryRowContext(ctx, stmt, originalURL).Scan(&shortURL)
	if err != nil {
		log.Fatal("Failed to find existing shotened url for this: "+originalURL, err)
	}
	return shortURL
}

func (dbData PostgreDB) Close() {
	dbData.DatabaseConnection.Close()
}
