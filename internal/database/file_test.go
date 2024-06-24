package database

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Mobrick/name-shortener/internal/filestorage"
	"github.com/Mobrick/name-shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileDB_PingDB(t *testing.T) {
	tests := []struct {
		name    string
		dbData  Storage
		wantErr bool
	}{
		{
			name:    "negative ping test #1",
			dbData:  NewDB("tmp/test.txt", ""),
			wantErr: true,
		}, {
			name:    "negative ping test #2",
			dbData:  &FileDB{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dbData.PingDB()
			assert.Equal(t, tt.wantErr, err != nil)

			os.Remove("tmp/test.txt")
			tt.dbData.Close()
		})
	}
}

func TestFileDB_Get(t *testing.T) {
	type args struct {
		ctx      context.Context
		shortURL string
	}
	tests := []struct {
		name    string
		dbData  FileDB
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			name: "positive get in memory #1",
			dbData: FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gg":    "https://www.go.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			args: args{
				ctx:      context.Background(),
				shortURL: "gog",
			},
			want:    "https://www.google.com/",
			want1:   false,
			wantErr: false,
		}, {
			name: "positive get in memory #1",
			dbData: FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gg":    "https://www.go.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			args: args{
				ctx:      context.Background(),
				shortURL: "google",
			},
			want:    "",
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.dbData.Get(tt.args.ctx, tt.args.shortURL)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestFileDB_Add(t *testing.T) {
	type args struct {
		ctx         context.Context
		shortURL    string
		originalURL string
		userID      string
	}
	tests := []struct {
		name       string
		dbData     *FileDB
		args       args
		wantRecord model.URLRecord
	}{
		{
			name: "positive add test #1",
			dbData: &FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			args: args{
				ctx:         context.Background(),
				shortURL:    "gg",
				originalURL: "https://www.go.com/",
				userID:      "1u",
			},
			wantRecord: model.URLRecord{
				ShortURL:    "gg",
				OriginalURL: "https://www.go.com/",
				UserID:      "1u",
				DeletedFlag: false,
			},
		},
		{
			name: "positive add test #2",
			dbData: &FileDB{
				URLRecords:  []model.URLRecord{},
				DatabaseMap: map[string]string{},
			},
			args: args{
				ctx:         context.Background(),
				shortURL:    "gg",
				originalURL: "https://www.go.com/",
				userID:      "1u",
			},
			wantRecord: model.URLRecord{
				ShortURL:    "gg",
				OriginalURL: "https://www.go.com/",
				UserID:      "1u",
				DeletedFlag: false,
			},
		},
		{
			name: "positive add test #3",
			dbData: &FileDB{
				URLRecords:  []model.URLRecord{},
				DatabaseMap: map[string]string{},
			},
			args: args{
				ctx:         context.Background(),
				shortURL:    "gg",
				originalURL: "go.com/",
				userID:      "",
			},
			wantRecord: model.URLRecord{
				ShortURL:    "gg",
				OriginalURL: "go.com/",
				UserID:      "",
				DeletedFlag: false,
			},
		},
		{
			name: "positive add test #4",
			dbData: &FileDB{
				URLRecords:  []model.URLRecord{},
				DatabaseMap: map[string]string{},
			},
			args: args{
				ctx:         nil,
				shortURL:    "gg",
				originalURL: "go.com/",
			},
			wantRecord: model.URLRecord{
				ShortURL:    "gg",
				OriginalURL: "go.com/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := filestorage.MakeFile("tmp/test.txt")
			require.NoError(t, err)
			tt.dbData.FileStorage = file
			tt.dbData.Add(tt.args.ctx, tt.args.shortURL, tt.args.originalURL, tt.args.userID)
			assert.Contains(t, tt.dbData.DatabaseMap, tt.wantRecord.ShortURL)
			os.Remove("tmp/test.txt")
			tt.dbData.Close()
		})
	}
}

func TestFileDB_AddMany(t *testing.T) {
	type args struct {
		ctx                context.Context
		shortURLRequestMap map[string]model.BatchRequestURL
		userID             string
	}
	tests := []struct {
		name      string
		dbData    *FileDB
		args      args
		wantErr   bool
		wantCount int
	}{
		{
			name: "positive add many test #1",
			dbData: &FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			args: args{
				ctx: context.Background(),
				shortURLRequestMap: map[string]model.BatchRequestURL{
					"gg": {
						CorrelationID: "1",
						OriginalURL:   "https://www.go.com/",
					}, "ggg": {
						CorrelationID: "4",
						OriginalURL:   "https://www.ggo.com/",
					},
				},
			},
			wantErr:   false,
			wantCount: 4,
		}, {
			name: "positive add many test #1",
			dbData: &FileDB{
				URLRecords:  []model.URLRecord{},
				DatabaseMap: map[string]string{},
			},
			args: args{
				ctx: context.Background(),
				shortURLRequestMap: map[string]model.BatchRequestURL{
					"gg": {
						CorrelationID: "1",
						OriginalURL:   "https://www.go.com/",
					}, "ggg": {
						CorrelationID: "4",
						OriginalURL:   "https://www.ggo.com/",
					},
				},
			},
			wantErr:   false,
			wantCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := filestorage.MakeFile("tmp/test.txt")
			require.NoError(t, err)
			tt.dbData.FileStorage = file
			err = tt.dbData.AddMany(tt.args.ctx, tt.args.shortURLRequestMap, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantCount, len(tt.dbData.URLRecords))
			os.Remove("tmp/test.txt")
			tt.dbData.Close()
		})
	}
}

func TestFileDB_Close(t *testing.T) {
	tests := []struct {
		name   string
		dbData Storage
	}{
		{
			name:   "postitive close test #1",
			dbData: &FileDB{},
		},
		{
			name:   "postitive close test #2",
			dbData: NewDB("tmp/test.txt", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dbData.Close()
			os.Remove("tmp/test.txt")
		})
	}
}

func TestFileDB_GetUrlsByUserID(t *testing.T) {
	type args struct {
		urlRecords      []model.URLRecord
		userID          string
		hostAndPathPart string
		req             *http.Request
	}
	tests := []struct {
		name    string
		args    args
		want    []model.SimpleURLRecord
		wantErr bool
		dbData  FileDB
	}{
		{
			name: "positive get urls by user test #1",
			dbData: FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gg":    "https://www.go.com/",
					"gogru": "https://www.google.ru/",
				},
			},
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
			wantErr: false,
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
			want:    []model.SimpleURLRecord{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dbData.GetUrlsByUserID(context.Background(), tt.args.userID, tt.args.hostAndPathPart, tt.args.req)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, len(tt.want), len(got))
			for _, record := range got {
				assert.Contains(t, tt.want, record)
			}
		})
	}
}

func TestFileDB_Delete(t *testing.T) {
	type args struct {
		ctx          context.Context
		urlsToDelete []string
		userID       string
	}
	tests := []struct {
		name     string
		dbData   *FileDB
		args     args
		wantErr  bool
		wantDiff int
	}{
		{
			name: "postivie delete test #1",
			args: args{
				ctx:          context.Background(),
				urlsToDelete: []string{"gg", "gog"},
				userID:       "1u",
			},
			dbData: &FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gg":    "https://www.go.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			wantErr:  false,
			wantDiff: 2,
		},
		{
			name: "postivie delete test #2",
			args: args{
				ctx:          context.Background(),
				urlsToDelete: []string{"gg", "gogru"},
				userID:       "2u",
			},
			dbData: &FileDB{
				URLRecords: []model.URLRecord{
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
				DatabaseMap: map[string]string{
					"gog":   "https://www.google.com/",
					"gg":    "https://www.go.com/",
					"gogru": "https://www.google.ru/",
				},
			},
			wantErr:  false,
			wantDiff: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startLen := len(tt.dbData.URLRecords)
			err := tt.dbData.Delete(tt.args.ctx, tt.args.urlsToDelete, tt.args.userID)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantDiff, startLen-len(tt.dbData.URLRecords))
		})
	}
}
