package booth

import (
	"errors"
	"strings"

	"venueops/internal/domain"
	"venueops/internal/quota"
	"venueops/internal/store"
)

type Service struct {
	store  *store.Store
	quotaS *quota.Service
}

func NewService(st *store.Store, quotaSvc *quota.Service) *Service {
	return &Service{store: st, quotaS: quotaSvc}
}

func (s *Service) Create(id string, hallID string, zoneID string, name string, loadWatts int) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("booth id is required")
	}
	if strings.TrimSpace(hallID) == "" {
		return errors.New("hall id is required")
	}
	if strings.TrimSpace(zoneID) == "" {
		return errors.New("zone id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("booth name is required")
	}
	if loadWatts < 0 {
		return errors.New("load watts cannot be negative")
	}
	if s.store.Exists(store.BoothStoreKey(id)) {
		return errors.New("booth already exists: " + id)
	}
	booth := domain.Booth{
		ID:        id,
		HallID:    hallID,
		ZoneID:    zoneID,
		Name:      name,
		Powered:   false,
		LoadWatts: loadWatts,
	}
	return s.store.Save(store.BoothStoreKey(id), booth)
}

func (s *Service) Get(id string) (domain.Booth, error) {
	var booth domain.Booth
	err := s.store.Load(store.BoothStoreKey(id), &booth)
	return booth, err
}

func (s *Service) List() ([]domain.Booth, error) {
	keys, err := s.store.Keys(store.BoothKey)
	if err != nil {
		return nil, err
	}
	booths := make([]domain.Booth, 0, len(keys))
	for _, key := range keys {
		var booth domain.Booth
		if err := s.store.Load(key, &booth); err != nil {
			return nil, err
		}
		booths = append(booths, booth)
	}
	return booths, nil
}

func (s *Service) BoothsOfHall(hallID string) ([]domain.Booth, error) {
	booths, err := s.List()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Booth, 0)
	for _, booth := range booths {
		if booth.HallID == hallID {
			filtered = append(filtered, booth)
		}
	}
	return filtered, nil
}

func (s *Service) SetLoad(id string, loadWatts int) error {
	if loadWatts < 0 {
		return errors.New("load watts cannot be negative")
	}
	booth, err := s.Get(id)
	if err != nil {
		return err
	}
	booth.LoadWatts = loadWatts
	return s.store.Save(store.BoothStoreKey(id), booth)
}

func (s *Service) Rename(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("booth name is required")
	}
	booth, err := s.Get(id)
	if err != nil {
		return err
	}
	booth.Name = name
	return s.store.Save(store.BoothStoreKey(id), booth)
}

func (s *Service) HasBooth(id string) bool {
	return s.store.Exists(store.BoothStoreKey(id))
}

func (s *Service) Count() (int, error) {
	keys, err := s.store.Keys(store.BoothKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}
