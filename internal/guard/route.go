package guard

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) RebuildRoute(routeID string, zoneID string, name string) error {
	if s.store.Exists(store.RouteStoreKey(routeID)) {
		return errors.New("route already exists: " + routeID)
	}
	checkpoints, err := s.buildCheckpoints(zoneID)
	if err != nil {
		return err
	}
	route := domain.Route{
		ID:          routeID,
		ZoneID:      zoneID,
		Name:        name,
		Checkpoints: checkpoints,
	}
	return s.store.Save(store.RouteStoreKey(routeID), route)
}

func (s *Service) buildCheckpoints(zoneID string) ([]domain.Checkpoint, error) {
	hallID, err := s.zoneSvc.HallOfZone(zoneID)
	if err != nil {
		return nil, err
	}
	booths, err := s.hallSvc.ListBooths(hallID)
	if err != nil {
		return nil, err
	}
	checkpoints := make([]domain.Checkpoint, 0)
	for index, booth := range booths {
		if booth.ZoneID != zoneID {
			continue
		}
		checkpoints = append(checkpoints, domain.Checkpoint{
			BoothID: booth.ID,
			ZoneID:  zoneID,
			Order:   index,
		})
	}
	return checkpoints, nil
}

func (s *Service) RebuildExisting(routeID string) error {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return err
	}
	checkpoints, err := s.buildCheckpoints(route.ZoneID)
	if err != nil {
		return err
	}
	route.Checkpoints = checkpoints
	return s.store.Save(store.RouteStoreKey(routeID), route)
}

func (s *Service) RouteAlarm(zoneID string) (string, error) {
	if partition, ok := s.partitions[zoneID]; ok {
		return partition, nil
	}
	partition, err := s.zoneSvc.PartitionOf(zoneID)
	if err != nil {
		return "", err
	}
	return partition, nil
}

func (s *Service) CheckpointCount(routeID string) (int, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return 0, err
	}
	return len(route.Checkpoints), nil
}

func (s *Service) Checkpoints(routeID string) ([]domain.Checkpoint, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return nil, err
	}
	return route.Checkpoints, nil
}

func (s *Service) BoothIDs(routeID string) ([]string, error) {
	route, err := s.GetRoute(routeID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(route.Checkpoints))
	for _, checkpoint := range route.Checkpoints {
		ids = append(ids, checkpoint.BoothID)
	}
	return ids, nil
}
