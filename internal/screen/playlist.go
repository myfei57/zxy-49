package screen

import (
	"errors"
	"time"

	"venueops/internal/domain"
	"venueops/internal/store"
)

func (s *Service) SetPlaylist(id string, clips []domain.Clip) error {
	if !s.store.Exists(store.ScreenStoreKey(id)) {
		return errors.New("screen not found: " + id)
	}
	if len(clips) == 0 {
		return errors.New("playlist cannot be empty")
	}
	if err := s.store.SaveClips(id, clips); err != nil {
		return err
	}
	return s.store.SavePosition(id, 0, 1)
}

func (s *Service) Tick(id string, now time.Time) (domain.Clip, error) {
	clips, err := s.store.LoadClips(id)
	if err != nil {
		return domain.Clip{}, err
	}
	if len(clips) == 0 {
		return domain.Clip{}, errors.New("playlist is empty for screen: " + id)
	}
	position, err := s.store.LoadPosition(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			position = store.PlaybackPosition{Position: 0, Epoch: 1}
		} else {
			return domain.Clip{}, err
		}
	}
	next := now.Minute() % len(clips)
	epoch := position.Epoch + 1
	if err := s.store.SavePosition(id, next, epoch); err != nil {
		return domain.Clip{}, err
	}
	return clips[next], nil
}

func (s *Service) Current(id string) (domain.Clip, error) {
	clips, err := s.store.LoadClips(id)
	if err != nil {
		return domain.Clip{}, err
	}
	if len(clips) == 0 {
		return domain.Clip{}, errors.New("playlist is empty for screen: " + id)
	}
	position, err := s.store.LoadPosition(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			position = store.PlaybackPosition{Position: 0, Epoch: 1}
		} else {
			return domain.Clip{}, err
		}
	}
	if position.Position < 0 || position.Position >= len(clips) {
		position.Position = 0
	}
	return clips[position.Position], nil
}

func (s *Service) Position(id string) (int, error) {
	position, err := s.store.LoadPosition(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return position.Position, nil
}

func (s *Service) PlaylistSize(id string) (int, error) {
	clips, err := s.store.LoadClips(id)
	if err != nil {
		return 0, err
	}
	return len(clips), nil
}

func (s *Service) Clips(id string) ([]domain.Clip, error) {
	return s.store.LoadClips(id)
}

func (s *Service) Sponsors(id string) ([]string, error) {
	clips, err := s.store.LoadClips(id)
	if err != nil {
		return nil, err
	}
	sponsors := make([]string, 0, len(clips))
	for _, clip := range clips {
		sponsors = append(sponsors, clip.Sponsor)
	}
	return sponsors, nil
}
