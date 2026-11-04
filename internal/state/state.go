// Package state persists delivery receipts under a single Linux writer lock.
package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

type Entry struct {
	Policy      string    `json:"policy,omitempty"`
	Link        string    `json:"link,omitempty"`
	Fingerprint string    `json:"fingerprint"`
	Active      bool      `json:"active"`
	Updated     time.Time `json:"updated"`
}
type Store struct {
	History map[string][]Record `json:"history,omitempty"`
	Entries map[string]Entry    `json:"entries"`
	Version int                 `json:"version"`
	dir     string
	lock    *os.File
}

func Open(dir string) (*Store, error) {
	if e := os.MkdirAll(dir, 0700); e != nil {
		return nil, errors.New("cannot create state directory")
	}
	f, e := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return nil, errors.New("cannot open state lock")
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return nil, errors.New("state directory already has a writer")
	}
	s := &Store{Entries: map[string]Entry{}, Version: 2, History: map[string][]Record{}, dir: dir, lock: f}
	b, e := os.ReadFile(filepath.Join(dir, "receipts.json"))
	if e == nil {
		if json.Unmarshal(b, s) != nil || (s.Version != 1 && s.Version != 2) || s.Entries == nil {
			s.Close()
			return nil, errors.New("invalid or unsupported state file")
		}
	} else if !os.IsNotExist(e) {
		s.Close()
		return nil, errors.New("cannot read state file")
	}
	s.Version = 2
	if s.History == nil {
		s.History = map[string][]Record{}
	}
	return s, nil
}
func (s *Store) Save() error {
	b, e := json.MarshalIndent(s, "", "  ")
	if e != nil {
		return e
	}
	f, e := os.CreateTemp(s.dir, "receipts-*")
	if e != nil {
		return errors.New("cannot create state temporary file")
	}
	defer os.Remove(f.Name())
	if _, e = f.Write(b); e == nil {
		e = f.Sync()
	}
	if closeErr := f.Close(); e == nil {
		e = closeErr
	}
	if e == nil {
		e = os.Rename(f.Name(), filepath.Join(s.dir, "receipts.json"))
	}
	if e != nil {
		return errors.New("cannot persist delivery state")
	}
	// Persist directory metadata after the atomic rename.
	d, e := os.Open(s.dir)
	if e != nil {
		return e
	}
	defer d.Close()
	return d.Sync()
}
func (s *Store) Close() {
	if s.lock != nil {
		syscall.Flock(int(s.lock.Fd()), syscall.LOCK_UN)
		s.lock.Close()
		s.lock = nil
	}
}
