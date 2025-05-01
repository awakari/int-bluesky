package config

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	os.Setenv("API_HTTP_HOST", "bluesky.awakari.com")
	os.Setenv("LOG_LEVEL", "4")
	os.Setenv("API_TOKEN_INTERNAL", "foo")
	os.Setenv("API_BLUESKY_APP_PASSWORD", "bar")
	cfg, err := NewConfigFromEnv()
	assert.Nil(t, err)
	assert.Equal(t, slog.LevelWarn, slog.Level(cfg.Log.Level))
}
