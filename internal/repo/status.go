package repo

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Status() error {
	// Read index
	indexPath := filepath.Join(".kommito", "index")
	indexData, _ := os.ReadFile(indexPath)
	indexLines := strings.Split(strings.TrimSpace(string(indexData)), "\n")
	staged := make(map[string]string) // file -> hash
	for _, line := range indexLines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			staged[parts[1]] = parts[0]
		}
	}

	fmt.Println("🗂️ Staged files:")
	if len(staged) == 0 {
		fmt.Println("  (none)")
	} else {
		for file := range staged {
			fmt.Printf("  ➕ %s\n", file)
		}
	}

	// Find modified but unstaged files
	fmt.Println("\n✏️ Modified but unstaged files:")
	modified := false
	for file, hash := range staged {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		h := sha1.New()
		h.Write(content)
		newHash := fmt.Sprintf("%x", h.Sum(nil))
		if newHash != hash {
			fmt.Printf("  📝 %s\n", file)
			modified = true
		}
	}
	if !modified {
		fmt.Println("  (none)")
	}

	// Find untracked files
	fmt.Println("\n❓ Untracked files:")
	entries, _ := os.ReadDir(".")
	untracked := false
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || name == ".kommito" || name == ".git" {
			continue
		}
		if _, ok := staged[name]; !ok {
			fmt.Printf("  ❔ %s\n", name)
			untracked = true
		}
	}
	if !untracked {
		fmt.Println("  (none)")
	}
	return nil
}
