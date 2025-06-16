package repo

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// AddFile stages a file: reads, hashes, stores as blob, and updates index
func AddFile(filePath string) error {
	// Read file contents
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	// Get hash as hex string
	hash := fmt.Sprintf("%x", h.Sum(nil))

	// Reset file pointer and read contents again
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("failed to rewind file: %w", err)
	}
	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Store blob
	blobPath := filepath.Join(".kommito", "objects", "blobs", hash)
	if err := os.WriteFile(blobPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write blob: %w", err)
	}

	// Update index (append line: <hash> <filepath>)
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
