package adapters

import (
	"context"
	"github.com/iamunnip/infra-drift-watch/internal/config"
	"strings"
	"testing"
)

func TestChangedDriftHasNewFingerprintWithoutValues(t *testing.T) {
	a := ParsePlan("a", []byte(`{"format_version":"1.2","resource_drift":[{"address":"x","change":{"actions":["update"],"after":{"value":"one"}}}]}`))
	b := ParsePlan("a", []byte(`{"format_version":"1.2","resource_drift":[{"address":"x","change":{"actions":["update"],"after":{"value":"two"}}}]}`))
	if a.Fingerprint() == b.Fingerprint() {
		t.Fatal("changed remote values must produce a new delivery event")
	}
}

func TestIncompletePlansCannotReportClean(t *testing.T) {
	for _, data := range []string{
		`{"format_version":"1.2","complete":false}`,
		`{"format_version":"1.2","deferred_changes":[{"reason":"unknown"}]}`,
	} {
		if r := ParsePlan("a", []byte(data)); r.Status != "incomplete" {
			t.Fatalf("incomplete plan reported %+v", r)
		}
	}
}

func TestKubernetesFingerprintIgnoresTemporaryDirectoryNames(t *testing.T) {
	output := "diff /tmp/LIVE-123/v1.ConfigMap.default.demo /tmp/MERGED-456/v1.ConfigMap.default.demo\n1c1\n< before\n---\n> after\n"
	runner := func(context.Context, string, string, []string, []string) ([]byte, int, error) {
		return []byte(output), 1, nil
	}
	s := Scanner{Run: runner}
	target := config.Target{ID: "k", Adapter: "kubernetes", Context: "demo", Manifests: []string{"demo.yaml"}, Settings: config.Settings{Timeout: "1m"}}
	first := s.Scan(context.Background(), target)
	output = strings.ReplaceAll(strings.ReplaceAll(output, "LIVE-123", "LIVE-999"), "MERGED-456", "MERGED-888")
	second := s.Scan(context.Background(), target)
	if first.Fingerprint() != second.Fingerprint() {
		t.Fatal("temporary filenames created a duplicate drift event")
	}
}
