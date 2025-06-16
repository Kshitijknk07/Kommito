package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CloneRepo copies a repository from source to destination
func CloneRepo(source, destination string) error {
	// Check if source is a Git repository URL
	if strings.HasPrefix(source, "http") || strings.HasPrefix(source, "git@") {
		return cloneGitRepo(source, destination)
	}

	// Handle local Kommito repository
	sourceKommito := filepath.Join(source, ".kommito")
	if _, err := os.Stat(sourceKommito); os.IsNotExist(err) {
		return fmt.Errorf("source is not a valid Kommito repository: %w", err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy .kommito directory
	if err := copyDir(sourceKommito, filepath.Join(destination, ".kommito")); err != nil {
		return fmt.Errorf("failed to copy .kommito directory: %w", err)
	}

	// Copy all tracked files
	indexPath := filepath.Join(sourceKommito, "index")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

	// Parse index and copy files
	lines := strings.Split(strings.TrimSpace(string(indexData)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		filePath := parts[1]

		// Copy file from source to destination
		sourceFile := filepath.Join(source, filePath)
		destFile := filepath.Join(destination, filePath)

		// Create parent directories if they don't exist
		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", filePath, err)
		}

		// Copy file
		if err := copyFile(sourceFile, destFile); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", filePath, err)
		}
	}

	return nil
}

// cloneGitRepo clones a Git repository and converts it to a Kommito repository
func cloneGitRepo(gitURL, destination string) error {
	// Create a temporary directory for Git clone
	tempDir, err := os.MkdirTemp("", "kommito-git-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the Git repository
	cmd := exec.Command("git", "clone", gitURL, tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone Git repository: %w", err)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Initialize new Kommito repository
	if err := InitRepo(); err != nil {
		return fmt.Errorf("failed to initialize Kommito repository: %w", err)
	}

	// Copy all files from Git repository to Kommito repository
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read Git repository: %w", err)
	}

	var addedFiles []string
	for _, entry := range entries {
		name := entry.Name()
		// Skip .git directory and system files
		if name == ".git" || isSystemFile(name) {
			continue
		}

		srcPath := filepath.Join(tempDir, name)
		dstPath := filepath.Join(destination, name)

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				fmt.Printf("(╥﹏╥) Could not copy directory %s: %v\n", name, err)
				continue
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				fmt.Printf("(╥﹏╥) Could not copy file %s: %v\n", name, err)
				continue
			}
		}
		addedFiles = append(addedFiles, name)
	}

	// Add all files to Kommito
	if err := AddFile("."); err != nil {
		return fmt.Errorf("failed to add files to Kommito: %w", err)
	}

	// Create initial commit
	if err := CommitStaged("Initial commit from Git repository"); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	fmt.Printf("(＾▽＾) Successfully added %d files!\n", len(addedFiles))
	return nil
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	// Create destination directory
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Copy each entry
	for _, entry := range entries {
		name := entry.Name()
		// Skip system files
		if isSystemFile(name) {
			continue
		}

		srcPath := filepath.Join(src, name)
		dstPath := filepath.Join(dst, name)

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write to destination
	return os.WriteFile(dst, data, 0644)
}
