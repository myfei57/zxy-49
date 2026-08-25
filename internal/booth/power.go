package booth

import (
	"venueops/internal/store"
)

func (s *Service) PowerOn(id string) error {
	booth, err := s.Get(id)
	if err != nil {
		return err
	}
	if booth.Powered {
		return nil
	}
	if err := s.quotaS.Reserve(booth.HallID, booth.LoadWatts); err != nil {
		return err
	}
	booth.Powered = true
	return s.store.Save(store.BoothStoreKey(id), booth)
}

func (s *Service) PowerOff(id string) error {
	booth, err := s.Get(id)
	if err != nil {
		return err
	}
	if !booth.Powered {
		return nil
	}
	booth.Powered = false
	if err := s.store.Save(store.BoothStoreKey(id), booth); err != nil {
		return err
	}
	return s.quotaS.Release(booth.HallID, booth.LoadWatts)
}

func (s *Service) PowerToggle(id string) (bool, error) {
	booth, err := s.Get(id)
	if err != nil {
		return false, err
	}
	if booth.Powered {
		return false, s.PowerOff(id)
	}
	if err := s.PowerOn(id); err != nil {
		return false, err
	}
	return true, nil
}

func (s *Service) Powered(id string) (bool, error) {
	booth, err := s.Get(id)
	if err != nil {
		return false, err
	}
	return booth.Powered, nil
}

func (s *Service) Load(id string) (int, error) {
	booth, err := s.Get(id)
	if err != nil {
		return 0, err
	}
	return booth.LoadWatts, nil
}

func (s *Service) PowerOffAll(hallID string) (int, error) {
	booths, err := s.BoothsOfHall(hallID)
	if err != nil {
		return 0, err
	}
	off := 0
	for _, booth := range booths {
		if booth.Powered {
			if err := s.PowerOff(booth.ID); err != nil {
				return off, err
			}
			off++
		}
	}
	return off, nil
}

func (s *Service) PowerOnAll(hallID string) (int, error) {
	booths, err := s.BoothsOfHall(hallID)
	if err != nil {
		return 0, err
	}
	on := 0
	for _, booth := range booths {
		if !booth.Powered {
			if err := s.PowerOn(booth.ID); err != nil {
				return on, err
			}
			on++
		}
	}
	return on, nil
}
