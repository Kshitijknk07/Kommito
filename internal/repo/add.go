package repo

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func isSystemFile(name string) bool {
	systemFiles := []string{
		".git",
		".kommito",
		"Application Data",
		"Cookies",
		"Local Settings",
		"My Documents",
		"NTUSER.DAT",
		"NetHood",
		"PrintHood",
		"Recent",
		"SendTo",
		"Start Menu",
		"Templates",
		"ntuser.dat.LOG1",
		"ntuser.dat.LOG2",
		"My Music",
		"My Pictures",
		"My Videos",
		"Kommito",
	}

	for _, sysFile := range systemFiles {
		if strings.EqualFold(name, sysFile) {
			return true
		}
	}

	return false
}

func AddFile(path string) error {
	if path == "." {
		entries, err := os.ReadDir(".")
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}

		var addedFiles []string
		for _, entry := range entries {
			name := entry.Name()

			if isSystemFile(name) {
				continue
			}

			if entry.IsDir() {
				// Skip directories
				continue
			}

			if err := addSingleFile(name); err != nil {
				fmt.Printf("(╥﹏╥) Could not add %s: %v\n", name, err)
				continue
			}
			addedFiles = append(addedFiles, name)
		}

		if len(addedFiles) == 0 {
			fmt.Println("(⊙_☉) No files to add!")
			return nil
		}

		fmt.Printf("(＾▽＾) Successfully added %d files!\n", len(addedFiles))
		return nil
	}

	return addSingleFile(path)
}

func addSingleFile(filePath string) error {
	if isSystemFile(filePath) {
		return fmt.Errorf("skipping system file: %s", filePath)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))

	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to rewind file: %w", err)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	blobPath := filepath.Join(".kommito", "objects", "blobs", hash)
	if err := os.WriteFile(blobPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}

	indexLine := fmt.Sprintf("%s %s\n", hash, filePath)
	indexPath := filepath.Join(".kommito", "index")
	fidx, err := os.OpenFile(indexPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open index: %w", err)
	}
	defer fidx.Close()
	if _, err := fidx.WriteString(indexLine); err != nil {
		return fmt.Errorf("failed to update index: %w", err)
	}

	return nil
}
