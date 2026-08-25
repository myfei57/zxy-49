package audit

import (
	"errors"
	"strings"

	"venueops/internal/domain"
	"venueops/internal/store"

	"github.com/google/uuid"
)

type Service struct {
	store *store.Store
}

func NewService(st *store.Store) *Service {
	return &Service{store: st}
}

func (s *Service) Record(action string, target string, detail string, at int64) error {
	if strings.TrimSpace(action) == "" {
		return errors.New("audit action is required")
	}
	records, err := s.loadAll()
	if err != nil {
		return err
	}
	record := domain.AuditRecord{
		ID:     uuid.NewString(),
		Action: action,
		Target: target,
		Detail: detail,
		At:     at,
	}
	records = append(records, record)
	return s.store.Save(store.AuditKey, records)
}

func (s *Service) Recent(limit int) ([]domain.AuditRecord, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit >= len(records) {
		return records, nil
	}
	return records[len(records)-limit:], nil
}

func (s *Service) Count(action string) (int, error) {
	records, err := s.loadAll()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, record := range records {
		if record.Action == action {
			count++
		}
	}
	return count, nil
}

func (s *Service) loadAll() ([]domain.AuditRecord, error) {
	var records []domain.AuditRecord
	if err := s.store.Load(store.AuditKey, &records); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return []domain.AuditRecord{}, nil
		}
		return nil, err
	}
	return records, nil
}

func (s *Service) All() ([]domain.AuditRecord, error) {
	return s.loadAll()
}

func (s *Service) Clear() error {
	return s.store.Delete(store.AuditKey)
}

func (s *Service) Size() (int, error) {
	records, err := s.loadAll()
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (s *Service) Actions() ([]string, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	actions := make([]string, 0)
	for _, record := range records {
		if !seen[record.Action] {
			seen[record.Action] = true
			actions = append(actions, record.Action)
		}
	}
	return actions, nil
}
