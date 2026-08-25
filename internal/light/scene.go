package light

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) UpdateScene(scene domain.Scene, at int64) error {
	if err := s.validateZone(scene.ZoneID); err != nil {
		return err
	}
	if scene.Brightness < 0 || scene.Brightness > 100 {
		return errors.New("brightness must be between 0 and 100")
	}
	scene.IsEmergency = false
	if err := s.store.Save(store.SceneStoreKey(scene.ZoneID), scene); err != nil {
		return err
	}
	if err := s.store.Save(store.EffectiveSceneKeyFor(scene.ZoneID), scene); err != nil {
		return err
	}
	return s.audit.Record("light.scene", scene.ZoneID, scene.Name, at)
}

func (s *Service) SceneFor(zoneID string) (domain.Scene, error) {
	var scene domain.Scene
	err := s.store.Load(store.SceneStoreKey(zoneID), &scene)
	return scene, err
}

func (s *Service) EffectiveScene(zoneID string) (domain.Scene, error) {
	var scene domain.Scene
	err := s.store.Load(store.EffectiveSceneKeyFor(zoneID), &scene)
	return scene, err
}

func (s *Service) ApplyScheduled(zoneID string, at int64) error {
	if err := s.validateZone(zoneID); err != nil {
		return err
	}
	hallID, err := s.hallOfZone(zoneID)
	if err != nil {
		return err
	}
	active, err := s.EvacuationActive(hallID)
	if err != nil {
		return err
	}
	if active {
		return nil
	}
	scene, err := s.SceneFor(zoneID)
	if err != nil {
		return err
	}
	if err := s.store.Save(store.EffectiveSceneKeyFor(zoneID), scene); err != nil {
		return err
	}
	return nil
}

func (s *Service) SceneNames() ([]string, error) {
	keys, err := s.store.Keys(store.SceneKey)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(keys))
	for _, key := range keys {
		var scene domain.Scene
		if err := s.store.Load(key, &scene); err != nil {
			return nil, err
		}
		names = append(names, scene.Name)
	}
	return names, nil
}

func (s *Service) Scenes() ([]domain.Scene, error) {
	keys, err := s.store.Keys(store.SceneKey)
	if err != nil {
		return nil, err
	}
	scenes := make([]domain.Scene, 0, len(keys))
	for _, key := range keys {
		var scene domain.Scene
		if err := s.store.Load(key, &scene); err != nil {
			return nil, err
		}
		scenes = append(scenes, scene)
	}
	return scenes, nil
}

func (s *Service) SceneCount() (int, error) {
	keys, err := s.store.Keys(store.SceneKey)
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}
