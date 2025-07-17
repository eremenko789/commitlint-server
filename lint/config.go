package lint

// RuleSetting represent config for a rule
type RuleSetting struct {
	Argument interface{}            `yaml:"argument"`
	Flags    map[string]interface{} `yaml:"flags,omitempty"`
}

// SeverityConfig represent severity levels for rules
type SeverityConfig struct {
	Default Severity            `yaml:"default"`
	Rules   map[string]Severity `yaml:"rules,omitempty"`
}

// ServerConfig represent webhook server configuration
type ServerConfig struct {
	// Server listening address
	Address string `yaml:"address"`
	
	// Server listening port
	Port int `yaml:"port"`
	
	// Webhook secret for verifying Gitea requests
	WebhookSecret string `yaml:"webhook_secret"`
	
	// Gitea instance URL
	GiteaURL string `yaml:"gitea_url"`
	
	// Gitea access token for API calls
	GiteaToken string `yaml:"gitea_token"`
	
	// SSL certificate file path (optional)
	CertFile string `yaml:"cert_file,omitempty"`
	
	// SSL private key file path (optional)
	KeyFile string `yaml:"key_file,omitempty"`
	
	// Enable debug logging
	Debug bool `yaml:"debug"`
}

// Config represent linter config
type Config struct {
	// MinVersion is the minimum version of commitlint required
	// should be in semver format
	MinVersion string `yaml:"version"`

	// Formatter of the lint result
	Formatter string `yaml:"formatter"`

	// Enabled Rules
	Rules []string `yaml:"rules"`

	// Severity
	Severity SeverityConfig `yaml:"severity"`

	// Settings is rule name to rule settings
	Settings map[string]RuleSetting `yaml:"settings"`
	
	// Server configuration for webhook server
	Server ServerConfig `yaml:"server,omitempty"`
}

// GetRule returns RuleConfig for given rule name
func (c *Config) GetRule(ruleName string) RuleSetting {
	return c.Settings[ruleName]
}

// GetSeverity returns Severity for given ruleName
func (c *Config) GetSeverity(ruleName string) Severity {
	s, ok := c.Severity.Rules[ruleName]
	if ok {
		return s
	}
	return c.Severity.Default
}
