package zone

import (
	"errors"
	"strings"

	"venueops/internal/store"
)

func (s *Service) PartitionOf(zoneID string) (string, error) {
	zone, err := s.GetZone(zoneID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(zone.Partition) == "" {
		return zone.ID, nil
	}
	return zone.Partition, nil
}

func (s *Service) Split(hallID string, partitionA string, partitionB string) error {
	if strings.TrimSpace(partitionA) == "" || strings.TrimSpace(partitionB) == "" {
		return errors.New("both partition names are required")
	}
	zones, err := s.ZonesOfHall(hallID)
	if err != nil {
		return err
	}
	mid := len(zones) / 2
	if mid == 0 {
		mid = 1
	}
	for index, zone := range zones {
		if index < mid {
			zone.Partition = partitionA
		} else {
			zone.Partition = partitionB
		}
		if err := s.store.Save(store.ZoneStoreKey(zone.ID), zone); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) PartitionsOfHall(hallID string) ([]string, error) {
	zones, err := s.ZonesOfHall(hallID)
	if err != nil {
		return nil, err
	}
	seen := map[string]bool{}
	partitions := make([]string, 0)
	for _, zone := range zones {
		partition := zone.Partition
		if partition == "" {
			partition = zone.ID
		}
		if !seen[partition] {
			seen[partition] = true
			partitions = append(partitions, partition)
		}
	}
	return partitions, nil
}

func (s *Service) ZonesInPartition(hallID string, partition string) ([]string, error) {
	zones, err := s.ZonesOfHall(hallID)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	for _, zone := range zones {
		current := zone.Partition
		if current == "" {
			current = zone.ID
		}
		if current == partition {
			ids = append(ids, zone.ID)
		}
	}
	return ids, nil
}

func (s *Service) SetPartition(zoneID string, partition string) error {
	if strings.TrimSpace(partition) == "" {
		return errors.New("partition is required")
	}
	zone, err := s.GetZone(zoneID)
	if err != nil {
		return err
	}
	zone.Partition = partition
	return s.store.Save(store.ZoneStoreKey(zoneID), zone)
}

func (s *Service) PartitionCount(hallID string) (int, error) {
	partitions, err := s.PartitionsOfHall(hallID)
	if err != nil {
		return 0, err
	}
	return len(partitions), nil
}
