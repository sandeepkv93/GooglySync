package config

import (
	"errors"
	"os"
	"path/filepath"
)

const appDirName = "drive-client"

// Config holds basic runtime configuration.
type Config struct {
	AppName      string
	ConfigDir    string
	DataDir      string
	LogLevel     string
	DatabasePath string
}

// NewConfig builds a default config from XDG paths and environment.
func NewConfig() (*Config, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		var err error
		configHome, err = os.UserConfigDir()
		if err != nil {
			return nil, err
		}
	}

	dataHome := os.Getenv("XDG_DATA_HOME")
	if dataHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dataHome = filepath.Join(home, ".local", "share")
	}

	if configHome == "" || dataHome == "" {
		return nil, errors.New("unable to resolve XDG directories")
	}

	configDir := filepath.Join(configHome, appDirName)
	dataDir := filepath.Join(dataHome, appDirName)

	return &Config{
		AppName:      "googlysync",
		ConfigDir:    configDir,
		DataDir:      dataDir,
		LogLevel:     "info",
		DatabasePath: filepath.Join(dataDir, "drive.db"),
	}, nil
}
