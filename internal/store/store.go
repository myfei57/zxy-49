package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Store struct {
	dir string
	mu  sync.Mutex
}

var ErrNotFound = errors.New("store key not found")

func New(dir string) (*Store, error) {
	if strings.TrimSpace(dir) == "" {
		return nil, errors.New("store directory is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

func (s *Store) path(key string) string {
	return filepath.Join(s.dir, key+".json")
}

func (s *Store) Save(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.saveLocked(s.dir, key, value)
}

func (s *Store) saveLocked(base string, key string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	target := filepath.Join(base, key+".json")
	tmp := target + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, target)
}

func (s *Store) Load(key string, dst any) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.loadLocked(s.dir, key, dst)
}

func (s *Store) loadLocked(base string, key string, dst any) error {
	data, err := os.ReadFile(filepath.Join(base, key+".json"))
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotFound
		}
		return err
	}
	return json.Unmarshal(data, dst)
}

func (s *Store) Exists(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := os.Stat(s.path(key))
	return err == nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) Keys(prefix string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		key := strings.TrimSuffix(name, ".json")
		if prefix != "" && !strings.HasPrefix(key, prefix) {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}
