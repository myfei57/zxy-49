package hall

import (
	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) BoothChange(hallID string, boothID string, zoneID string, name string, add bool) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	if add {
		if _, exists := hall.Layout[boothID]; !exists {
			if hall.Layout == nil {
				hall.Layout = map[string]domain.Point{}
			}
			hall.Layout[boothID] = domain.Point{X: len(hall.Layout) * 10, Y: 20}
		}
		booth := domain.Booth{
			ID:     boothID,
			HallID: hallID,
			ZoneID: zoneID,
			Name:   name,
		}
		if err := s.store.Save(store.BoothStoreKey(boothID), booth); err != nil {
			return err
		}
		return s.store.Save(store.HallStoreKey(hallID), hall)
	}
	delete(hall.Layout, boothID)
	if err := s.store.Delete(store.BoothStoreKey(boothID)); err != nil {
		return err
	}
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) ListBooths(hallID string) ([]domain.Booth, error) {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return nil, err
	}
	booths := make([]domain.Booth, 0, len(hall.Layout))
	for boothID := range hall.Layout {
		var booth domain.Booth
		if err := s.store.Load(store.BoothStoreKey(boothID), &booth); err != nil {
			return nil, err
		}
		booths = append(booths, booth)
	}
	return booths, nil
}

func (s *Service) BoothCount(hallID string) (int, error) {
	booths, err := s.ListBooths(hallID)
	if err != nil {
		return 0, err
	}
	return len(booths), nil
}

func (s *Service) BoothExists(hallID string, boothID string) bool {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return false
	}
	_, ok := hall.Layout[boothID]
	return ok
}

func (s *Service) BoothZones(hallID string) ([]string, error) {
	booths, err := s.ListBooths(hallID)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	zones := make([]string, 0)
	for _, booth := range booths {
		if !seen[booth.ZoneID] {
			seen[booth.ZoneID] = true
			zones = append(zones, booth.ZoneID)
		}
	}
	return zones, nil
}
