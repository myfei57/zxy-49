package store

type Season struct {
	Name      string `json:"name"`
	NightHour int    `json:"night_hour"`
	NightMin  int    `json:"night_min"`
}

func (s *Store) SaveSeason(season Season) error {
	return s.Save("season", season)
}

func (s *Store) LoadSeason() (Season, error) {
	var season Season
	err := s.Load("season", &season)
	return season, err
}

func (s *Store) SeasonExists() bool {
	return s.Exists("season")
}
