package store

import "os"

type gateMode struct {
	Mode      string `json:"mode"`
	UpdatedAt int64  `json:"updated_at"`
}

func (s *Store) SaveGateMode(gateID string, mode string, at int64) error {
	return s.SaveIn("journal", "gate-mode-"+gateID, gateMode{Mode: mode, UpdatedAt: at})
}

func (s *Store) LoadGateMode(gateID string) (gateMode, error) {
	var mode gateMode
	err := s.LoadIn("journal", "gate-mode-"+gateID, &mode)
	return mode, err
}

func (s *Store) GateModeExists(gateID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.dir + "/journal/gate-mode-" + gateID + ".json"
	_, err := os.Stat(path)
	return err == nil
}
