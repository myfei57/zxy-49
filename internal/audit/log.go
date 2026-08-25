package audit

import (
	"errors"
	"time"

	"venueops/internal/domain"
)

func (s *Service) RecordNow(action string, target string, detail string) error {
	return s.Record(action, target, detail, time.Now().Unix())
}

func (s *Service) Last() (domain.AuditRecord, error) {
	records, err := s.loadAll()
	if err != nil {
		return domain.AuditRecord{}, err
	}
	if len(records) == 0 {
		return domain.AuditRecord{}, errors.New("no audit records")
	}
	return records[len(records)-1], nil
}

func (s *Service) ByAction(action string) ([]domain.AuditRecord, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.AuditRecord, 0)
	for _, record := range records {
		if record.Action == action {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Service) ByTarget(target string) ([]domain.AuditRecord, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.AuditRecord, 0)
	for _, record := range records {
		if record.Target == target {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Service) Since(at int64) ([]domain.AuditRecord, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.AuditRecord, 0)
	for _, record := range records {
		if record.At >= at {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}

func (s *Service) DistinctTargets() ([]string, error) {
	records, err := s.loadAll()
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	targets := make([]string, 0)
	for _, record := range records {
		if !seen[record.Target] {
			seen[record.Target] = true
			targets = append(targets, record.Target)
		}
	}
	return targets, nil
}
