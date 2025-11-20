package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Collectors CollectorsConfig `yaml:"collectors"`
	AI         AIConfig         `yaml:"ai"`
	Logging    LoggingConfig    `yaml:"logging"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Port        int    `yaml:"port"`
	MetricsPath string `yaml:"metrics_path"`
	APIPath     string `yaml:"api_path"`
}

// CollectorsConfig represents collectors configuration
type CollectorsConfig struct {
	FileLog FileLogConfig `yaml:"file_log"`
	HTTPOFD HTTPOFDConfig `yaml:"http_ofd"`
}

// FileLogConfig represents file log collector configuration
type FileLogConfig struct {
	Enabled      bool          `yaml:"enabled"`
	Path         string        `yaml:"path"`
	Format       string        `yaml:"format"`
	PollInterval time.Duration `yaml:"poll_interval"`
}

// HTTPOFDConfig represents HTTP OFD collector configuration
type HTTPOFDConfig struct {
	Enabled      bool          `yaml:"enabled"`
	URL          string        `yaml:"url"`
	APIKey       string        `yaml:"api_key"`
	PollInterval time.Duration `yaml:"poll_interval"`
	Timeout      time.Duration `yaml:"timeout"`
}

// AIConfig represents AI subsystem configuration
type AIConfig struct {
	Provider         string                  `yaml:"provider"`
	ErrorClustering  ErrorClusteringConfig   `yaml:"error_clustering"`
	AlertAdvisor     AlertAdvisorConfig      `yaml:"alert_advisor"`
}

// ErrorClusteringConfig represents error clustering configuration
type ErrorClusteringConfig struct {
	Enabled         bool    `yaml:"enabled"`
	MinClusterSize  int     `yaml:"min_cluster_size"`
	SimilarityThreshold float64 `yaml:"similarity_threshold"`
}

// AlertAdvisorConfig represents alert advisor configuration
type AlertAdvisorConfig struct {
	Enabled bool `yaml:"enabled"`
	LookbackPeriod time.Duration `yaml:"lookback_period"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Load loads configuration from file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Expand environment variables
	data = []byte(os.ExpandEnv(string(data)))

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.MetricsPath == "" {
		c.Server.MetricsPath = "/metrics"
	}

	if c.Server.APIPath == "" {
		c.Server.APIPath = "/api/v1"
	}

	if c.Collectors.FileLog.Enabled {
		if c.Collectors.FileLog.Path == "" {
			return fmt.Errorf("file_log path is required when enabled")
		}
		if c.Collectors.FileLog.Format == "" {
			c.Collectors.FileLog.Format = "json"
		}
		if c.Collectors.FileLog.PollInterval == 0 {
			c.Collectors.FileLog.PollInterval = 10 * time.Second
		}
	}

	if c.Collectors.HTTPOFD.Enabled {
		if c.Collectors.HTTPOFD.URL == "" {
			return fmt.Errorf("http_ofd url is required when enabled")
		}
		if c.Collectors.HTTPOFD.PollInterval == 0 {
			c.Collectors.HTTPOFD.PollInterval = 30 * time.Second
		}
		if c.Collectors.HTTPOFD.Timeout == 0 {
			c.Collectors.HTTPOFD.Timeout = 10 * time.Second
		}
	}

	if c.AI.Provider == "" {
		c.AI.Provider = "mock"
	}

	if c.AI.ErrorClustering.MinClusterSize == 0 {
		c.AI.ErrorClustering.MinClusterSize = 5
	}

	if c.AI.ErrorClustering.SimilarityThreshold == 0 {
		c.AI.ErrorClustering.SimilarityThreshold = 0.7
	}

	if c.AI.AlertAdvisor.LookbackPeriod == 0 {
		c.AI.AlertAdvisor.LookbackPeriod = 7 * 24 * time.Hour // 7 days
	}

	if c.Logging.Level == "" {
		c.Logging.Level = "info"
	}

	if c.Logging.Format == "" {
		c.Logging.Format = "json"
	}

	return nil
}
