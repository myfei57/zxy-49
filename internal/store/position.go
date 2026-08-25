package store

type PlaybackPosition struct {
	Position int `json:"position"`
	Epoch    int `json:"epoch"`
}

func (s *Store) SavePosition(screenID string, position int, epoch int) error {
	return s.Save("position-"+screenID, PlaybackPosition{Position: position, Epoch: epoch})
}

func (s *Store) LoadPosition(screenID string) (PlaybackPosition, error) {
	var position PlaybackPosition
	err := s.Load("position-"+screenID, &position)
	return position, err
}

func (s *Store) ClearPosition(screenID string) error {
	return s.Delete("position-" + screenID)
}
