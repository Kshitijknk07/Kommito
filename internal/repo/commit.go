package repo

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Commit struct {
	Author    string   `json:"author"`
	Timestamp string   `json:"timestamp"`
	Message   string   `json:"message"`
	Blobs     []string `json:"blobs"`
}

func CommitStaged(message string) error {

	indexPath := filepath.Join(".kommito", "index")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(indexData)), "\n")
	var blobs []string
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) > 0 {
			blobs = append(blobs, parts[0])
		}
	}

	author := "Kommito User"
	configPath := filepath.Join(".kommito", "config.json")
	if configData, err := os.ReadFile(configPath); err == nil {
		var cfg struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(configData, &cfg); err == nil && cfg.Name != "" {
			author = cfg.Name
		}
	}

	commit := Commit{
		Author:    author,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   message,
		Blobs:     blobs,
	}

	commitBytes, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %w", err)
	}

	h := sha1.New()
	h.Write(commitBytes)
	commitHash := fmt.Sprintf("%x", h.Sum(nil))

	commitPath := filepath.Join(".kommito", "objects", "commits", commitHash)
	if err := os.WriteFile(commitPath, commitBytes, 0644); err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}

	headPath := filepath.Join(".kommito", "HEAD")
	if err := os.WriteFile(headPath, []byte(commitHash), 0644); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	return nil
}
