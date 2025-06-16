package repo

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func InitRepo() error {
	if err := os.MkdirAll(".kommito", 0755); err != nil {
		return fmt.Errorf("failed to create .kommito directory: %w", err)
	}
	dirs := []string{
		".kommito/objects/commits",
		".kommito/objects/blobs",
		".kommito/refs/heads",
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	if err := os.WriteFile(".kommito/HEAD", []byte("ref: refs/heads/main"), 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %w", err)
	}
	if err := os.WriteFile(".kommito/index", []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create index file: %w", err)
	}
	config := Config{
		Name:    "kommito",
		Version: "0.1.0",
	}
	configData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	if err := os.WriteFile(".kommito/config.json", configData, 0644); err != nil {
		return fmt.Errorf("failed to create config.json: %w", err)
	}
	return nil
}
