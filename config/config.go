package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appName        = "gix"
	configFileName = "config.json"
)

// Config contains all persistent settings for gix
type Config struct {
	OpenAIKey string `json:"openai_key,omitempty"`
	GeminiKey string `json:"gemini_key,omitempty"`
	Provider  string `json:"provider,omitempty"`

	DisableUpdateCheck bool `json:"disable_update_check,omitempty"`
}

// ResolveProvider returns the active provider name, defaulting to "openai".
func (c Config) ResolveProvider() string {
	if c.Provider != "" {
		return c.Provider
	}
	return "openai"
}

func (c Config) APIKey() string {
	switch c.ResolveProvider() {
	case "gemini":
		return c.GeminiKey
	default:
		return c.OpenAIKey
	}
}

// Load reads configuration from disk
func Load() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, fmt.Errorf("resolving config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("config not found - run `gix config set-key`")
		}
		return Config{}, fmt.Errorf("reading config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func Save(cfg Config) error {
	dir, err := configDir()
	if err != nil {
		return fmt.Errorf("resolving config dir: %w", err)
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling config: %w", err)
	}

	dest := filepath.Join(dir, configFileName)

	tmp := dest + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("saving config: %w", err)
	}

	return nil
}

// Path returns the resolved path to the config file
func Path() (string, error) {
	return configPath()
}

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, appName), nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}
