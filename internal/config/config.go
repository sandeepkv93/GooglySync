package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
)

const appDirName = "drive-client"

// Config holds basic runtime configuration.
type Config struct {
	AppName            string
	ConfigDir          string
	DataDir            string
	LogLevel           string
	DatabasePath       string
	ConfigFile         string
	LogFilePath        string
	LogFileMaxMB       int
	LogFileMaxBackups  int
	LogFileMaxAgeDays  int
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
		AppName:           "googlysync",
		ConfigDir:         configDir,
		DataDir:           dataDir,
		LogLevel:          "info",
		DatabasePath:      filepath.Join(dataDir, "googlysync.db"),
		LogFilePath:       filepath.Join(dataDir, "logs", "daemon.jsonl"),
		LogFileMaxMB:      10,
		LogFileMaxBackups: 5,
		LogFileMaxAgeDays: 7,
	}, nil
}

// Options defines runtime overrides for config resolution.
type Options struct {
	ConfigPath string
	LogLevel   string
}

type fileConfig struct {
	AppName           string `json:"app_name"`
	ConfigDir         string `json:"config_dir"`
	DataDir           string `json:"data_dir"`
	LogLevel          string `json:"log_level"`
	DatabasePath      string `json:"database_path"`
	LogFilePath       string `json:"log_file_path"`
	LogFileMaxMB      int    `json:"log_file_max_mb"`
	LogFileMaxBackups int    `json:"log_file_max_backups"`
	LogFileMaxAgeDays int    `json:"log_file_max_age_days"`
}

// NewConfigWithOptions resolves config and applies overrides from options and environment.
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

	applyEnv(cfg)

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
	if fc.LogFilePath != "" {
		cfg.LogFilePath = fc.LogFilePath
	}
	if fc.LogFileMaxMB > 0 {
		cfg.LogFileMaxMB = fc.LogFileMaxMB
	}
	if fc.LogFileMaxBackups > 0 {
		cfg.LogFileMaxBackups = fc.LogFileMaxBackups
	}
	if fc.LogFileMaxAgeDays > 0 {
		cfg.LogFileMaxAgeDays = fc.LogFileMaxAgeDays
	}

	return nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("GOOGLYSYNC_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}
	if v := os.Getenv("GOOGLYSYNC_LOG_FILE"); v != "" {
		cfg.LogFilePath = v
	}
	if v := os.Getenv("GOOGLYSYNC_LOG_MAX_MB"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			cfg.LogFileMaxMB = i
		}
	}
	if v := os.Getenv("GOOGLYSYNC_LOG_MAX_BACKUPS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			cfg.LogFileMaxBackups = i
		}
	}
	if v := os.Getenv("GOOGLYSYNC_LOG_MAX_AGE_DAYS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			cfg.LogFileMaxAgeDays = i
		}
	}
}
