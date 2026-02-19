package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	OpenAIKey string `json:"openai_key"`
	GeminiKey string `json:"gemini_key,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

const configFileName = "config.json"
const appName = "gix"

func getConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, appName), nil
}

func getConfigPath() (string, error) {
	dir, err := getConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func Save(cfg Config) error {
	dir, err := getConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path := filepath.Join(dir, configFileName)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0600)
	if err != nil {
		fmt.Println("Failed to write config file:", err)
	} else {
		fmt.Println("Config file written to:", path)
	}
	return err
}

func Load() (Config, error) {
	path, err := getConfigPath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return Config{}, errors.New("config file not found. Please run `gix config set-key` to set your API key")
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// ResolveProvider returns the active provider name, defaulting to "openai".
func (c Config) ResolveProvider() string {
	if c.Provider != "" {
		return c.Provider
	}
	return "openai"
}

// APIKey returns the API key for the active provider.
func (c Config) APIKey() string {
	switch c.ResolveProvider() {
	case "gemini":
		return c.GeminiKey
	default:
		return c.OpenAIKey
	}
}
