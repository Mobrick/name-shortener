package goconst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzer(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "positive new goconst analyser #1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, NewAnalyzer())
		})
	}
}
