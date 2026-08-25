package camera

import (
	"errors"

	"venueops/internal/store"
)

func (s *Service) DayNightSwitch(id string, hour int, minute int) (string, error) {
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return "", errors.New("invalid clock time")
	}
	camera, err := s.Get(id)
	if err != nil {
		return "", err
	}
	season := s.season
	current := hour*60 + minute
	threshold := season.NightHour*60 + season.NightMin
	mode := "day"
	if current >= threshold {
		mode = "night"
	}
	camera.DayNight = mode
	if err := s.store.Save(store.CameraStoreKey(id), camera); err != nil {
		return "", err
	}
	return mode, nil
}

func (s *Service) DayNightState(id string) (string, error) {
	camera, err := s.Get(id)
	if err != nil {
		return "", err
	}
	return camera.DayNight, nil
}

func (s *Service) SetDayNight(id string, mode string) error {
	if mode != "day" && mode != "night" {
		return errors.New("invalid day/night mode")
	}
	camera, err := s.Get(id)
	if err != nil {
		return err
	}
	camera.DayNight = mode
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) NightCameras() ([]string, error) {
	cameras, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, camera := range cameras {
		if camera.DayNight == "night" {
			ids = append(ids, camera.ID)
		}
	}
	return ids, nil
}

func (s *Service) DayCameras() ([]string, error) {
	cameras, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, camera := range cameras {
		if camera.DayNight == "day" {
			ids = append(ids, camera.ID)
		}
	}
	return ids, nil
}
