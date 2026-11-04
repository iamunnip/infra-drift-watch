//go:build kubernetes

package integration

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/adapters"
	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestDisposableKubernetes(t *testing.T) {
	kind, err := exec.LookPath("kind")
	if err != nil {
		t.Fatal("kind is required")
	}
	kubectl, err := exec.LookPath("kubectl")
	if err != nil {
		t.Fatal("kubectl is required")
	}
	dir := t.TempDir()
	kubeconfig := filepath.Join(dir, "kubeconfig")
	t.Setenv("KUBECONFIG", kubeconfig)
	name := fmt.Sprintf("drift-watch-%d", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	t.Cleanup(func() {
		cleanup, stop := context.WithTimeout(context.Background(), time.Minute)
		defer stop()
		cmd := exec.CommandContext(cleanup, kind, "delete", "cluster", "--name", name, "--kubeconfig", kubeconfig)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("cleanup failed: %s %v", out, err)
		}
	})
	cmd := exec.CommandContext(ctx, kind, "create", "cluster", "--name", name, "--kubeconfig", kubeconfig, "--image", "kindest/node:v1.36.4", "--wait", "90s")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create disposable cluster: %s %v", out, err)
	}
	manifest := filepath.Join(dir, "configmap.yaml")
	os.WriteFile(manifest, []byte("apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: drift-demo\n  namespace: default\ndata:\n  greeting: hello\n"), 0600)
	invoke := func(args ...string) {
		t.Helper()
		cmd := exec.CommandContext(ctx, kubectl, args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("fixture: %s %v", out, err)
		}
	}
	invoke("apply", "--server-side", "--field-manager=kubectl-client-side-apply", "-f", manifest)
	target := config.Target{ID: "cluster", Adapter: "kubernetes", Context: "kind-" + name, Namespace: "default", Manifests: []string{manifest}, Settings: config.Settings{Timeout: "30s", Executable: kubectl}}
	if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "clean" {
		t.Fatalf("expected clean: %+v", r)
	}
	// A separate deployment using the same apply manager changes desired ownership values.
	changed := filepath.Join(dir, "changed.yaml")
	source, _ := os.ReadFile(manifest)
	os.WriteFile(changed, []byte(strings.ReplaceAll(string(source), "greeting: hello", "greeting: changed")), 0600)
	invoke("apply", "--server-side", "--field-manager=kubectl-client-side-apply", "-f", changed)
	if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "drift" {
		out, _ := exec.CommandContext(ctx, kubectl, "--context", target.Context, "diff", "--server-side=true", "--field-manager=kubectl-client-side-apply", "-f", manifest).CombinedOutput()
		t.Fatalf("expected drift: %+v; kubectl diagnostic: %s", r, out)
	}
	first := (adapters.Scanner{}).Scan(ctx, target)
	second := (adapters.Scanner{}).Scan(ctx, target)
	if first.Fingerprint() != second.Fingerprint() {
		t.Fatalf("unchanged drift fingerprint changed: %s vs %s", first.Fingerprint(), second.Fingerprint())
	}
	invoke("apply", "--server-side", "--field-manager=kubectl-client-side-apply", "-f", manifest)
	if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "clean" {
		t.Fatalf("expected resolved: %+v", r)
	}
	invoke("patch", "configmap", "drift-demo", "--field-manager=external-editor", "--type=merge", "-p", `{"data":{"greeting":"conflicting"}}`)
	if r := (adapters.Scanner{}).Scan(ctx, target); r.Status != "failed" {
		t.Fatalf("ownership conflict must fail: %+v", r)
	}

}
