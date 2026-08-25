package gate

import (
	"errors"
	"strings"

	"venueops/internal/audit"
	"venueops/internal/domain"
	"venueops/internal/store"
)

type Service struct {
	store *store.Store
	audit *audit.Service
}

func NewService(st *store.Store, auditSvc *audit.Service) *Service {
	return &Service{store: st, audit: auditSvc}
}

func (s *Service) Create(id string, zoneID string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("gate id is required")
	}
	if strings.TrimSpace(zoneID) == "" {
		return errors.New("zone id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("gate name is required")
	}
	if s.store.Exists(store.GateStoreKey(id)) {
		return errors.New("gate already exists: " + id)
	}
	gate := domain.Gate{
		ID:       id,
		ZoneID:   zoneID,
		Name:     name,
		Mode:     "armed",
		Released: false,
	}
	if err := s.store.Save(store.GateStoreKey(id), gate); err != nil {
		return err
	}
	return s.audit.Record("gate.create", id, "armed", 0)
}

func (s *Service) Get(id string) (domain.Gate, error) {
	var gate domain.Gate
	err := s.store.Load(store.GateStoreKey(id), &gate)
	return gate, err
}

func (s *Service) List() ([]domain.Gate, error) {
	keys, err := s.store.Keys(store.GateKey)
	if err != nil {
		return nil, err
	}
	gates := make([]domain.Gate, 0, len(keys))
	for _, key := range keys {
		var gate domain.Gate
		if err := s.store.Load(key, &gate); err != nil {
			return nil, err
		}
		gates = append(gates, gate)
	}
	return gates, nil
}

func (s *Service) Rename(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("gate name is required")
	}
	gate, err := s.Get(id)
	if err != nil {
		return err
	}
	gate.Name = name
	return s.store.Save(store.GateStoreKey(id), gate)
}

func (s *Service) Count() (int, error) {
	keys, err := s.store.Keys(store.GateKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) GatesOfZone(zoneID string) ([]domain.Gate, error) {
	gates, err := s.List()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Gate, 0)
	for _, gate := range gates {
		if gate.ZoneID == zoneID {
			filtered = append(filtered, gate)
		}
	}
	return filtered, nil
}

func (s *Service) HasGate(id string) bool {
	return s.store.Exists(store.GateStoreKey(id))
}

func (s *Service) ZoneOf(id string) (string, error) {
	gate, err := s.Get(id)
	if err != nil {
		return "", err
	}
	return gate.ZoneID, nil
}
