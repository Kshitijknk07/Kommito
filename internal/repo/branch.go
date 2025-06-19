package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

type Branch struct {
	Name   string `json:"name"`
	Commit string `json:"commit"`
}

type BranchManager struct {
	repoPath string
}

func NewBranchManager(repoPath string) *BranchManager {
	return &BranchManager{
		repoPath: repoPath,
	}
}

func (bm *BranchManager) CreateBranch(name string) error {

	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	branches, err := bm.ListBranches()
	if err != nil {
		return err
	}
	for _, branch := range branches {
		if branch.Name == name {
			return fmt.Errorf("branch '%s' already exists", name)
		}
	}

	headPath := filepath.Join(bm.repoPath, ".kommito", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %v", err)
	}

	branchPath := filepath.Join(bm.repoPath, ".kommito", "refs", "heads", name)
	if err := os.MkdirAll(filepath.Dir(branchPath), 0755); err != nil {
		return fmt.Errorf("failed to create branch directory: %v", err)
	}

	if err := os.WriteFile(branchPath, headContent, 0644); err != nil {
		return fmt.Errorf("failed to create branch: %v", err)
	}

	return nil
}

func (bm *BranchManager) SwitchBranch(name string) error {

	branchPath := filepath.Join(bm.repoPath, ".kommito", "refs", "heads", name)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", name)
	}

	branchContent, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("failed to read branch: %v", err)
	}

	headPath := filepath.Join(bm.repoPath, ".kommito", "HEAD")
	if err := os.WriteFile(headPath, branchContent, 0644); err != nil {
		return fmt.Errorf("failed to switch branch: %v", err)
	}

	return nil
}

func (bm *BranchManager) ListBranches() ([]Branch, error) {
	headsPath := filepath.Join(bm.repoPath, ".kommito", "refs", "heads")
	entries, err := os.ReadDir(headsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Branch{}, nil
		}
		return nil, fmt.Errorf("failed to read branches: %v", err)
	}

	var branches []Branch
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		branchPath := filepath.Join(headsPath, entry.Name())
		commit, err := os.ReadFile(branchPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read branch '%s': %v", entry.Name(), err)
		}

		branches = append(branches, Branch{
			Name:   entry.Name(),
			Commit: string(commit),
		})
	}

	return branches, nil
}

func (bm *BranchManager) DeleteBranch(name string) error {

	branchPath := filepath.Join(bm.repoPath, ".kommito", "refs", "heads", name)
	if _, err := os.Stat(branchPath); os.IsNotExist(err) {
		return fmt.Errorf("branch '%s' does not exist", name)
	}

	headPath := filepath.Join(bm.repoPath, ".kommito", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return fmt.Errorf("failed to read HEAD: %v", err)
	}

	branchContent, err := os.ReadFile(branchPath)
	if err != nil {
		return fmt.Errorf("failed to read branch: %v", err)
	}

	if string(headContent) == string(branchContent) {
		return fmt.Errorf("cannot delete current branch")
	}

	if err := os.Remove(branchPath); err != nil {
		return fmt.Errorf("failed to delete branch: %v", err)
	}

	return nil
}

func (bm *BranchManager) GetCurrentBranch() (string, error) {
	headPath := filepath.Join(bm.repoPath, ".kommito", "HEAD")
	headContent, err := os.ReadFile(headPath)
	if err != nil {
		return "", fmt.Errorf("failed to read HEAD: %v", err)
	}

	branches, err := bm.ListBranches()
	if err != nil {
		return "", err
	}

	for _, branch := range branches {
		if branch.Commit == string(headContent) {
			return branch.Name, nil
		}
	}

	return "", fmt.Errorf("not on any branch")
}

func (bm *BranchManager) GetBranchCommit(name string) (string, error) {
	branchPath := filepath.Join(bm.repoPath, ".kommito", "refs", "heads", name)
	commit, err := os.ReadFile(branchPath)
	if err != nil {
		return "", err
	}
	return string(commit), nil
}
