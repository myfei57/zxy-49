package camera

import (
	"errors"
	"sort"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) SetPreset(id string, boothID string, pan int, tilt int) error {
	camera, err := s.Get(id)
	if err != nil {
		return err
	}
	position, ok := s.hall.BoothPosition(camera.HallID, boothID)
	if !ok {
		return errors.New("booth not found in hall layout: " + boothID)
	}
	if camera.Presets == nil {
		camera.Presets = map[string]domain.Preset{}
	}
	camera.Presets[boothID] = domain.Preset{
		BoothID: boothID,
		X:       position.X,
		Y:       position.Y,
		Pan:     pan,
		Tilt:    tilt,
	}
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) RefreshPresets(id string, at int64) error {
	camera, err := s.Get(id)
	if err != nil {
		return err
	}
	boothIDs, err := s.hall.LayoutBooths(camera.HallID)
	if err != nil {
		return err
	}
	if camera.Presets == nil {
		camera.Presets = map[string]domain.Preset{}
	}
	for _, boothID := range boothIDs {
		position, ok := s.hall.BoothPosition(camera.HallID, boothID)
		if !ok {
			continue
		}
		preset := camera.Presets[boothID]
		preset.BoothID = boothID
		preset.X = position.X
		preset.Y = position.Y
		camera.Presets[boothID] = preset
	}
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) PresetFor(id string, boothID string) (domain.Preset, error) {
	camera, err := s.Get(id)
	if err != nil {
		return domain.Preset{}, err
	}
	preset, ok := camera.Presets[boothID]
	if !ok {
		return domain.Preset{}, errors.New("preset not found for booth: " + boothID)
	}
	return preset, nil
}

func (s *Service) PresetCount(id string) (int, error) {
	camera, err := s.Get(id)
	if err != nil {
		return 0, err
	}
	return len(camera.Presets), nil
}

func (s *Service) PresetBooths(id string) ([]string, error) {
	camera, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	boothIDs := make([]string, 0, len(camera.Presets))
	for boothID := range camera.Presets {
		boothIDs = append(boothIDs, boothID)
	}
	sort.Strings(boothIDs)
	return boothIDs, nil
}

func (s *Service) RemovePreset(id string, boothID string) error {
	camera, err := s.Get(id)
	if err != nil {
		return err
	}
	delete(camera.Presets, boothID)
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) PresetPosition(id string, boothID string) (domain.Point, error) {
	preset, err := s.PresetFor(id, boothID)
	if err != nil {
		return domain.Point{}, err
	}
	return domain.Point{X: preset.X, Y: preset.Y}, nil
}
