package hall

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) LayoutChange(hallID string, positions map[string]domain.Point) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	if positions == nil {
		return errors.New("positions map is required")
	}
	hall.Layout = positions
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) BoothPosition(hallID string, boothID string) (domain.Point, bool) {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return domain.Point{}, false
	}
	position, ok := hall.Layout[boothID]
	return position, ok
}

func (s *Service) SetBoothPosition(hallID string, boothID string, position domain.Point) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	if hall.Layout == nil {
		hall.Layout = map[string]domain.Point{}
	}
	hall.Layout[boothID] = position
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) RemoveBoothPosition(hallID string, boothID string) error {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return err
	}
	delete(hall.Layout, boothID)
	return s.store.Save(store.HallStoreKey(hallID), hall)
}

func (s *Service) LayoutBooths(hallID string) ([]string, error) {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return nil, err
	}
	boothIDs := make([]string, 0, len(hall.Layout))
	for boothID := range hall.Layout {
		boothIDs = append(boothIDs, boothID)
	}
	return boothIDs, nil
}

func (s *Service) LayoutSize(hallID string) (int, error) {
	hall, err := s.GetHall(hallID)
	if err != nil {
		return 0, err
	}
	return len(hall.Layout), nil
}

func (s *Service) MoveBooth(hallID string, boothID string, position domain.Point) error {
	return s.SetBoothPosition(hallID, boothID, position)
}
