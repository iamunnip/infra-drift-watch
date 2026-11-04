package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDryRunScansWithoutCredentialsDeliveryOrState(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "fake-tofu")
	script := `#!/bin/sh
case "$1" in
 plan) exit 2 ;;
 show) echo '{"format_version":"1.2","resource_drift":[{"address":"aws_s3_bucket.demo","change":{"actions":["update"],"before":{"secret":"must-not-leak"},"after":{"secret":"also-private"}}}]}' ;;
 *) exit 1 ;;
esac
`
	if err := os.WriteFile(exe, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	global := filepath.Join(dir, "config.yaml")
	inventory := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(global, []byte("defaults:\n  github:\n    enabled: true\n    repository: a/b\n    token_env: MISSING_DRIFT_TEST_TOKEN"), 0600)
	os.WriteFile(inventory, []byte("targets:\n - id: aws\n   adapter: opentofu\n   path: .\n   settings:\n     executable: "+exe), 0600)
	var out bytes.Buffer
	code, err := Scan(context.Background(), Options{Config: global, Inventory: inventory, DryRun: true, Format: "json"}, &out)
	if err != nil || code != 2 {
		t.Fatalf("code=%d err=%v out=%s", code, err, &out)
	}
	if !json.Valid(out.Bytes()) || strings.Contains(out.String(), "must-not-leak") || strings.Contains(out.String(), "also-private") {
		t.Fatal("invalid or leaking JSON report")
	}
	if _, err := os.Stat(filepath.Join(dir, ".state")); !os.IsNotExist(err) {
		t.Fatal("dry run wrote state")
	}
}

func TestMissingCredentialFailsBeforeExecution(t *testing.T) {
	dir := t.TempDir()
	global := filepath.Join(dir, "config.yaml")
	inventory := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(global, []byte("defaults:\n  github:\n    enabled: true\n    repository: a/b\n    token_env: DRIFT_ABSENT_TOKEN"), 0600)
	os.WriteFile(inventory, []byte("targets:\n - {id: aws, adapter: terraform, path: .}"), 0600)
	t.Setenv("DRIFT_ABSENT_TOKEN", "")
	code, err := Scan(context.Background(), Options{Config: global, Inventory: inventory}, &bytes.Buffer{})
	if code != 1 || err == nil || !strings.Contains(err.Error(), "credential") {
		t.Fatalf("%d %v", code, err)
	}
}

func TestFailureTakesPrecedenceAndLaterTargetStillScans(t *testing.T) {
	dir := t.TempDir()
	exe := filepath.Join(dir, "tool")
	os.WriteFile(exe, []byte("#!/bin/sh\nif [ \"$1\" = plan ]; then exit 2; fi\necho '{\"format_version\":\"1.2\",\"resource_drift\":[{\"address\":\"x\",\"change\":{\"actions\":[\"delete\"]}}]}'\n"), 0700)
	g := filepath.Join(dir, "config.yaml")
	i := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(g, []byte("{}"), 0600)
	os.WriteFile(i, []byte("targets:\n - id: a-failure\n   adapter: terraform\n   path: .\n   settings:\n     executable: /does-not-exist/terraform\n - id: z-drift\n   adapter: opentofu\n   path: .\n   settings:\n     executable: "+exe), 0600)
	var out bytes.Buffer
	code, err := Scan(context.Background(), Options{Config: g, Inventory: i, DryRun: true, Format: "json"}, &out)
	if code != 1 || err != nil || !strings.Contains(out.String(), `"status": "failed"`) || !strings.Contains(out.String(), `"status": "drift"`) {
		t.Fatalf("code=%d err=%v out=%s", code, err, &out)
	}
}
