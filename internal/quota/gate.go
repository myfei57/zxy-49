package quota

import (
	"errors"
	"fmt"

	"venueops/internal/store"
)

func (s *Service) Reserve(hallID string, watts int) error {
	if watts < 0 {
		return errors.New("watts cannot be negative")
	}
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return err
	}
	if quota.UsedWatts+watts > quota.LimitWatts {
		return fmt.Errorf("quota exceeded for hall %s: used %d limit %d", hallID, quota.UsedWatts, quota.LimitWatts)
	}
	quota.UsedWatts += watts
	return s.store.Save(store.QuotaStoreKey(hallID), quota)
}

func (s *Service) Release(hallID string, watts int) error {
	if watts < 0 {
		return errors.New("watts cannot be negative")
	}
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return err
	}
	if watts > quota.UsedWatts {
		quota.UsedWatts = 0
	} else {
		quota.UsedWatts -= watts
	}
	return s.store.Save(store.QuotaStoreKey(hallID), quota)
}

func (s *Service) CanReserve(hallID string, watts int) (bool, error) {
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return false, err
	}
	return quota.UsedWatts+watts <= quota.LimitWatts, nil
}

func (s *Service) Available(hallID string) (int, error) {
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return 0, err
	}
	available := quota.LimitWatts - quota.UsedWatts
	if available < 0 {
		return 0, nil
	}
	return available, nil
}

func (s *Service) Usage(hallID string) (int, int, error) {
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return 0, 0, err
	}
	return quota.UsedWatts, quota.LimitWatts, nil
}
