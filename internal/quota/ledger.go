package quota

import (
	"errors"

	"venueops/internal/domain"
	"venueops/internal/store"
)

type Ledger struct {
	Entries []domain.Quota `json:"entries"`
}

func (s *Service) Ledger(hallID string) (Ledger, error) {
	var ledger Ledger
	if err := s.store.Load("quota-ledger-"+hallID, &ledger); err != nil {
		return Ledger{}, err
	}
	return ledger, nil
}

func (s *Service) AppendLedger(hallID string, quota domain.Quota) error {
	ledger, err := s.Ledger(hallID)
	if err != nil {
		if !errors.Is(err, store.ErrNotFound) {
			return err
		}
		ledger = Ledger{}
	}
	ledger.Entries = append(ledger.Entries, quota)
	return s.store.Save("quota-ledger-"+hallID, ledger)
}

func (s *Service) LedgerSize(hallID string) (int, error) {
	ledger, err := s.Ledger(hallID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return 0, nil
		}
		return 0, err
	}
	return len(ledger.Entries), nil
}
