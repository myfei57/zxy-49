package light

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) Evacuate(hallID string, at int64) error {
	if err := s.validateHall(hallID); err != nil {
		return err
	}
	if err := s.store.Save(store.EvacuationKeyFor(hallID), domain.Evacuation{HallID: hallID, Active: true, At: at}); err != nil {
		return err
	}
	zones, err := s.zonesOfHall(hallID)
	if err != nil {
		return err
	}
	for _, zoneID := range zones {
		scene := domain.Scene{
			ZoneID:      zoneID,
			Name:        "emergency",
			Brightness:  100,
			Color:       "white",
			IsEmergency: true,
		}
		if err := s.store.Save(store.EffectiveSceneKeyFor(zoneID), scene); err != nil {
			return err
		}
	}
	if s.handover != nil {
		for _, screenID := range s.screenIDs {
			if err := s.handover(screenID, true); err != nil {
				return err
			}
		}
	}
	return s.audit.Record("light.evacuate", hallID, "active", at)
}

func (s *Service) ClearEvacuation(hallID string, at int64) error {
	if err := s.validateHall(hallID); err != nil {
		return err
	}
	if err := s.store.Save(store.EvacuationKeyFor(hallID), domain.Evacuation{HallID: hallID, Active: false, At: at}); err != nil {
		return err
	}
	if s.handover != nil {
		for _, screenID := range s.screenIDs {
			if err := s.handover(screenID, false); err != nil {
				return err
			}
		}
	}
	return s.audit.Record("light.evacuate", hallID, "cleared", at)
}

func (s *Service) EvacuationActive(hallID string) (bool, error) {
	var evacuation domain.Evacuation
	if err := s.store.Load(store.EvacuationKeyFor(hallID), &evacuation); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return evacuation.Active, nil
}

func (s *Service) EvacuationState(hallID string) (domain.Evacuation, error) {
	var evacuation domain.Evacuation
	err := s.store.Load(store.EvacuationKeyFor(hallID), &evacuation)
	if errors.Is(err, store.ErrNotFound) {
		return domain.Evacuation{HallID: hallID, Active: false}, nil
	}
	return evacuation, err
}

func (s *Service) EvacuatingHalls() ([]string, error) {
	keys, err := s.store.Keys(store.EvacuationKey)
	if err != nil {
		return nil, err
	}
	halls := make([]string, 0)
	for _, key := range keys {
		var evacuation domain.Evacuation
		if err := s.store.Load(key, &evacuation); err != nil {
			return nil, err
		}
		if evacuation.Active {
			halls = append(halls, evacuation.HallID)
		}
	}
	return halls, nil
}
