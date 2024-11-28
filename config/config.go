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
	ConnectionModeInMemory ConnectionMode = "in-memory" // SQLITE in-memory database
	ConnectionModeFile     ConnectionMode = "file"      // SQLITE file database
)

// Contains the SQL connection configuration
type SqlConnectionConfig struct {
	Mode       ConnectionMode `yaml:"mode"`       // in-memory, file, http
	Path       string         `yaml:"connection"` // connection string
	BackupPath string         `yaml:"backupPath"` // Location to backup the database (sqlite disk only)
}

// DSN returns the data source name for the connection e.g "sqlite://:memory:",
// or "sqlite:///path/to/file.db", Defaults to in-memory if not set
func (c *SqlConnectionConfig) DSN() string {
	switch c.Mode {
	case ConnectionModeInMemory:
		return ":memory:"
	case ConnectionModeFile:
		return c.Path
	default:
		return ":memory:"
	}
}

// Contains the user configuration
type Config struct {
	Connection SqlConnectionConfig `yaml:"connection"` // How to connect to the sqlite database
	Log        Log                 `yaml:"log"`        // How to log
}

type ConfigErrorVariant string

const (
	ConfigErrorNoHomeDir             ConfigErrorVariant = "no-home-dir"               // User home directory not found
	ConfigErrorFileNotFound          ConfigErrorVariant = "file-not-found"            // Config file not found
	ConfigErrorConfigDirDoesNotExist ConfigErrorVariant = "config-dir-does-not-exist" // Config directory does not exist
	ConfigErrorReadFile              ConfigErrorVariant = "read-file"                 // Failed to read config file
	ConfigErrorUnmarshal             ConfigErrorVariant = "unmarshal"                 // Failed to unmarshal config
)

// ConfigError is an error that occurs when loading the configuration
type ConfigError struct {
	Variant ConfigErrorVariant // The type of error
	Err     error              // The error that occurred
}

// CanBootstrap returns true if the error can be resolved by bootstrapping the configuration
func (e ConfigError) CanBootstrap() bool {
	return e.Variant == ConfigErrorFileNotFound || e.Variant == ConfigErrorConfigDirDoesNotExist
}

// Error returns the error message
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

// Bootstrap creates the initial configuration file if it does not exist
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

func setDefaults() {
	viper.SetDefault("connection.mode", ConnectionModeInMemory)
	viper.SetDefault("connection.path", "")
	viper.SetDefault("log.level", LogLevelInfo)
	viper.SetDefault("log.mode", LogModePretty)
	viper.SetDefault("log.path", DefaultLogPath)
}
