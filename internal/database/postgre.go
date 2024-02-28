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

func (dbData PostgreDB) Get(ctx context.Context, shortURL string) (string, bool, error) {
	var location string

	// TODO: проверить, флаг is_deleted, если true, ответить 410 Gone
	row := dbData.DatabaseConnection.QueryRowContext(ctx, "SELECT original_url FROM url_records WHERE short_url = $1", shortURL)
	err := row.Scan(&location)
	if err == sql.ErrNoRows {
		return location, false, nil
	} else if err != nil {
		return location, false, err
	}
	return location, true, nil
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
	stmt := "UPDATE url_records SET is_deleted = true WHERE short_url = ANY($1)"
	_, err := dbData.DatabaseConnection.ExecContext(ctx, stmt, urlsToDelete)
	if err != nil {
		return err
	}
	return nil
}

// fanOut принимает канал данных, порождает 10 горутин
func fanOut(doneCh chan struct{}, inputCh chan int) []chan int {
    // количество горутин add
    numWorkers := 10
    // каналы, в которые отправляются результаты
    channels := make([]chan int, numWorkers)


    for i := 0; i < numWorkers; i++ {
        // получаем канал из горутины add
        addResultCh := add(doneCh, inputCh)
        // отправляем его в слайс каналов
        channels[i] = addResultCh
    }

    // возвращаем слайс каналов
    return channels
} 

// fanIn принимает несколько каналов, в которых итоговые значения
func fanIn(chs ...chan int) chan int {
    var wg sync.WaitGroup
    outCh := make(chan int)

    // определяем функцию output для каждого канала в chs
    // функция output копирует значения из канала с в канал outCh, пока с не будет закрыт 
    output := func(c chan int) {
        for n := range c {
            outCh <- n
        }
        wg.Done()
    }

    // добавляем в группу столько горутин, сколько каналов пришло в fanIn
    wg.Add(len(chs))
    // перебираем все каналы, которые пришли и отправляем каждый в отдельную горутину 
    for _, c := range chs {
        go output(c)
    }

    // запускаем горутину для закрытия outCh после того, как все горутины отработают 
    go func() {
        wg.Wait()
        close(outCh)
    }()

    // возвращаем общий канал
    return outCh
} 