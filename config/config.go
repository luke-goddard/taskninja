package config

import (
	"fmt"
	"os"
	"path"

	"github.com/luke-goddard/taskninja/assert"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// ConnectionMode is the type of connection to the sqlite database
type ConnectionMode string

const (
	ConnectionModeInMemory ConnectionMode = "in-memory"
	ConnectionModeFile     ConnectionMode = "file"
	ConnectionModeHTTP     ConnectionMode = "http"
)

type SqlConnectionConfig struct {
	Mode       ConnectionMode `yaml:"mode"`       // in-memory, file, http
	Path       string         `yaml:"connection"` // connection string
	BackupPath string         `yaml:"backupPath"` // Location to backup the database (sqlite disk only)
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
	assert.Nil(err, "Failed to get user home directory, cannot load/create config")
	var configDir = path.Join(home, ".config", "taskninja")

	if _, err = os.Stat(configDir); os.IsNotExist(err) {
		log.Info().Str("directory", configDir).Msg("Creating initial config directory")
		if err = os.MkdirAll(configDir, 0750); err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to create config directory")
		}
	}
	var configLocation = path.Join(home, ".config", "taskninja", "config.yaml")
	log.Info().Str("file", configLocation).Msg("Creating initial config file")
	_, err = os.Create(configLocation) // #nosec G304
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
			Msg("Failed to write initial config file")
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
