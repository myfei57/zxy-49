package hall

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

func (s *Service) CreateHall(id string, venueID string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("hall id is required")
	}
	if strings.TrimSpace(venueID) == "" {
		return errors.New("venue id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("hall name is required")
	}
	if s.store.Exists(store.HallStoreKey(id)) {
		return errors.New("hall already exists: " + id)
	}
	hall := domain.Hall{
		ID:      id,
		VenueID: venueID,
		Name:    name,
		ZoneIDs: []string{},
		Layout:  map[string]domain.Point{},
	}
	return s.store.Save(store.HallStoreKey(id), hall)
}

func (s *Service) GetHall(id string) (domain.Hall, error) {
	var hall domain.Hall
	err := s.store.Load(store.HallStoreKey(id), &hall)
	return hall, err
}

func (s *Service) ListHalls() ([]domain.Hall, error) {
	keys, err := s.store.Keys(store.HallKey)
	if err != nil {
		return nil, err
	}
	halls := make([]domain.Hall, 0, len(keys))
	for _, key := range keys {
		var hall domain.Hall
		if err := s.store.Load(key, &hall); err != nil {
			return nil, err
		}
		halls = append(halls, hall)
	}
	return halls, nil
}

func (s *Service) HallsOfVenue(venueID string) ([]domain.Hall, error) {
	halls, err := s.ListHalls()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Hall, 0)
	for _, hall := range halls {
		if hall.VenueID == venueID {
			filtered = append(filtered, hall)
		}
	}
	return filtered, nil
}

func (s *Service) HasHall(id string) bool {
	return s.store.Exists(store.HallStoreKey(id))
}

func (s *Service) CountHalls() (int, error) {
	keys, err := s.store.Keys(store.HallKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) RenameHall(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("hall name is required")
	}
	hall, err := s.GetHall(id)
	if err != nil {
		return err
	}
	hall.Name = name
	return s.store.Save(store.HallStoreKey(id), hall)
}

func (s *Service) ZoneCount(id string) (int, error) {
	hall, err := s.GetHall(id)
	if err != nil {
		return 0, err
	}
	return len(hall.ZoneIDs), nil
}
