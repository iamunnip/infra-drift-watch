package adapters

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestRefreshOnlyAndPrivateArtifactCleanup(t *testing.T) {
	var artifact string
	calls := 0
	runner := func(ctx context.Context, dir, exe string, args, env []string) ([]byte, int, error) {
		calls++
		if calls == 1 {
			joined := strings.Join(args, " ")
			if !strings.Contains(joined, "-refresh-only") || strings.Contains(joined, "apply") {
				t.Fatal(joined)
			}
			for _, a := range args {
				if strings.HasPrefix(a, "-out=") {
					artifact = strings.TrimPrefix(a, "-out=")
				}
			}
			info, err := os.Stat(filepath.Dir(artifact))
			if err != nil || info.Mode().Perm() != 0700 {
				t.Fatal("plan directory not private")
			}
			return nil, 2, nil
		}
		return []byte(`{"format_version":"1.2","resource_drift":[{"address":"x","change":{"actions":["delete"]}}]}`), 0, nil
	}
	result := (Scanner{Run: runner}).Scan(context.Background(), config.Target{ID: "a", Adapter: "terraform", Path: ".", Settings: config.Settings{Timeout: "1m"}})
	if result.Status != "drift" || calls != 2 {
		t.Fatalf("%+v calls=%d", result, calls)
	}
	if _, err := os.Stat(filepath.Dir(artifact)); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("private artifact directory left behind")
	}
}

func TestExpiredSuppressionDoesNotHideDrift(t *testing.T) {
	r := ParsePlan("a", []byte(`{"format_version":"1.2","resource_drift":[{"address":"x","change":{"actions":["update"]}}]}`))
	target := config.Target{Suppressions: []config.Suppression{{Resource: "x", Reason: "review", Until: "2000-01-01"}}}
	if got := suppress(target, r); got.Status != "drift" {
		t.Fatal(got)
	}
	target.Suppressions[0].Until = "2999-01-01"
	if got := suppress(target, r); got.Status != "incomplete" || got.Suppressed != 1 {
		t.Fatal(got)
	}
}

func TestExecuteCancellationAndSecretStderr(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := Execute(ctx, "", "sh", []string{"-c", "sleep 10"}, nil); err == nil {
		t.Fatal("cancel ignored")
	}
	out, code, err := Execute(context.Background(), "", "sh", []string{"-c", "echo secret >&2; exit 1"}, nil)
	if err != nil || code != 1 || len(out) != 0 {
		t.Fatalf("out=%s code=%d err=%v", out, code, err)
	}
}

func TestExecuteRejectsOversizedOutputEvenWithDriftExitCode(t *testing.T) {
	out, code, err := Execute(context.Background(), "", "sh", []string{"-c", "head -c 33554433 /dev/zero; exit 1"}, nil)
	if err == nil || code >= 0 || out != nil {
		t.Fatalf("truncated output accepted: bytes=%d code=%d err=%v", len(out), code, err)
	}
}
