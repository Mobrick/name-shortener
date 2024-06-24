package database

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/Mobrick/name-shortener/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestNewDB(t *testing.T) {
	tests := []struct {
		name             string
		fileName         string
		connectionString string
		wantType         reflect.Type
	}{
		{
			name:             "positive db storage creation test #1",
			fileName:         "",
			connectionString: "",
			wantType:         reflect.TypeOf(InMemoryDB{}),
		},
		{
			name:             "positive db storage creation test #2",
			fileName:         "tmp/test.txt",
			connectionString: "",
			wantType:         reflect.TypeOf(FileDB{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDB(tt.fileName, tt.connectionString)
			assert.IsType(t, tt.wantType, reflect.TypeOf(got))
		})
	}
}

func Test_newDBFromFile(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		wantNotNil bool
	}{
		{
			name:       "positive make file #1",
			fileName:   "tmp/test/storage.json",
			wantNotNil: true,
		},
		{
			name:       "positive make file #2",
			fileName:   "tmp/test/storage.txt",
			wantNotNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newDBFromFile(tt.fileName)
			assert.Equal(t, tt.wantNotNil, got != nil)
			got.Close()
		})
	}
}

func TestCreateRecordAndUpdateDBMap(t *testing.T) {
	type args struct {
		dbMap       map[string]string
		originalURL string
		shortURL    string
		id          string
		userID      string
	}
	tests := []struct {
		name string
		args args
		want model.URLRecord
	}{
		{
			name: "positive create record #1",
			args: args{
				dbMap:       make(map[string]string),
				id:          "1",
				shortURL:    "gg",
				originalURL: "https://www.go.com/",
				userID:      "1u",
			},
			want: model.URLRecord{
				UUID:        "1",
				ShortURL:    "gg",
				OriginalURL: "https://www.go.com/",
				UserID:      "1u",
			},
		},
		{
			name: "positive create record #1",
			args: args{
				dbMap:       make(map[string]string),
				id:          "2",
				shortURL:    "gog",
				originalURL: "https://www.google.com/",
				userID:      "1u",
			},
			want: model.URLRecord{
				UUID:        "2",
				ShortURL:    "gog",
				OriginalURL: "https://www.google.com/",
				UserID:      "1u",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRecord := CreateRecordAndUpdateDBMap(tt.args.dbMap, tt.args.originalURL, tt.args.shortURL, tt.args.id, tt.args.userID)
			assert.EqualValues(t, tt.want, gotRecord)
		})
	}
}

func TestGetUrlsCreatedByUser(t *testing.T) {
	type args struct {
		urlRecords      []model.URLRecord
		userID          string
		hostAndPathPart string
		req             *http.Request
	}
	tests := []struct {
		name string
		args args
		want []model.SimpleURLRecord
	}{
		{
			name: "positive get urls by user test #1",
			args: args{
				urlRecords: []model.URLRecord{
					{
						UUID:        "1",
						ShortURL:    "gg",
						OriginalURL: "https://www.go.com/",
						UserID:      "1u",
						DeletedFlag: false,
					},
					{
						UUID:        "2",
						ShortURL:    "gog",
						OriginalURL: "https://www.google.com/",
						UserID:      "1u",
						DeletedFlag: false,
					},
					{
						UUID:        "3",
						ShortURL:    "gogru",
						OriginalURL: "https://www.google.ru/",
						UserID:      "2u",
						DeletedFlag: false,
					},
				},
				userID:          "1u",
				hostAndPathPart: "http://shortener/",
				req:             httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.go.com/")),
			},
			want: []model.SimpleURLRecord{
				{
					ShortURL:    "http://shortener/gg",
					OriginalURL: "https://www.go.com/",
				},
				{
					ShortURL:    "http://shortener/gog",
					OriginalURL: "https://www.google.com/",
				},
			},
		},
		{
			name: "positive get urls by user test #1",
			args: args{
				urlRecords: []model.URLRecord{
					{
						UUID:        "1",
						ShortURL:    "gg",
						OriginalURL: "https://www.go.com/",
						UserID:      "2u",
						DeletedFlag: false,
					},
					{
						UUID:        "2",
						ShortURL:    "gog",
						OriginalURL: "https://www.google.com/",
						UserID:      "2u",
						DeletedFlag: false,
					},
					{
						UUID:        "3",
						ShortURL:    "gogru",
						OriginalURL: "https://www.google.ru/",
						UserID:      "2u",
						DeletedFlag: false,
					},
				},
				userID:          "1u",
				hostAndPathPart: "http://shortener/",
				req:             httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://www.go.com/")),
			},
			want: []model.SimpleURLRecord{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUrlsCreatedByUser(tt.args.urlRecords, tt.args.userID, tt.args.hostAndPathPart, tt.args.req)
			assert.Equal(t, len(tt.want), len(got))
			for _, record := range got {
				assert.Contains(t, tt.want, record)
			}
		})
	}
}
