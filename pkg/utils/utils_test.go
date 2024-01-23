package utils

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeys(t *testing.T) {
	tests := []struct {
		name  string
		dbMap map[string]string
		want  []string
	}{
		{
			name: "positive GetKeys test #1",
			dbMap: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
			},
			want: []string{
				"key1",
				"key2",
				"key3",
				"key4",
				"key5",
			},
		},
		{
			name: "positive GetKeys test #2",
			dbMap: map[string]string{
				"key2": "value3",
				"key3": "value3",
				"key5": "value3",
				"key4": "value3",
			},
			want: []string{
				"key3",
				"key2",
				"key5",
				"key4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := GetKeys(tt.dbMap)
			wantSlice := tt.want
			slices.Sort(keys)
			slices.Sort(wantSlice)
			assert.Equal(t, wantSlice, keys)
		})
	}
}
