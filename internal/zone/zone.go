package zone

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

func (s *Service) CreateZone(id string, hallID string, venueID string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("zone id is required")
	}
	if strings.TrimSpace(hallID) == "" {
		return errors.New("hall id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name is required")
	}
	if s.store.Exists(store.ZoneStoreKey(id)) {
		return errors.New("zone already exists: " + id)
	}
	zone := domain.Zone{
		ID:        id,
		HallID:    hallID,
		VenueID:   venueID,
		Name:      name,
		Partition: id,
	}
	return s.store.Save(store.ZoneStoreKey(id), zone)
}

func (s *Service) GetZone(id string) (domain.Zone, error) {
	var zone domain.Zone
	err := s.store.Load(store.ZoneStoreKey(id), &zone)
	return zone, err
}

func (s *Service) ListZones() ([]domain.Zone, error) {
	keys, err := s.store.Keys(store.ZoneKey)
	if err != nil {
		return nil, err
	}
	zones := make([]domain.Zone, 0, len(keys))
	for _, key := range keys {
		var zone domain.Zone
		if err := s.store.Load(key, &zone); err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}
	return zones, nil
}

func (s *Service) ZonesOfHall(hallID string) ([]domain.Zone, error) {
	zones, err := s.ListZones()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Zone, 0)
	for _, zone := range zones {
		if zone.HallID == hallID {
			filtered = append(filtered, zone)
		}
	}
	return filtered, nil
}

func (s *Service) HasZone(id string) bool {
	return s.store.Exists(store.ZoneStoreKey(id))
}

func (s *Service) CountZones() (int, error) {
	keys, err := s.store.Keys(store.ZoneKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) RenameZone(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("zone name is required")
	}
	zone, err := s.GetZone(id)
	if err != nil {
		return err
	}
	zone.Name = name
	return s.store.Save(store.ZoneStoreKey(id), zone)
}

func (s *Service) HallOfZone(id string) (string, error) {
	zone, err := s.GetZone(id)
	if err != nil {
		return "", err
	}
	return zone.HallID, nil
}

func (s *Service) VenueOfZone(id string) (string, error) {
	zone, err := s.GetZone(id)
	if err != nil {
		return "", err
	}
	return zone.VenueID, nil
}
