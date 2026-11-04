package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

// Record stores only public result metadata and an opaque change fingerprint.
type Record struct {
	At             time.Time    `json:"at"`
	DurationMS     int64        `json:"duration_ms"`
	Result         drift.Result `json:"result"`
	DeliveryErrors []string     `json:"delivery_errors"`
	Fingerprint    string       `json:"fingerprint"`
}

// Record atomically saves a target's latest bounded scan history alongside receipts.
func (s *Store) Record(record Record, limit int) error {
	if limit < 1 || limit > 10000 || record.Result.Target == "" {
		return errors.New("invalid scan history record or retention limit")
	}
	target := record.Result.Target
	record.At = record.At.UTC()
	record.Fingerprint = record.Result.Fingerprint()
	record.Result.Digest = ""
	record.Result.Normalize()
	if record.DeliveryErrors == nil {
		record.DeliveryErrors = []string{}
	}
	previous := s.History[target]
	updated := append(append([]Record{}, previous...), record)
	if len(updated) > limit {
		updated = updated[len(updated)-limit:]
	}
	s.History[target] = updated
	if err := s.Save(); err != nil {
		if previous == nil {
			delete(s.History, target)
		} else {
			s.History[target] = previous
		}
		return errors.New("cannot persist scan history")
	}
	return nil
}

// ReadHistory reads an atomic snapshot without creating directories or locking a writer.
func ReadHistory(dir, target string, limit int) ([]Record, error) {
	if limit < 1 || limit > 10000 {
		return nil, errors.New("history limit must be between 1 and 10000")
	}
	data, err := os.ReadFile(filepath.Join(dir, "receipts.json"))
	if err != nil {
		return nil, errors.New("cannot read history state file")
	}
	var store Store
	if json.Unmarshal(data, &store) != nil || (store.Version != 1 && store.Version != 2) || store.Entries == nil {
		return nil, errors.New("invalid or unsupported history state file")
	}
	records := []Record{}
	for id, list := range store.History {
		if target == "" || target == id {
			records = append(records, list...)
		}
	}
	sort.SliceStable(records, func(i, j int) bool {
		if records[i].At.Equal(records[j].At) {
			return records[i].Result.Target < records[j].Result.Target
		}
		return records[i].At.After(records[j].At)
	})
	if len(records) > limit {
		records = records[:limit]
	}
	return records, nil
}
