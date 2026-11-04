//go:build integration

package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/adapters"
	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestRealStateTools(t *testing.T) {
	for _, adapter := range []string{"opentofu", "terraform"} {
		t.Run(adapter, func(t *testing.T) {
			name := adapter
			if name == "opentofu" {
				name = "tofu"
			}
			exe, err := exec.LookPath(name)
			if err != nil {
				t.Fatalf("install %s before integration checks", name)
			}
			dir := t.TempDir()
			// Built-in terraform_data needs no cloud access or external provider downloads.
			source := `resource "terraform_data" "demo" { input = "first" }`
			if err := os.WriteFile(filepath.Join(dir, "main.tf"), []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			for _, args := range [][]string{{"init", "-input=false"}, {"apply", "-auto-approve", "-input=false"}} {
				cmd := exec.CommandContext(ctx, exe, args...)
				cmd.Dir = dir
				if out, err := cmd.CombinedOutput(); err != nil {
					t.Fatalf("fixture setup: %s: %v", out, err)
				}
			}
			target := config.Target{ID: "local", Adapter: adapter, Path: dir, Settings: config.Settings{Timeout: "1m", Executable: exe}}
			for _, input := range []string{"first", "pending-configuration"} {
				os.WriteFile(filepath.Join(dir, "main.tf"), []byte(`resource "terraform_data" "demo" { input = "`+input+`" }`), 0600)
				result := (adapters.Scanner{}).Scan(ctx, target)
				if result.Status != "clean" {
					t.Fatalf("pending configuration incorrectly treated as observed drift: %+v", result)
				}
			}
		})
	}
}
