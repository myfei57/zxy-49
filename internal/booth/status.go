package booth

import "venueops/internal/store"

func (s *Service) PoweredCount() (int, error) {
	booths, err := s.List()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, booth := range booths {
		if booth.Powered {
			count++
		}
	}
	return count, nil
}

func (s *Service) TotalLoad() (int, error) {
	booths, err := s.List()
	if err != nil {
		return 0, err
	}
	total := 0
	for _, booth := range booths {
		total += booth.LoadWatts
	}
	return total, nil
}

func (s *Service) PoweredLoad() (int, error) {
	booths, err := s.List()
	if err != nil {
		return 0, err
	}
	total := 0
	for _, booth := range booths {
		if booth.Powered {
			total += booth.LoadWatts
		}
	}
	return total, nil
}

func (s *Service) HallLoad(hallID string) (int, error) {
	booths, err := s.BoothsOfHall(hallID)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, booth := range booths {
		if booth.Powered {
			total += booth.LoadWatts
		}
	}
	return total, nil
}

func (s *Service) BoothIDs() ([]string, error) {
	booths, err := s.List()
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(booths))
	for _, booth := range booths {
		ids = append(ids, booth.ID)
	}
	return ids, nil
}

func (s *Service) BoothCount() (int, error) {
	return s.Count()
}

func (s *Service) Store() *store.Store {
	return s.store
}
