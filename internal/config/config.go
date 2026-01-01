package config

import (
	"encoding/json"
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
	ConfigFile   string
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

// Options defines runtime overrides for config resolution.
type Options struct {
	ConfigPath string
	LogLevel   string
}

type fileConfig struct {
	AppName      string `json:"app_name"`
	ConfigDir    string `json:"config_dir"`
	DataDir      string `json:"data_dir"`
	LogLevel     string `json:"log_level"`
	DatabasePath string `json:"database_path"`
}

// NewConfigWithOptions resolves config and applies overrides from options.
func NewConfigWithOptions(opts Options) (*Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if opts.ConfigPath != "" {
		if err := applyConfigFile(cfg, opts.ConfigPath); err != nil {
			return nil, err
		}
		cfg.ConfigFile = opts.ConfigPath
	}

	if opts.LogLevel != "" {
		cfg.LogLevel = opts.LogLevel
	}

	return cfg, nil
}

func applyConfigFile(cfg *Config, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var fc fileConfig
	if err := json.Unmarshal(data, &fc); err != nil {
		return err
	}

	if fc.AppName != "" {
		cfg.AppName = fc.AppName
	}
	if fc.ConfigDir != "" {
		cfg.ConfigDir = fc.ConfigDir
	}
	if fc.DataDir != "" {
		cfg.DataDir = fc.DataDir
	}
	if fc.LogLevel != "" {
		cfg.LogLevel = fc.LogLevel
	}
	if fc.DatabasePath != "" {
		cfg.DatabasePath = fc.DatabasePath
	}

	return nil
}
