package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)


func CheckoutTarget(target string) error {
	bm := NewBranchManager(".")
	commitHash := target

	
	branches, err := bm.ListBranches()
	if err == nil {
		for _, branch := range branches {
			if branch.Name == target {
				commitHash = branch.Commit
				break
			}
		}
	}

	commit, err := LoadCommit(commitHash)
	if err != nil {
		return fmt.Errorf("could not find commit or branch '%s': %w", target, err)
	}

	
	blobToPath, err := getBlobToPathMap(commit)
	if err != nil {
		return err
	}
	filesInCommit := make(map[string]string) 
	for hash, path := range blobToPath {
		filesInCommit[path] = hash
	}

	
	entries, err := ioutil.ReadDir(".")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || isSystemFile(name) {
			continue
		}
		if _, ok := filesInCommit[name]; !ok {
			_ = os.Remove(name)
		}
	}

	
	for path, hash := range filesInCommit {
		blobPath := filepath.Join(".kommito", "objects", "blobs", hash)
		content, err := os.ReadFile(blobPath)
		if err != nil {
			return fmt.Errorf("failed to read blob for %s: %w", path, err)
		}
		if err := ioutil.WriteFile(path, content, 0644); err != nil {
			return fmt.Errorf("failed to restore file %s: %w", path, err)
		}
	}

	
	for _, branch := range branches {
		if branch.Name == target {
			headPath := filepath.Join(".kommito", "HEAD")
			if err := os.WriteFile(headPath, []byte(branch.Commit), 0644); err != nil {
				return fmt.Errorf("failed to update HEAD: %w", err)
			}
			break
		}
	}

	fmt.Printf("Checked out %s\n", target)
	return nil
} 