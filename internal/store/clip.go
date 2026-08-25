package store

import "venueops/internal/domain"

type clipSet struct {
	Clips []domain.Clip `json:"clips"`
}

func (s *Store) SaveClips(screenID string, clips []domain.Clip) error {
	return s.Save("clips-"+screenID, clipSet{Clips: clips})
}

func (s *Store) LoadClips(screenID string) ([]domain.Clip, error) {
	var set clipSet
	if err := s.Load("clips-"+screenID, &set); err != nil {
		return nil, err
	}
	return set.Clips, nil
}

func (s *Store) ClipsExist(screenID string) bool {
	return s.Exists("clips-" + screenID)
}
