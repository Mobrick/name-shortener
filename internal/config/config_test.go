package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeConfig(t *testing.T) {
	type want struct {
		envRunAddr  string
		envBaseAddr string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive MakeConfig test #1",
			want: want{
				envRunAddr:  ":8081",
				envBaseAddr: "http://just.an.example.com/",
			},
		},
		{
			name: "positive MakeConfig test #2",
			want: want{
				envRunAddr:  ":8083",
				envBaseAddr: "http://ok.example.com/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SERVER_ADDRESS", tt.want.envRunAddr)
			os.Setenv("BASE_URL", tt.want.envBaseAddr)

			config := MakeConfig()

			resultRunAddr := config.FlagRunAddr
			resultBaseAddr := config.FlagShortURLBaseAddr

			assert.Equal(t, tt.want.envRunAddr, resultRunAddr)
			assert.Equal(t, tt.want.envBaseAddr, resultBaseAddr)

			os.Unsetenv("SERVER_ADDRESS")
			os.Unsetenv("BASE_URL")
		})
	}
}
