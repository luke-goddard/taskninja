package config

import (
	"fmt"
	"os"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConnectionMode is the type of connection to the sqlite database
type ConnectionMode string

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

const (
	ConnectionModeInMemory ConnectionMode = "in-memory"
	ConnectionModeFile     ConnectionMode = "file"
	ConnectionModeHTTP     ConnectionMode = "http"
)

type SqlConnectionConfig struct {
	Mode ConnectionMode `yaml:"mode"`       // in-memory, file, http
	Path string         `yaml:"connection"` // connection string
}

func (c *SqlConnectionConfig) DSN() string {
	switch c.Mode {
	case ConnectionModeInMemory:
		return ":memory:"
	case ConnectionModeFile:
		return c.Path
	case ConnectionModeHTTP:
		return c.Path
	default:
		return ":memory:"
	}
}

type Log struct {
	Level string `yaml:"level"` // log level
	Mode  string `yaml:"mode"`  // log mode
	Path  string `yaml:"path"`  // log path
}

type Config struct {
	Connection SqlConnectionConfig `yaml:"connection"` // How to connect to the sqlite database
	Log        Log                 `yaml:"log"`        // How to log
}

type ConfigErrorVariant string

const (
	ConfigErrorNoHomeDir             ConfigErrorVariant = "no-home-dir"
	ConfigErrorFileNotFound          ConfigErrorVariant = "file-not-found"
	ConfigErrorConfigDirDoesNotExist ConfigErrorVariant = "config-dir-does-not-exist"
	ConfigErrorReadFile              ConfigErrorVariant = "read-file"
	ConfigErrorUnmarshal             ConfigErrorVariant = "unmarshal"
)

type ConfigError struct {
	Variant ConfigErrorVariant
	Err     error
}

func (e ConfigError) CanBootstrap() bool {
	return e.Variant == ConfigErrorFileNotFound || e.Variant == ConfigErrorConfigDirDoesNotExist
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Failed to load config file: %v", e.Err)
}

// Used to load the configuration from the config file
func GetConfig() (*Config, *ConfigError) {
	var home, err = os.UserHomeDir()
	if err != nil {
		var err = fmt.Errorf("Failed to get user home directory: %v", err)
		return nil, &ConfigError{Err: err, Variant: ConfigErrorNoHomeDir}
	}
	var configDir = path.Join(home, ".config", "taskninja")
	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		var err = fmt.Errorf("Config directory does not exist: %v", err)
		return nil, &ConfigError{Err: err, Variant: ConfigErrorConfigDirDoesNotExist}
	}
	var configLocation = path.Join(home, ".config", "taskninja", "config.yaml")

	if _, err = os.Stat(configLocation); os.IsNotExist(err) {
		var err = fmt.Errorf("Config file not found: %v", err)
		return nil, &ConfigError{Err: err, Variant: ConfigErrorFileNotFound}
	}

	viper.SetConfigFile(configLocation)
	if err = viper.ReadInConfig(); err != nil {
		var err = fmt.Errorf("Failed to read config file: %v", err)
		return nil, &ConfigError{Err: err, Variant: ConfigErrorReadFile}
	}

	setDefaults()

	var config = Config{}
	if err = viper.Unmarshal(&config); err != nil {
		var err = fmt.Errorf("Failed to unmarshal config: %v", err)
		return nil, &ConfigError{Err: err, Variant: ConfigErrorUnmarshal}
	}
	return &config, nil
}

func setDefaults() {
	viper.SetDefault("connection.mode", ConnectionModeInMemory)
	viper.SetDefault("connection.path", "")
	viper.SetDefault("log.level", LogLevelInfo)
	viper.SetDefault("log.mode", LogModePretty)
	viper.SetDefault("log.path", DefaultLogPath)
}

func Bootstrap() *Config {
	var home, err = os.UserHomeDir()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to get user home directory, cannot bootstrap config")
	}
	var configDir = path.Join(home, ".config", "taskninja")

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		log.Info().Str("directory", configDir).Msg("Creating inital config directory")
		if err = os.MkdirAll(configDir, 0755); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to create config directory")
		}
	}
	var configLocation = path.Join(home, ".config", "taskninja", "config.yaml")
	log.Info().Str("file", configLocation).Msg("Creating inital config file")
	_, err = os.Create(configLocation)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to create config file")
	}

	viper.SetConfigFile(configLocation)

	err = viper.WriteConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to write intial config file")
	}
	var conf *Config
	var confErr *ConfigError
	conf, confErr = GetConfig()
	if confErr != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to load config after bootstrapping")
	}
	return conf
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

	var file, err = os.OpenFile(c.Log.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
	zerolog.SetGlobalLevel(level)
}
