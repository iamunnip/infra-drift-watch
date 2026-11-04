package state

import "testing"

func TestPersistenceAndExclusiveWriter(t *testing.T) {
	dir := t.TempDir()
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Open(dir); err == nil {
		t.Fatal("allowed concurrent writer")
	}
	s.Entries["one"] = Entry{Fingerprint: "abc", Active: true}
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
	s.Close()
	s, err = Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if s.Entries["one"].Fingerprint != "abc" {
		t.Fatal("state lost")
	}
}
