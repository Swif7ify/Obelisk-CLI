package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zalando/go-keyring"
)

// Config holds persistent user configuration.
type Config struct {
	// APIKey is explicitly excluded from JSON serialization to prevent saving to disk.
	APIKey       string `json:"-"`
	Model        string `json:"model,omitempty"`
	DefaultPath  string `json:"default_path,omitempty"`
	NoColor      bool   `json:"no_color,omitempty"`
	ReportFormat string `json:"report_format,omitempty"`
}

const configDirName = ".obelisk"
const configFileName = "config.json"
const keyringService = "obelisk-cli"
const keyringUser = "gemini-api-key"

// GetConfigDir returns the path to the Obelisk config directory.
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}

// GetConfigPath returns the full path to the config file.
func GetConfigPath() (string, error) {
	dir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

// Load reads the config from disk. Returns a default config if the file doesn't exist.
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return defaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Security Migration: Check if an old plaintext api_key exists and migrate it to the secure OS keyring
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		if oldKey, ok := raw["api_key"].(string); ok && oldKey != "" {
			_ = keyring.Set(keyringService, keyringUser, oldKey) // Secure it in the OS vault
			delete(raw, "api_key") // Remove from plaintext
			cleanData, _ := json.MarshalIndent(raw, "", "  ")
			_ = os.WriteFile(configPath, cleanData, 0600) // Overwrite file safely
			data = cleanData
		}
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig(), nil
	}

	// Populate APIKey from OS keyring if available
	if k, err := keyring.Get(keyringService, keyringUser); err == nil {
		cfg.APIKey = k
	}

	return &cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// SetAPIKey stores the API key in the OS keyring.
func (c *Config) SetAPIKey(key string) {
	c.APIKey = strings.TrimSpace(key)
	if c.APIKey != "" {
		_ = keyring.Set(keyringService, keyringUser, c.APIKey)
	}
}

// GetAPIKey returns the API key, checking env vars then the OS keyring.
func (c *Config) GetAPIKey() string {
	if key := os.Getenv("GOOGLE_API_KEY"); key != "" {
		return key
	}
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		return key
	}
	// Fallback to loaded keyring value
	return c.APIKey
}

// ClearAPIKey removes the stored API key from the OS keyring.
func (c *Config) ClearAPIKey() {
	c.APIKey = ""
	_ = keyring.Delete(keyringService, keyringUser)
}

// MaskedAPIKey returns a masked version of the API key for display.
func (c *Config) MaskedAPIKey() string {
	key := c.GetAPIKey()
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return strings.Repeat("•", len(key))
	}
	return key[:4] + strings.Repeat("•", len(key)-8) + key[len(key)-4:]
}

// GetModel returns the configured model or the default.
func (c *Config) GetModel() string {
	if c.Model != "" {
		return c.Model
	}
	return "gemini-2.5-flash"
}

// GetReportFormat returns the configured report format or the default ("md").
func (c *Config) GetReportFormat() string {
	if c.ReportFormat == "txt" {
		return "txt"
	}
	return "md"
}

// Reset clears all config values back to defaults.
func (c *Config) Reset() {
	c.APIKey = ""
	c.Model = ""
	c.DefaultPath = ""
	c.NoColor = false
}

func defaultConfig() *Config {
	return &Config{
		Model: "gemini-2.5-flash",
	}
}
