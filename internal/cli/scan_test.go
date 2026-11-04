package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateInventoryWithoutCredentials(t *testing.T) {
	dir := t.TempDir()
	global := filepath.Join(dir, "config.yaml")
	inv := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(global, []byte("defaults:\n  github:\n    enabled: true\n    repository: owner/issues\n"), 0600)
	os.WriteFile(inv, []byte("targets:\n  - id: aws\n    adapter: opentofu\n    path: .\n"), 0600)
	var out, errOut bytes.Buffer
	code := Run([]string{"validate", "--config", global, "--inventory", inv}, &out, &errOut)
	if code != 0 || !strings.Contains(out.String(), "1 target") {
		t.Fatalf("code=%d out=%s err=%s", code, &out, &errOut)
	}
	if _, err := os.Stat(filepath.Join(dir, ".state")); !os.IsNotExist(err) {
		t.Fatal("validation wrote delivery state")
	}
}
func TestScanInvalidTargetSelection(t *testing.T) {
	dir := t.TempDir()
	g := filepath.Join(dir, "config.yaml")
	i := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(g, []byte("{}"), 0600)
	os.WriteFile(i, []byte("targets:\n - {id: aws, adapter: opentofu, path: .}"), 0600)
	var out, err bytes.Buffer
	code := Run([]string{"scan", "--config", g, "--inventory", i, "--target", "missing", "--dry-run"}, &out, &err)
	if code != 1 || !strings.Contains(err.String(), "target selection") {
		t.Fatalf("code=%d err=%s", code, &err)
	}
}

func TestValidateAWSCLISelectorsWithoutCredentials(t *testing.T) {
	dir := t.TempDir()
	global := filepath.Join(dir, "config.yaml")
	inventory := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(global, []byte("{}"), 0600)
	os.WriteFile(inventory, []byte("targets:\n - {id: aws, adapter: opentofu, path: .}"), 0600)
	for _, tt := range []struct {
		profile, region string
		code            int
	}{{"default", "us-east-1", 0}, {"", "ap-south-1", 0}, {"default", "invalid", 1}} {
		var out, errOut bytes.Buffer
		got := Run([]string{"validate", "--config", global, "--inventory", inventory, "--aws-profile", tt.profile, "--aws-region", tt.region}, &out, &errOut)
		if got != tt.code {
			t.Fatalf("code=%d want=%d stderr=%s", got, tt.code, &errOut)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".state")); !os.IsNotExist(err) {
		t.Fatal("validation wrote state")
	}
}
