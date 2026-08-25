package store

import (
	"os"
	"path/filepath"
)

func (s *Store) SaveIn(sub string, key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	base := filepath.Join(s.dir, sub)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return err
	}
	return s.saveLocked(base, key, value)
}

func (s *Store) LoadIn(sub string, key string, dst any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked(filepath.Join(s.dir, sub), key, dst)
}
