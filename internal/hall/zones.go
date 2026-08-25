package hall

import (
	"errors"

	"venueops/internal/store"
)

func (s *Service) AddZone(hallID string, zoneID string) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	for _, existing := range hall.ZoneIDs {
		if existing == zoneID {
			return nil
		}
	}
	hall.ZoneIDs = append(hall.ZoneIDs, zoneID)
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) RemoveZone(hallID string, zoneID string) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	kept := make([]string, 0, len(hall.ZoneIDs))
	found := false
	for _, existing := range hall.ZoneIDs {
		if existing == zoneID {
			found = true
			continue
		}
		kept = append(kept, existing)
	}
	if !found {
		return errors.New("zone not attached to hall: " + zoneID)
	}
	hall.ZoneIDs = kept
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) ZonesOfHall(hallID string) ([]string, error) {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return nil, err
	}
	return hall.ZoneIDs, nil
}

func (s *Service) AttachVenueZone(hallID string, zoneID string) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	hall.ZoneIDs = append(hall.ZoneIDs, zoneID)
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) HallZoneIDs(hallID string) ([]string, error) {
	return s.ZonesOfHall(hallID)
}
