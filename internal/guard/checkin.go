package guard

import (
	"errors"

	"venueops/internal/store"
)

func (s *Service) CheckIn(routeID string, boothID string, at int64) error {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return err
	}
	found := false
	for index := range route.Checkpoints {
		if route.Checkpoints[index].BoothID == boothID {
			route.Checkpoints[index].Visited = true
			found = true
			break
		}
	}
	if !found {
		return errors.New("booth is not a checkpoint of route: " + boothID)
	}
	if err := s.store.Save(store.RouteStoreKey(routeID), route); err != nil {
		return err
	}
	return s.audit.Record("guard.checkin", routeID, boothID, at)
}

func (s *Service) CheckedIn(routeID string, boothID string) (bool, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return false, err
	}
	for _, checkpoint := range route.Checkpoints {
		if checkpoint.BoothID == boothID {
			return checkpoint.Visited, nil
		}
	}
	return false, errors.New("booth is not a checkpoint of route: " + boothID)
}

func (s *Service) Overdue(routeID string, now int64, limitSec int64) (bool, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return false, err
	}
	if len(route.Checkpoints) == 0 {
		return false, nil
	}
	visitedCount := 0
	for _, checkpoint := range route.Checkpoints {
		if checkpoint.Visited {
			visitedCount++
		}
	}
	if visitedCount == len(route.Checkpoints) {
		return false, nil
	}
	lastAt, err := s.lastCheckinAt(routeID)
	if err != nil {
		return false, err
	}
	if lastAt == 0 {
		return true, nil
	}
	return now-lastAt > limitSec, nil
}

func (s *Service) lastCheckinAt(routeID string) (int64, error) {
	records, err := s.audit.ByAction("guard.checkin")
	if err != nil {
		return 0, err
	}
	var last int64
	for _, record := range records {
		if record.Target == routeID && record.At > last {
			last = record.At
		}
	}
	return last, nil
}

func (s *Service) Progress(routeID string) (int, int, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return 0, 0, err
	}
	visited := 0
	for _, checkpoint := range route.Checkpoints {
		if checkpoint.Visited {
			visited++
		}
	}
	return visited, len(route.Checkpoints), nil
}
