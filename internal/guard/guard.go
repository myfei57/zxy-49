package guard

import (
	"errors"
	"strings"

	"venueops/internal/audit"
	"venueops/internal/domain"
	"venueops/internal/hall"
	"venueops/internal/store"
	"venueops/internal/zone"
)

type Service struct {
	store   *store.Store
	zoneSvc *zone.Service
	hallSvc *hall.Service
	audit   *audit.Service
}

func NewService(st *store.Store, zoneSvc *zone.Service, hallSvc *hall.Service, auditSvc *audit.Service) *Service {
	return &Service{store: st, zoneSvc: zoneSvc, hallSvc: hallSvc, audit: auditSvc}
}

func (s *Service) Create(routeID string, zoneID string, name string) error {
	if strings.TrimSpace(routeID) == "" {
		return errors.New("route id is required")
	}
	if strings.TrimSpace(zoneID) == "" {
		return errors.New("zone id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("route name is required")
	}
	if s.store.Exists(store.RouteStoreKey(routeID)) {
		return errors.New("route already exists: " + routeID)
	}
	route := domain.Route{
		ID:          routeID,
		ZoneID:      zoneID,
		Name:        name,
		Checkpoints: []domain.Checkpoint{},
	}
	return s.store.Save(store.RouteStoreKey(routeID), route)
}

func (s *Service) GetRoute(routeID string) (domain.Route, error) {
	var route domain.Route
	err := s.store.Load(store.RouteStoreKey(routeID), &route)
	return route, err
}

func (s *Service) ListRoutes() ([]domain.Route, error) {
	keys, err := s.store.Keys(store.RouteKey)
	if err != nil {
		return nil, err
	}
	routes := make([]domain.Route, 0, len(keys))
	for _, key := range keys {
		var route domain.Route
		if err := s.store.Load(key, &route); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (s *Service) Rename(routeID string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("route name is required")
	}
	route, err := s.GetRoute(routeID)
	if err != nil {
		return err
	}
	route.Name = name
	return s.store.Save(store.RouteStoreKey(routeID), route)
}

func (s *Service) Count() (int, error) {
	keys, err := s.store.Keys(store.RouteKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) HasRoute(routeID string) bool {
	return s.store.Exists(store.RouteStoreKey(routeID))
}

func (s *Service) RoutesOfZone(zoneID string) ([]domain.Route, error) {
	routes, err := s.ListRoutes()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Route, 0)
	for _, route := range routes {
		if route.ZoneID == zoneID {
			filtered = append(filtered, route)
		}
	}
	return filtered, nil
}
