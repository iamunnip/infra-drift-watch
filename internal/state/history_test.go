package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

func TestHistoryRetentionAndVersionOneMigration(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "receipts.json"), []byte(`{"version":1,"entries":{"old":{"fingerprint":"old","active":true}}}`), 0600)
	s, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	for i, status := range []string{"drift", "failed", "clean"} {
		if err := s.Record(Record{At: time.Unix(int64(i), 0), Result: drift.Result{Target: "aws", Status: status}}, 2); err != nil {
			t.Fatal(err)
		}
	}
	s.Close()
	records, err := ReadHistory(dir, "aws", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || records[0].Result.Status != "clean" || records[1].Result.Status != "failed" {
		t.Fatalf("%+v", records)
	}
	s, err = Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	if !s.Entries["old"].Active || s.Version != 2 {
		t.Fatal("migration lost delivery receipts")
	}
}
func TestReadingAbsentHistoryDoesNotCreateState(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "absent")
	if _, err := ReadHistory(dir, "", 10); err == nil {
		t.Fatal("missing state accepted")
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("read created state")
	}
}

func TestFailedHistorySavePreservesPreviousRecord(t *testing.T) {
	store, err := Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	first := Record{At: time.Now(), Result: drift.Result{Target: "aws", Status: "drift"}}
	if err := store.Record(first, 2); err != nil {
		t.Fatal(err)
	}
	store.dir = filepath.Join(store.dir, "missing-directory")
	if err := store.Record(Record{At: time.Now(), Result: drift.Result{Target: "aws", Status: "clean"}}, 2); err == nil {
		t.Fatal("failed persistence hidden")
	}
	if len(store.History["aws"]) != 1 || store.History["aws"][0].Result.Status != "drift" {
		t.Fatal("failed save replaced valid history")
	}
}
