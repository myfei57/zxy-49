package screen

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) HandoverToEmergency(screenID string, emergency bool) error {
	if !s.store.Exists(store.ScreenStoreKey(screenID)) {
		return errors.New("screen not found: " + screenID)
	}
	return s.store.Save(store.HandoverKeyFor(screenID), domain.Handover{ScreenID: screenID, Emergency: emergency})
}

func (s *Service) HandoverState(screenID string) (bool, error) {
	var handover domain.Handover
	if err := s.store.Load(store.HandoverKeyFor(screenID), &handover); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return handover.Emergency, nil
}

func (s *Service) HandoverScreens() ([]string, error) {
	keys, err := s.store.Keys(store.HandoverKey)
	if err != nil {
		return nil, err
	}
	screens := make([]string, 0)
	for _, key := range keys {
		var handover domain.Handover
		if err := s.store.Load(key, &handover); err != nil {
			return nil, err
		}
		if handover.Emergency {
			screens = append(screens, handover.ScreenID)
		}
	}
	return screens, nil
}

func (s *Service) ClearHandover(screenID string) error {
	return s.store.Delete(store.HandoverKeyFor(screenID))
}
