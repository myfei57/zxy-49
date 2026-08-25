package screen

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) Reset(id string, at int64) error {
	if !s.store.Exists(store.ScreenStoreKey(id)) {
		return errors.New("screen not found: " + id)
	}
	return s.store.Save("screen-reset-"+id, domain.Handover{ScreenID: id, Emergency: false, At: at})
}

func (s *Service) ResetAt(id string) (domain.Handover, error) {
	var reset domain.Handover
	err := s.store.Load("screen-reset-"+id, &reset)
	if errors.Is(err, store.ErrNotFound) {
		return domain.Handover{ScreenID: id}, nil
	}
	return reset, err
}

func (s *Service) ResetCount() (int, error) {
	keys, err := s.store.Keys("screen-reset-")
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}
