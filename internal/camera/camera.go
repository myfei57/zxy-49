package camera

import (
	"errors"
	"strings"

	"venueops/internal/domain"
	"venueops/internal/hall"
	"venueops/internal/store"
)

type Service struct {
	store *store.Store
	hall  *hall.Service
	season store.Season
}

func NewService(st *store.Store, hallSvc *hall.Service) *Service {
	svc := &Service{store: st, hall: hallSvc}
	if season, err := st.LoadSeason(); err == nil {
		svc.season = season
	}
	return svc
}

func (s *Service) Create(id string, hallID string, name string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("camera id is required")
	}
	if strings.TrimSpace(hallID) == "" {
		return errors.New("hall id is required")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("camera name is required")
	}
	if s.store.Exists(store.CameraStoreKey(id)) {
		return errors.New("camera already exists: " + id)
	}
	camera := domain.Camera{
		ID:       id,
		HallID:   hallID,
		Name:     name,
		Presets:  map[string]domain.Preset{},
		DayNight: "day",
	}
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) Get(id string) (domain.Camera, error) {
	var camera domain.Camera
	err := s.store.Load(store.CameraStoreKey(id), &camera)
	return camera, err
}

func (s *Service) List() ([]domain.Camera, error) {
	keys, err := s.store.Keys(store.CameraKey)
	if err != nil {
		return nil, err
	}
	cameras := make([]domain.Camera, 0, len(keys))
	for _, key := range keys {
		var camera domain.Camera
		if err := s.store.Load(key, &camera); err != nil {
			return nil, err
		}
		cameras = append(cameras, camera)
	}
	return cameras, nil
}

func (s *Service) Rename(id string, name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("camera name is required")
	}
	camera, err := s.Get(id)
	if err != nil {
		return err
	}
	camera.Name = name
	return s.store.Save(store.CameraStoreKey(id), camera)
}

func (s *Service) Count() (int, error) {
	keys, err := s.store.Keys(store.CameraKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

func (s *Service) CamerasOfHall(hallID string) ([]domain.Camera, error) {
	cameras, err := s.List()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.Camera, 0)
	for _, camera := range cameras {
		if camera.HallID == hallID {
			filtered = append(filtered, camera)
		}
	}
	return filtered, nil
}

func (s *Service) HasCamera(id string) bool {
	return s.store.Exists(store.CameraStoreKey(id))
}

func (s *Service) HallOf(id string) (string, error) {
	camera, err := s.Get(id)
	if err != nil {
		return "", err
	}
	return camera.HallID, nil
}
