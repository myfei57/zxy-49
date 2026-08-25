package light

import (
	"errors"
	"strings"

	"venueops/internal/audit"
	"venueops/internal/domain"
	"venueops/internal/store"
)

type HandoverFunc func(screenID string, emergency bool) error

type Service struct {
	store     *store.Store
	audit     *audit.Service
	handover  HandoverFunc
	screenIDs []string
}

func NewService(st *store.Store, auditSvc *audit.Service) *Service {
	return &Service{store: st, audit: auditSvc}
}

func (s *Service) SetHandover(fn HandoverFunc) {
	s.handover = fn
}

func (s *Service) RegisterScreen(screenID string) {
	for _, existing := range s.screenIDs {
		if existing == screenID {
			return
		}
	}
	s.screenIDs = append(s.screenIDs, screenID)
}

func (s *Service) ScreenCount() int {
	return len(s.screenIDs)
}

func (s *Service) Screens() []string {
	return append([]string{}, s.screenIDs...)
}

func (s *Service) HandoverAll(emergency bool, at int64) error {
	if s.handover == nil {
		return errors.New("handover callback is not configured")
	}
	for _, screenID := range s.screenIDs {
		if err := s.handover(screenID, emergency); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) zoneFor(zoneID string) (domain.Zone, error) {
	var zone domain.Zone
	if err := s.store.Load(store.ZoneStoreKey(zoneID), &zone); err != nil {
		return domain.Zone{}, err
	}
	return zone, nil
}

func (s *Service) zonesOfHall(hallID string) ([]string, error) {
	keys, err := s.store.Keys(store.ZoneKey)
	if err != nil {
		return nil, err
	}
	zones := make([]string, 0)
	for _, key := range keys {
		var zone domain.Zone
		if err := s.store.Load(key, &zone); err != nil {
			return nil, err
		}
		if zone.HallID == hallID {
			zones = append(zones, zone.ID)
		}
	}
	return zones, nil
}

func (s *Service) hallOfZone(zoneID string) (string, error) {
	zone, err := s.zoneFor(zoneID)
	if err != nil {
		return "", err
	}
	return zone.HallID, nil
}

func (s *Service) validateZone(zoneID string) error {
	if strings.TrimSpace(zoneID) == "" {
		return errors.New("zone id is required")
	}
	if !s.store.Exists(store.ZoneStoreKey(zoneID)) {
		return errors.New("zone not found: " + zoneID)
	}
	return nil
}

func (s *Service) validateHall(hallID string) error {
	if strings.TrimSpace(hallID) == "" {
		return errors.New("hall id is required")
	}
	if !s.store.Exists(store.HallStoreKey(hallID)) {
		return errors.New("hall not found: " + hallID)
	}
	return nil
}
