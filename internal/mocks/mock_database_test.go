package mocks

import (
	"context"
	"reflect"
	"testing"

	"github.com/Mobrick/name-shortener/internal/database"
	"github.com/Mobrick/name-shortener/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMockDB(t *testing.T) {
	tests := []struct {
		name string
		want reflect.Type
	}{
		{
			name: "positive new mock db test #1",
			want: reflect.TypeOf(MockDB{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMockDB()
			assert.IsType(t, tt.want, reflect.TypeOf(got))
		})
	}
}

func TestMockDB_PingDB(t *testing.T) {
	tests := []struct {
		name    string
		dbData  database.Storage
		wantErr bool
	}{
		{
			name:    "negative ping test #1",
			dbData:  NewMockDB(),
			wantErr: true,
		}, {
			name:    "negative ping test #2",
			dbData:  &MockDB{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dbData.PingDB()
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestMockDB_GetUrlsByUserID(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.SimpleURLRecord
		wantErr bool
		dbData  MockDB
	}{
		{
			name:   "positive get urls by user test #1",
			dbData: MockDB{},
			args: args{
				userID: "1a91a181-80ec-45cb-a576-14db11505700",
			},
			want: []models.SimpleURLRecord{
				{
					ShortURL:    "DDDDdddd",
					OriginalURL: "https://www.google.com/",
				},
				{
					ShortURL:    "vvvv4444",
					OriginalURL: "https://www.go.com/",
				},
			},
			wantErr: false,
		},
		{
			name: "positive get urls by user test #1",
			args: args{
				userID: "1u",
			},
			want:    []models.SimpleURLRecord{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dbData.GetUrlsByUserID(context.Background(), tt.args.userID, "", nil)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, len(tt.want), len(got))
			for _, record := range got {
				assert.Contains(t, tt.want, record)
			}
		})
	}
}
