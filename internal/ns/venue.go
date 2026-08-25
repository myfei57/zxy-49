package ns

import (
	"errors"
	"strings"

	"venueops/internal/store"
)

func (s *Service) RenameVenue(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("venue name is required")
	}
	venue, err := s.GetVenue(id)
	if err != nil {
		return err
	}
	venue.Name = name
	return s.store.Save(store.VenueStoreKey(id), venue)
}

func (s *Service) ZonesOfVenue(id string) ([]string, error) {
	venue, err := s.GetVenue(id)
	if err != nil {
		return nil, err
	}
	return venue.ZoneIDs, nil
}

func (s *Service) HasVenue(id string) bool {
	return s.store.Exists(store.VenueStoreKey(id))
}

func (s *Service) CountVenues() (int, error) {
	keys, err := s.store.Keys(store.VenueKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) VenueNames() ([]string, error) {
	venues, err := s.ListVenues()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(venues))
	for _, venue := range venues {
		names = append(names, venue.Name)
	}
	return names, nil
}

func (s *Service) EnsureVenue(id string, name string) error {
	if s.HasVenue(id) {
		return nil
	}
	return s.CreateVenue(id, name)
}

func (s *Service) ValidateVenue(id string) error {
	if !s.HasVenue(id) {
		return errors.New("venue not found: " + id)
	}
	return nil
}
