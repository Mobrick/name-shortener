package filestorage

import (
	"io"
	"os"
	"testing"

	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeFile(t *testing.T) {
	tests := []struct {
		name           string
		fileName       string
		wantNotNilFile bool
		wantHasNoErr   bool
	}{
		{
			name:           "positive make file #1",
			fileName:       "tmp/test/storage",
			wantNotNilFile: true,
			wantHasNoErr:   true,
		},
		{
			name:           "positive make file #2",
			fileName:       "",
			wantNotNilFile: false,
			wantHasNoErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MakeFile(tt.fileName)
			assert.Equal(t, tt.wantHasNoErr, err == nil)
			assert.Equal(t, tt.wantNotNilFile, got != nil)
			os.Remove(tt.fileName)
		})
	}
}

func TestUploadNewURLRecord(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		record       models.URLRecord
		wantHasNoErr bool
	}{
		{
			name:         "positive upload file #1",
			fileName:     "tmp/test/storage",
			wantHasNoErr: true,
			record: models.URLRecord{
				UUID:        "1",
				ShortURL:    "gg",
				OriginalURL: "https://www.go.com/",
				UserID:      "1u",
				DeletedFlag: false,
			},
		},
		{
			name:         "positive upload file #2",
			fileName:     "",
			wantHasNoErr: true,
			record: models.URLRecord{
				UUID:        "1",
				ShortURL:    "gg",
				OriginalURL: "https://www.go.com/",
				UserID:      "1u",
				DeletedFlag: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := MakeFile(tt.fileName)
			require.NoError(t, err)
			err = UploadNewURLRecord(tt.record, file)
			assert.Equal(t, tt.wantHasNoErr, err == nil)
			os.Remove(tt.fileName)
		})
	}
}

func TestLoadURLRecords(t *testing.T) {
	tests := []struct {
		name             string
		fileName         string
		records          []models.URLRecord
		wantHasNoErr     bool
		wantRecordsCount int
	}{
		{
			name:             "positive load file #1",
			fileName:         "tmp/test/storage.json",
			wantHasNoErr:     true,
			wantRecordsCount: 2,
			records: []models.URLRecord{
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
			},
		},
		{
			name:             "positive load file #2",
			fileName:         "",
			wantHasNoErr:     true,
			wantRecordsCount: 0,
			records:          []models.URLRecord{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := MakeFile(tt.fileName)
			require.NoError(t, err)
			for _, record := range tt.records {
				err = UploadNewURLRecord(record, file)
				require.NoError(t, err)
			}
			file.Seek(0, io.SeekStart)
			records, err := LoadURLRecords(file)
			assert.Equal(t, tt.wantRecordsCount, len(records))
			assert.Equal(t, tt.wantHasNoErr, err == nil)
			file.Close()
			os.Remove(tt.fileName)
		})
	}
}
