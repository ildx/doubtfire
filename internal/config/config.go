package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	LastCleanupDate      time.Time `json:"last_cleanup_date"`
	DestinationDirectory string    `json:"destination_directory"`
}

func LoadConfig() (*Config, error) {
	configPath := getConfigPath()
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath := getConfigPath()
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".doubtfire.json")
}
