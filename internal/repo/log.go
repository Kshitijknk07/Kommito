package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func LogCommits() error {
	headPath := filepath.Join(".kommito", "HEAD")
	headData, err := ioutil.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %w", err)
	}
	commitHash := string(headData)
	commitPath := filepath.Join(".kommito", "objects", "commits", commitHash)
	commitData, err := ioutil.ReadFile(commitPath)
	if err != nil {
		return fmt.Errorf("failed to read commit: %w", err)
	}
	var commit Commit
	if err := json.Unmarshal(commitData, &commit); err != nil {
		return fmt.Errorf("failed to parse commit: %w", err)
	}
	fmt.Printf("🕐 Commit: %s\n📜 Message: %s\n👤 Author: %s\n🕰️ Date: %s\n", commitHash, commit.Message, commit.Author, commit.Timestamp)
	return nil
}
