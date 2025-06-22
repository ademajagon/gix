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
}

const configFileName = "config.json"
const appName = "toka"

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
		return Config{}, errors.New("config file not found. Please run `toka config set-key <key>` to set your OpenAI key")
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.OpenAIKey == "" {
		return Config{}, errors.New("OpenAI key is missing in config")
	}

	return cfg, nil
}
