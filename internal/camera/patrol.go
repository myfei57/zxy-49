package camera

import (
	"errors"
	"sort"

	"venueops/internal/domain"
)

func (s *Service) Patrol(id string, at int64) (domain.Preset, error) {
	camera, err := s.Get(id)
	if err != nil {
		return domain.Preset{}, err
	}
	boothIDs := make([]string, 0, len(camera.Presets))
	for boothID := range camera.Presets {
		boothIDs = append(boothIDs, boothID)
	}
	if len(boothIDs) == 0 {
		return domain.Preset{}, errors.New("camera has no presets: " + id)
	}
	sort.Strings(boothIDs)
	index := int(at) % len(boothIDs)
	return camera.Presets[boothIDs[index]], nil
}

func (s *Service) PatrolTarget(id string, at int64) (string, error) {
	preset, err := s.Patrol(id, at)
	if err != nil {
		return "", err
	}
	return preset.BoothID, nil
}

func (s *Service) PatrolSequence(id string) ([]domain.Preset, error) {
	camera, err := s.Get(id)
	if err != nil {
		return nil, err
	}
	boothIDs := make([]string, 0, len(camera.Presets))
	for boothID := range camera.Presets {
		boothIDs = append(boothIDs, boothID)
	}
	sort.Strings(boothIDs)
	presets := make([]domain.Preset, 0, len(boothIDs))
	for _, boothID := range boothIDs {
		presets = append(presets, camera.Presets[boothID])
	}
	return presets, nil
}

func (s *Service) PatrolRounds(id string, count int) ([]domain.Preset, error) {
	sequence, err := s.PatrolSequence(id)
	if err != nil {
		return nil, err
	}
	if len(sequence) == 0 {
		return nil, errors.New("camera has no presets: " + id)
	}
	rounds := make([]domain.Preset, 0, count)
	for index := 0; index < count; index++ {
		rounds = append(rounds, sequence[index%len(sequence)])
	}
	return rounds, nil
}

func (s *Service) PatrolBooths(id string) ([]string, error) {
	return s.PresetBooths(id)
}
