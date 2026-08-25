package gate

import (
	"errors"

	"venueops/internal/store"
)

func validMode(mode string) bool {
	return mode == "armed" || mode == "disarmed" || mode == "pass"
}

func (s *Service) SetMode(id string, mode string, at int64) error {
	if !validMode(mode) {
		return errors.New("invalid gate mode: " + mode)
	}
	gate, err := s.Get(id)
	if err != nil {
		return err
	}
	if gate.Mode == mode {
		return nil
	}
	gate.Mode = mode
	gate.Released = mode == "pass"
	if err := s.store.Save(store.GateStoreKey(id), gate); err != nil {
		return err
	}
	if err := s.store.SaveGateMode(id, mode, at); err != nil {
		return err
	}
	return s.audit.Record("gate.mode", id, mode, at)
}

func (s *Service) Mode(id string) (string, error) {
	gate, err := s.Get(id)
	if err != nil {
		return "", err
	}
	return gate.Mode, nil
}

func (s *Service) DurableMode(id string) (string, error) {
	mode, err := s.store.LoadGateMode(id)
	if err != nil {
		return "", err
	}
	return mode.Mode, nil
}

func (s *Service) SetModeByName(name string, mode string, at int64) error {
	gates, err := s.List()
	if err != nil {
		return err
	}
	for _, gate := range gates {
		if gate.Name == name {
			return s.SetMode(gate.ID, mode, at)
		}
	}
	return errors.New("gate not found by name: " + name)
}

func (s *Service) ArmingState(id string) (string, error) {
	return s.Mode(id)
}
