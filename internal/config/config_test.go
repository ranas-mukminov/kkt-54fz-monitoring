package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
server:
  port: 9090
  metrics_path: /metrics

collectors:
  file_log:
    enabled: true
    path: /var/log/kkt/*.log
    format: json
    poll_interval: 10s
  
  http_ofd:
    enabled: true
    url: https://ofd.example.ru/api/v1
    api_key: test-key
    poll_interval: 30s
    timeout: 10s

ai:
  provider: mock
  error_clustering:
    enabled: true
    min_cluster_size: 5
  alert_advisor:
    enabled: true

logging:
  level: info
  format: json
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test server config
	if cfg.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", cfg.Server.Port)
	}
	if cfg.Server.MetricsPath != "/metrics" {
		t.Errorf("Expected metrics_path /metrics, got %s", cfg.Server.MetricsPath)
	}

	// Test collectors config
	if !cfg.Collectors.FileLog.Enabled {
		t.Error("Expected file_log to be enabled")
	}
	if cfg.Collectors.FileLog.PollInterval != 10*time.Second {
		t.Errorf("Expected poll_interval 10s, got %v", cfg.Collectors.FileLog.PollInterval)
	}

	// Test AI config
	if cfg.AI.Provider != "mock" {
		t.Errorf("Expected AI provider mock, got %s", cfg.AI.Provider)
	}

	// Test logging config
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected log level info, got %s", cfg.Logging.Level)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Server: ServerConfig{
					Port: 9090,
				},
				Collectors: CollectorsConfig{
					FileLog: FileLogConfig{
						Enabled: true,
						Path:    "/var/log/kkt/*.log",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			cfg: Config{
				Server: ServerConfig{
					Port: -1,
				},
			},
			wantErr: true,
		},
		{
			name: "file_log enabled without path",
			cfg: Config{
				Server: ServerConfig{
					Port: 9090,
				},
				Collectors: CollectorsConfig{
					FileLog: FileLogConfig{
						Enabled: true,
						Path:    "",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
