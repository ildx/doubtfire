package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DestinationDir string `json:"destination_dir"`
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(configPath)
	if os.IsNotExist(err) {
		return &Config{}, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(cfg *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cfg)
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".doubtfire_config.json"), nil
}
