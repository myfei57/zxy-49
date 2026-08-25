package ns

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

func (s *Service) CreateVenue(id string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("venue id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("venue name is required")
	}
	if s.store.Exists(store.VenueStoreKey(id)) {
		return errors.New("venue already exists: " + id)
	}
	venue := domain.Venue{ID: id, Name: name, ZoneIDs: []string{}}
	return s.store.Save(store.VenueStoreKey(id), venue)
}

func (s *Service) GetVenue(id string) (domain.Venue, error) {
	var venue domain.Venue
	err := s.store.Load(store.VenueStoreKey(id), &venue)
	return venue, err
}

func (s *Service) ListVenues() ([]domain.Venue, error) {
	keys, err := s.store.Keys(store.VenueKey)
	if err != nil {
		return nil, err
	}
	venues := make([]domain.Venue, 0, len(keys))
	for _, key := range keys {
		var venue domain.Venue
		if err := s.store.Load(key, &venue); err != nil {
			return nil, err
		}
		venues = append(venues, venue)
	}
	return venues, nil
}

func (s *Service) AttachZone(venueID string, zoneID string) error {
	venue, err := s.GetVenue(venueID)
	if err != nil {
		return err
	}
	for _, existing := range venue.ZoneIDs {
		if existing == zoneID {
			return nil
		}
	}
	venue.ZoneIDs = append(venue.ZoneIDs, zoneID)
	return s.store.Save(store.VenueStoreKey(venueID), venue)
}
