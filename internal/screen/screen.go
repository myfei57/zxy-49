package screen

import (
	"errors"
	"strings"

	"venueops/internal/domain"
	"venueops/internal/store"
)

type Service struct {
	store *store.Store
}

func NewService(st *store.Store) *Service {
	return &Service{store: st}
}

func (s *Service) Create(id string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("screen id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("screen name is required")
	}
	if s.store.Exists(store.ScreenStoreKey(id)) {
		return errors.New("screen already exists: " + id)
	}
	return s.store.Save(store.ScreenStoreKey(id), domain.Screen{ID: id, Name: name})
}

func (s *Service) Get(id string) (domain.Screen, error) {
	var screen domain.Screen
	err := s.store.Load(store.ScreenStoreKey(id), &screen)
	return screen, err
}

func (s *Service) List() ([]domain.Screen, error) {
	keys, err := s.store.Keys(store.ScreenKey)
	if err != nil {
		return nil, err
	}
	screens := make([]domain.Screen, 0, len(keys))
	for _, key := range keys {
		var screen domain.Screen
		if err := s.store.Load(key, &screen); err != nil {
			return nil, err
		}
		screens = append(screens, screen)
	}
	return screens, nil
}

func (s *Service) Rename(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("screen name is required")
	}
	screen, err := s.Get(id)
	if err != nil {
		return err
	}
	screen.Name = name
	return s.store.Save(store.ScreenStoreKey(id), screen)
}

func (s *Service) Count() (int, error) {
	keys, err := s.store.Keys(store.ScreenKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) HasScreen(id string) bool {
	return s.store.Exists(store.ScreenStoreKey(id))
}

func (s *Service) ScreenIDs() ([]string, error) {
	screens, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(screens))
	for _, screen := range screens {
		ids = append(ids, screen.ID)
	}
	return ids, nil
}

func (s *Service) Names() ([]string, error) {
	screens, err := s.List()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(screens))
	for _, screen := range screens {
		names = append(names, screen.Name)
	}
	return names, nil
}
