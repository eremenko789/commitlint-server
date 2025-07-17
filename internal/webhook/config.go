package webhook

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Config represents the webhook server configuration
type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Gitea       GiteaConfig       `yaml:"gitea"`
	Commitlint  CommitlintConfig  `yaml:"commitlint"`
	Webhook     WebhookConfig     `yaml:"webhook"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Address      string `yaml:"address"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// GiteaConfig contains Gitea API settings
type GiteaConfig struct {
	BaseURL     string `yaml:"base_url"`
	Token       string `yaml:"token"`
	Username    string `yaml:"username"`
}

// CommitlintConfig contains commitlint settings
type CommitlintConfig struct {
	ConfigPath string `yaml:"config_path"`
}

// WebhookConfig contains webhook-specific settings
type WebhookConfig struct {
	Secret string   `yaml:"secret"`
	Path   string   `yaml:"path"`
	Events []string `yaml:"events"`
}

// LoadConfig loads configuration from file or environment
func LoadConfig() (*Config, error) {
	// Default configuration
	config := &Config{
		Server: ServerConfig{
			Address:      ":8080",
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Webhook: WebhookConfig{
			Path: "/webhook",
			Events: []string{
				"pull_request",
				"pull_request_sync",
			},
		},
		Commitlint: CommitlintConfig{
			ConfigPath: ".commitlintrc.yml",
		},
	}

	// Try to load from config file
	configPaths := []string{
		"webhook-server.yml",
		"webhook-server.yaml",
		"/etc/commitlint/webhook-server.yml",
		filepath.Join(os.Getenv("HOME"), ".config/commitlint/webhook-server.yml"),
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if addr := os.Getenv("WEBHOOK_SERVER_ADDRESS"); addr != "" {
		config.Server.Address = addr
	}
	if url := os.Getenv("GITEA_BASE_URL"); url != "" {
		config.Gitea.BaseURL = url
	}
	if token := os.Getenv("GITEA_TOKEN"); token != "" {
		config.Gitea.Token = token
	}
	if username := os.Getenv("GITEA_USERNAME"); username != "" {
		config.Gitea.Username = username
	}
	if secret := os.Getenv("WEBHOOK_SECRET"); secret != "" {
		config.Webhook.Secret = secret
	}

	// Validate configuration
	if config.Gitea.BaseURL == "" {
		return nil, fmt.Errorf("gitea base URL is required")
	}
	if config.Gitea.Token == "" {
		return nil, fmt.Errorf("gitea token is required")
	}

	return config, nil
}