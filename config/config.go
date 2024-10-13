package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

// ConnectionMode is the type of connection to the sqlite database
type ConnectionMode string

const (
	ConnectionModeInMemory ConnectionMode = "in-memory"
	ConnectionModeFile     ConnectionMode = "file"
	ConnectionModeHTTP     ConnectionMode = "http"
)

type Connection struct {
	Mode ConnectionMode `yaml:"mode"`       // in-memory, file, http
	Path string         `yaml:"connection"` // connection string
}

type Config struct {
	Connection Connection `yaml:"connection"` // How to connect to the sqlite database
}

type ConfigError struct {
	Err error
}

func (e ConfigError) Error() string {
	return fmt.Sprintf("Failed to load config file: %v", e.Err)
}

// Used to load the configuration from the config file
func GetConfig() (*Config, error) {
	var home, err = os.UserHomeDir()
	if err != nil {
		var err = fmt.Errorf("Failed to get user home directory: %v", err)
		return nil, ConfigError{Err: err}
	}
	var configLocation = path.Join(home, ".config", "taskninja", "config.yaml")

	if _, err = os.Stat(configLocation); os.IsNotExist(err) {
		var err = fmt.Errorf("Config file not found: %v", err)
		return nil, ConfigError{Err: err}
	}

	viper.SetConfigFile(configLocation)
	if err = viper.ReadInConfig(); err != nil {
		var err = fmt.Errorf("Failed to read config file: %v", err)
		return nil, ConfigError{Err: err}
	}
	viper.SetDefault("connection.mode", ConnectionModeInMemory)
	viper.SetDefault("connection.path", "")

	var config = Config{}
	if err = viper.Unmarshal(&config); err != nil {
		var err = fmt.Errorf("Failed to unmarshal config: %v", err)
		return nil, ConfigError{Err: err}
	}
	return &config, nil
}
