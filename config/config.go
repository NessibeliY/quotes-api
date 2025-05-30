package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	fileName = "config.json"
)

type Config struct {
	Port    string `json:"port"`
	LogFile string `json:"log_file"`
}

func LoadConfig() (*Config, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("open config file: %v", err)
	}
	defer f.Close()

	cfg := &Config{}
	err = json.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("decode file: %v", err)
	}

	return cfg, nil
}
