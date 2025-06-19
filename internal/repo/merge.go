package repo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileEntry struct {
	Hash string
	Path string
}


func LoadCommit(hash string) (*Commit, error) {
	commitPath := filepath.Join(".kommito", "objects", "commits", hash)
	data, err := os.ReadFile(commitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read commit object: %w", err)
	}
	var commit Commit
	if err := json.Unmarshal(data, &commit); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commit: %w", err)
	}
	return &commit, nil
}


func getBlobToPathMap(commit *Commit) (map[string]string, error) {
	indexPath := filepath.Join(".kommito", "index")
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index: %w", err)
	}
	blobToPath := make(map[string]string)
	lines := string(indexData)
	for _, line := range splitLines(lines) {
		if line == "" {
			continue
		}
		var hash, path string
		fmt.Sscanf(line, "%s %s", &hash, &path)
		blobToPath[hash] = path
	}
	return blobToPath, nil
}


func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := range s {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}


func MergeBranches(targetBranch string) error {
	bm := NewBranchManager(".")
	currentBranch, err := bm.GetCurrentBranch()
	if err != nil {
		return err
	}
	if currentBranch == targetBranch {
		return fmt.Errorf("cannot merge branch '%s' into itself", targetBranch)
	}
	currentCommitHash, err := bm.GetBranchCommit(currentBranch)
	if err != nil {
		return err
	}
	targetCommitHash, err := bm.GetBranchCommit(targetBranch)
	if err != nil {
		return err
	}
	currentCommit, err := LoadCommit(currentCommitHash)
	if err != nil {
		return err
	}
	targetCommit, err := LoadCommit(targetCommitHash)
	if err != nil {
		return err
	}
	
	blobToPath, err := getBlobToPathMap(currentCommit)
	if err != nil {
		return err
	}
	conflicts := []string{}
	for _, targetBlob := range targetCommit.Blobs {
		path, ok := blobToPath[targetBlob]
		if !ok {
			
			blobPath := filepath.Join(".kommito", "objects", "blobs", targetBlob)
			content, err := os.ReadFile(blobPath)
			if err != nil {
				return fmt.Errorf("failed to read blob: %w", err)
			}
			if err := ioutil.WriteFile(path, content, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
			continue
		}
		
		currentBlobPath := filepath.Join(".kommito", "objects", "blobs", targetBlob)
		currentContent, _ := os.ReadFile(currentBlobPath)
		targetBlobPath := filepath.Join(".kommito", "objects", "blobs", targetBlob)
		targetContent, _ := os.ReadFile(targetBlobPath)
		if string(currentContent) != string(targetContent) {
			
			conflictContent := []byte(
				"<<<<<<< " + currentBranch + "\n" +
				string(currentContent) +
				"\n=======\n" +
				string(targetContent) +
				"\n>>>>>>> " + targetBranch + "\n")
			if err := ioutil.WriteFile(path, conflictContent, 0644); err != nil {
				return fmt.Errorf("failed to write conflict file: %w", err)
			}
			conflicts = append(conflicts, path)
		} else {
			
			if err := ioutil.WriteFile(path, currentContent, 0644); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	}
	if len(conflicts) > 0 {
		fmt.Println("Merge completed with conflicts in:")
		for _, c := range conflicts {
			fmt.Println("  ", c)
		}
		fmt.Println("Please resolve conflicts and commit.")
	} else {
		fmt.Println("Merge completed successfully. No conflicts detected.")
	}
	return nil
} 