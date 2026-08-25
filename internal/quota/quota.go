package quota

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

func (s *Service) SetQuota(id string, hallID string, limitWatts int) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("quota id is required")
	}
	if strings.TrimSpace(hallID) == "" {
		return errors.New("hall id is required")
	}
	if limitWatts < 0 {
		return errors.New("quota limit cannot be negative")
	}
	quota := domain.Quota{
		ID:         id,
		HallID:     hallID,
		LimitWatts: limitWatts,
		UsedWatts:  0,
	}
	return s.store.Save(store.QuotaStoreKey(hallID), quota)
}

func (s *Service) GetQuota(hallID string) (domain.Quota, error) {
	var quota domain.Quota
	err := s.store.Load(store.QuotaStoreKey(hallID), &quota)
	return quota, err
}

func (s *Service) HasQuota(hallID string) bool {
	return s.store.Exists(store.QuotaStoreKey(hallID))
}

func (s *Service) ListQuotas() ([]domain.Quota, error) {
	keys, err := s.store.Keys(store.QuotaKey)
	if err != nil {
		return nil, err
	}
	quotas := make([]domain.Quota, 0, len(keys))
	for _, key := range keys {
		var quota domain.Quota
		if err := s.store.Load(key, &quota); err != nil {
			return nil, err
		}
		quotas = append(quotas, quota)
	}
	return quotas, nil
}

func (s *Service) ChangeLimit(hallID string, limitWatts int) error {
	if limitWatts < 0 {
		return errors.New("quota limit cannot be negative")
	}
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return err
	}
	quota.LimitWatts = limitWatts
	return s.store.Save(store.QuotaStoreKey(hallID), quota)
}

func (s *Service) ResetUsage(hallID string) error {
	quota, err := s.GetQuota(hallID)
	if err != nil {
		return err
	}
	quota.UsedWatts = 0
	return s.store.Save(store.QuotaStoreKey(hallID), quota)
}
