//go:build integration

package integration

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/adapters"
	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestRealObservedDriftDoesNotRefreshState(t *testing.T) {
	for _, adapter := range []string{"opentofu", "terraform"} {
		t.Run(adapter, func(t *testing.T) {
			name := adapter
			if name == "opentofu" {
				name = "tofu"
			}
			exe, err := exec.LookPath(name)
			if err != nil {
				t.Fatal(err)
			}
			dir := t.TempDir()
			managed := filepath.Join(dir, "managed.txt")
			source := fmt.Sprintf(`terraform {
 required_providers {
  local = { source = "hashicorp/local", version = "2.5.3" }
 }
}
resource "local_file" "demo" {
 filename = %q
 content = "disposable-test-content"
}
`, managed)
			os.WriteFile(filepath.Join(dir, "main.tf"), []byte(source), 0600)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			for _, args := range [][]string{{"init", "-input=false"}, {"apply", "-auto-approve", "-input=false"}} {
				cmd := exec.CommandContext(ctx, exe, args...)
				cmd.Dir = dir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("local fixture setup: %s %v", out, err)
				}
			}
			before, err := os.ReadFile(filepath.Join(dir, "terraform.tfstate"))
			if err != nil {
				t.Fatal(err)
			}
			target := config.Target{ID: "local-file", Adapter: adapter, Path: dir, Settings: config.Settings{Timeout: "1m", Executable: exe}}
			if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "clean" {
				t.Fatalf("expected clean: %+v", r)
			}
			os.Remove(managed)
			r := (adapters.Scanner{}).Scan(ctx, target)
			if r.Status != "drift" || len(r.Findings) != 1 || r.Findings[0].Resource != "local_file.demo" {
				t.Fatalf("expected observed deletion: %+v", r)
			}
			after, err := os.ReadFile(filepath.Join(dir, "terraform.tfstate"))
			if err != nil || !bytes.Equal(before, after) {
				t.Fatal("scan persisted refreshed state")
			}
			os.WriteFile(managed, []byte("disposable-test-content"), 0755)
			if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "clean" {
				t.Fatalf("expected resolution: %+v", r)
			}
		})
	}
}
