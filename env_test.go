package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fellippemendonca/genv"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	Port             int64   `name:"GO_ENV_PORT" required:"true" default:"3000"`
	SentryDsn        string  `name:"GO_ENV_SENTRY_DSN" required:"false"`
	SentryEnv        string  `name:"GO_ENV_SENTRY_ENV" required:"true"`
	SentrySampleRate float64 `name:"GO_ENV_SENTRY_SAMPLE_RATE" default:"0.01"`
}

func TestConfigLoad(t *testing.T) {
	expectedValues := &Config{
		Port:             3000,
		SentryDsn:        "https://123@o456.ingest.sentry.io/789",
		SentryEnv:        "development",
		SentrySampleRate: 1,
	}
	setEnv(expectedValues)
	defer unsetEnv()
	loadedConfig := &Config{}
	assert.NoError(t, env.Load(loadedConfig))
	assert.Equal(t, expectedValues, loadedConfig)
}

func TestConfigLoad_defaults(t *testing.T) {
	expectedValues := &Config{
		Port:             3000,
		SentryDsn:        "",
		SentryEnv:        "development",
		SentrySampleRate: 0.01,
	}
	setEnv(expectedValues)
	os.Unsetenv("GO_ENV_PORT")
	os.Unsetenv("GO_ENV_SENTRY_SAMPLE_RATE")
	defer unsetEnv()
	loadedConfig := &Config{}
	assert.NoError(t, env.Load(loadedConfig))
	assert.Equal(t, expectedValues, loadedConfig)
}

func TestConfigLoad_required(t *testing.T) {
	envValues := &Config{}
	setEnv(envValues)
	defer unsetEnv()
	loadedConfig := &Config{}
	assert.EqualError(t, env.Load(loadedConfig), "required env var not set: GO_ENV_SENTRY_ENV")
}

func setEnv(cfg *Config) {
	os.Setenv("GO_ENV_PORT", fmt.Sprint(cfg.Port))
	os.Setenv("GO_ENV_SENTRY_DSN", cfg.SentryDsn)
	os.Setenv("GO_ENV_SENTRY_ENV", cfg.SentryEnv)
	os.Setenv("GO_ENV_SENTRY_SAMPLE_RATE", fmt.Sprint(cfg.SentrySampleRate))
}

func unsetEnv() {
	os.Unsetenv("GO_ENV_PORT")
	os.Unsetenv("GO_ENV_SENTRY_DSN")
	os.Unsetenv("GO_ENV_SENTRY_ENV")
	os.Unsetenv("GO_ENV_SENTRY_SAMPLE_RATE")
}
