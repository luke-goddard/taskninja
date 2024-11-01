package config

import (
	"os"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type LogLevel string

const (
	LogLevelTrace LogLevel = "trace"
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

const DefaultLogPath = "/tmp/taskninja.log"

type LogMode string

const (
	LogModePretty LogMode = "pretty"
	LogModeJson   LogMode = "json"
)

type Log struct {
	Level string `yaml:"level"` // log level
	Mode  string `yaml:"mode"`  // log mode
	Path  string `yaml:"path"`  // log path
}

func (c *Config) InitLogger() {
	switch LogMode(c.Log.Mode) {
	case LogModePretty:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case LogModeJson:
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	var level zerolog.Level
	switch LogLevel(c.Log.Level) {
	case LogLevelTrace:
		level = zerolog.TraceLevel
	case LogLevelDebug:
		level = zerolog.DebugLevel
		log.Logger = log.With().Caller().Logger()
	case LogLevelInfo:
		level = zerolog.InfoLevel
	case LogLevelWarn:
		level = zerolog.WarnLevel
	case LogLevelError:
		level = zerolog.ErrorLevel
	default:
		log.Warn().Msg("Unknown log level set in config file, defaulting to info")
		level = zerolog.InfoLevel
	}

	var file, err = os.OpenFile(c.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	assert.NotNil(file, "Failed to open log file")

	err = file.Truncate(0)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to truncate log file")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(level)
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}
