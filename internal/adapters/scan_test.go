package adapters

import (
	"context"
	"github.com/iamunnip/infra-drift-watch/internal/config"
	"testing"
)

func TestPlanInterpretation(t *testing.T) {
	for _, tt := range []struct {
		name, data, status string
		count              int
	}{
		{"clean", `{"format_version":"1.2","resource_drift":[]}`, "clean", 0},
		{"pending only", `{"format_version":"1.2","resource_changes":[{"address":"x","change":{"actions":["create"]}}]}`, "clean", 0},
		{"drift", `{"format_version":"1.2","resource_drift":[{"address":"aws_s3_bucket.demo","change":{"actions":["update"],"before":{"secret":"never-print"},"after":{"secret":"changed"}}}]}`, "drift", 1},
		{"unsupported", `{"format_version":"99.0"}`, "failed", 0},
		{"invalid", `{`, "failed", 0},
		{"errored", `{"format_version":"1.2","errored":true}`, "failed", 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := ParsePlan("aws", []byte(tt.data))
			if r.Status != tt.status || len(r.Findings) != tt.count {
				t.Fatalf("%+v", r)
			}
		})
	}
}

func TestKubernetesDiffExitCodes(t *testing.T) {
	for _, code := range []int{0, 1, 2} {
		runner := func(ctx context.Context, dir, exe string, args []string, env []string) ([]byte, int, error) {
			return []byte("secret output"), code, nil
		}
		s := Scanner{Run: runner}
		r := s.Scan(context.Background(), config.Target{ID: "cluster", Adapter: "kubernetes", Context: "demo", Manifests: []string{"demo.yaml"}, Settings: config.Settings{Timeout: "1m"}})
		want := map[int]string{0: "clean", 1: "drift", 2: "failed"}[code]
		if r.Status != want {
			t.Fatalf("code %d: %+v", code, r)
		}
	}
}
