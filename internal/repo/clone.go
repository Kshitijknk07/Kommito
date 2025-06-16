package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CloneRepo(source, destination string) error {
	if strings.HasPrefix(source, "http") || strings.HasPrefix(source, "git@") {
		return cloneGitRepo(source, destination)
	}

	sourceKommito := filepath.Join(source, ".kommito")
	if _, err := os.Stat(sourceKommito); os.IsNotExist(err) {
		return fmt.Errorf("source is not a valid Kommito repository: %w", err)
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	if err := copyDir(sourceKommito, filepath.Join(destination, ".kommito")); err != nil {
		return fmt.Errorf("failed to copy .kommito directory: %w", err)
	}

	indexPath := filepath.Join(sourceKommito, "index")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index: %w", err)
	}

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

		sourceFile := filepath.Join(source, filePath)
		destFile := filepath.Join(destination, filePath)

		if err := os.MkdirAll(filepath.Dir(destFile), 0755); err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %w", filePath, err)
		}

		if err := copyFile(sourceFile, destFile); err != nil {
			return fmt.Errorf("failed to copy file %s: %w", filePath, err)
		}
	}

	return nil
}

func cloneGitRepo(gitURL, destination string) error {

	tempDir, err := os.MkdirTemp("", "kommito-git-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	cmd := exec.Command("git", "clone", gitURL, tempDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone Git repository: %w", err)
	}

	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	if err := os.Chdir(destination); err != nil {
		return fmt.Errorf("failed to change to destination directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := InitRepo(); err != nil {
		return fmt.Errorf("failed to initialize Kommito repository: %w", err)
	}

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("failed to read Git repository: %w", err)
	}

	var addedFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if name == ".git" || isSystemFile(name) {
			continue
		}

		srcPath := filepath.Join(tempDir, name)
		dstPath := name

		if entry.IsDir() {

			if err := os.MkdirAll(dstPath, 0755); err != nil {
				fmt.Printf("(╥﹏╥) Could not create directory %s: %v\n", name, err)
				continue
			}

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

	if err := AddFile("."); err != nil {
		return fmt.Errorf("failed to add files to Kommito: %w", err)
	}

	if err := CommitStaged("Initial commit from Git repository"); err != nil {
		return fmt.Errorf("failed to create initial commit: %w", err)
	}

	fmt.Printf("(＾▽＾) Successfully added %d files!\n", len(addedFiles))
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
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

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, data, 0644)
}
