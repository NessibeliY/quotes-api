package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test_config_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	cfgData := Config{
		Port:    ":8080",
		LogFile: "test.log",
	}
	if err := json.NewEncoder(tmpFile).Encode(cfgData); err != nil {
		t.Fatalf("encode config: %v", err)
	}

	tmpFile.Close()

	originalFileName := fileName
	fileName = tmpFile.Name()
	defer func() {
		fileName = originalFileName
	}()

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.Port != ":8080" {
		t.Fatalf("expected port :8080, got: %s", cfg.Port)
	}
	if cfg.LogFile != "test.log" {
		t.Fatalf("expected log_file test.log, got: %s", cfg.LogFile)
	}
}
