package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	BaseUrl string `json:"baseUrl"`
	Model   string `json:"model"`
}

func LoadConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(home, ".config", "q", "config.json")
	if _, err := os.Stat(configPath); err != nil {
		return nil, err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}

	return config, nil
}
