package storage

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git-mini-commit/internal/types"
)

const (
	MiniCommitsDir = ".git/mini-commits"
	IndexFile      = "index.json"
)

// Storage manages mini-commit storage
type Storage struct {
	basePath string
	mutex    sync.RWMutex
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	// Check if current directory is a Git repository
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil, fmt.Errorf("git repository not found")
	}

	// Create mini-commits directory
	miniCommitsPath := filepath.Join(".git", "mini-commits")
	if err := os.MkdirAll(miniCommitsPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create mini-commits directory: %v", err)
	}

	return &Storage{
		basePath: miniCommitsPath,
	}, nil
}

// SaveMiniCommit saves a mini-commit
func (s *Storage) SaveMiniCommit(mc *types.MiniCommit) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()


	// Load existing index
	index, err := s.loadIndex()
	if err != nil {
		return fmt.Errorf("failed to load index: %v", err)
	}

	// Add new mini-commit
	index = append(index, *mc)

	// Save index
	if err := s.saveIndex(index); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	// Save patch file
	patchPath := filepath.Join(s.basePath, mc.ID+".patch")
	if err := os.WriteFile(patchPath, []byte(mc.Patch), 0644); err != nil {
		return fmt.Errorf("failed to save patch file: %v", err)
	}

	return nil
}

// LoadMiniCommits loads all mini-commits
func (s *Storage) LoadMiniCommits() (types.MiniCommitList, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.loadIndex()
}

// GetMiniCommit gets a mini-commit by ID
func (s *Storage) GetMiniCommit(id string) (*types.MiniCommit, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	index, err := s.loadIndex()
	if err != nil {
		return nil, err
	}

	for _, mc := range index {
		if mc.ID == id {
			return &mc, nil
		}
	}

	return nil, fmt.Errorf("mini-commit '%s' not found", id)
}

// DeleteMiniCommit deletes a mini-commit by ID
func (s *Storage) DeleteMiniCommit(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	index, err := s.loadIndex()
	if err != nil {
		return err
	}

	// Remove from index
	var newIndex types.MiniCommitList
	found := false
	for _, mc := range index {
		if mc.ID != id {
			newIndex = append(newIndex, mc)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("mini-commit '%s' not found", id)
	}

	// Save index
	if err := s.saveIndex(newIndex); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	// Delete patch file
	patchPath := filepath.Join(s.basePath, id+".patch")
	if err := os.Remove(patchPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete patch file: %v", err)
	}

	return nil
}

// ClearAllMiniCommits deletes all mini-commits
func (s *Storage) ClearAllMiniCommits() error {
	index, err := s.loadIndex()
	if err != nil {
		return err
	}

	// Delete all patch files
	for _, mc := range index {
		patchPath := filepath.Join(s.basePath, mc.ID+".patch")
		if err := os.Remove(patchPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete patch file '%s': %v", patchPath, err)
		}
	}

	// Save empty index
	if err := s.saveIndex(types.MiniCommitList{}); err != nil {
		return fmt.Errorf("failed to save index: %v", err)
	}

	return nil
}

// loadIndex loads the index file
func (s *Storage) loadIndex() (types.MiniCommitList, error) {
	indexPath := filepath.Join(s.basePath, IndexFile)

	// Return empty list if file doesn't exist
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return types.MiniCommitList{}, nil
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read index file: %v", err)
	}

	var index types.MiniCommitList
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index: %v", err)
	}

	return index, nil
}

// saveIndex saves the index file
func (s *Storage) saveIndex(index types.MiniCommitList) error {
	indexPath := filepath.Join(s.basePath, IndexFile)

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize index: %v", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save index file: %v", err)
	}

	return nil
}

// GenerateID generates ID from patch content and timestamp
func (s *Storage) GenerateID(patch string, timestamp time.Time) string {
	h := sha1.New()
	_, _ = io.WriteString(h, patch)
	_, _ = io.WriteString(h, timestamp.Format(time.RFC3339Nano))
	return fmt.Sprintf("%x", h.Sum(nil))
}
