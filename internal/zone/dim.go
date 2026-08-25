package zone

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

type SceneProvider func(zoneID string) (domain.Scene, error)

func (s *Service) Dim(zoneID string, provider SceneProvider) (domain.Scene, error) {
	if provider == nil {
		return domain.Scene{}, errors.New("scene provider is required")
	}
	scene, err := provider(zoneID)
	if err != nil {
		return domain.Scene{}, err
	}
	if scene.IsEmergency {
		scene.Brightness = 100
		scene.Color = "white"
	}
	if err := s.store.Save(store.DimKeyFor(zoneID), scene); err != nil {
		return domain.Scene{}, err
	}
	return scene, nil
}

func (s *Service) DimLevel(zoneID string) (domain.Scene, error) {
	var scene domain.Scene
	err := s.store.Load(store.DimKeyFor(zoneID), &scene)
	return scene, err
}

func (s *Service) SetDimLevel(zoneID string, brightness int) error {
	scene, err := s.DimLevel(zoneID)
	if err != nil {
		return err
	}
	scene.Brightness = brightness
	return s.store.Save(store.DimKeyFor(zoneID), scene)
}

func (s *Service) DimAll(hallID string, provider SceneProvider) (map[string]domain.Scene, error) {
	zones, err := s.ZonesOfHall(hallID)
	if err != nil {
		return nil, err
	}
	results := map[string]domain.Scene{}
	for _, zone := range zones {
		scene, err := s.Dim(zone.ID, provider)
		if err != nil {
			return nil, err
		}
		results[zone.ID] = scene
	}
	return results, nil
}

func (s *Service) ClearDim(zoneID string) error {
	return s.store.Delete(store.DimKeyFor(zoneID))
}
