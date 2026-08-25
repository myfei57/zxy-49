package gate

import (
	"errors"

	"venueops/internal/store"
)

func (s *Service) Release(id string) error {
	gate, err := s.Get(id)
	if err != nil {
		return err
	}
	if gate.Mode != "pass" {
		return errors.New("gate cannot release unless mode is pass")
	}
	gate.Released = true
	if err := s.store.Save(store.GateStoreKey(id), gate); err != nil {
		return err
	}
	return s.audit.Record("gate.release", id, "released", 0)
}

func (s *Service) Released(id string) (bool, error) {
	gate, err := s.Get(id)
	if err != nil {
		return false, err
	}
	return gate.Released, nil
}

func (s *Service) Disarm(id string) error {
	return s.SetMode(id, "disarmed", 0)
}

func (s *Service) ArmedGates() ([]string, error) {
	gates, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, gate := range gates {
		if gate.Mode == "armed" {
			ids = append(ids, gate.ID)
		}
	}
	return ids, nil
}

func (s *Service) ReleasedGates() ([]string, error) {
	gates, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, gate := range gates {
		if gate.Released {
			ids = append(ids, gate.ID)
		}
	}
	return ids, nil
}
