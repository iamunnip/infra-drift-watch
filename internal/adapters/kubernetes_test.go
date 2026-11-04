package adapters

import (
	"context"
	"strings"
	"testing"

	"github.com/iamunnip/infra-drift-watch/internal/config"
)

func TestKubernetesFindingsIdentifyChangedResources(t *testing.T) {
	output := `diff /tmp/LIVE-123/v1.ConfigMap.default.app.settings /tmp/MERGED-456/v1.ConfigMap.default.app.settings
1c1
< private-old-value
---
> private-new-value
diff /tmp/LIVE-123/apps.v1.Deployment.production.web /tmp/MERGED-456/apps.v1.Deployment.production.web
1c1
< before
---
> after
`
	run := func(context.Context, string, string, []string, []string) ([]byte, int, error) {
		return []byte(output), 1, nil
	}
	r := (Scanner{Run: run}).Scan(context.Background(), config.Target{ID: "k", Adapter: "kubernetes", Context: "test", Manifests: []string{"demo.yaml"}, Settings: config.Settings{Timeout: "1m"}})
	if len(r.Findings) != 2 || r.Findings[0].Resource != "apps/v1/Deployment/production/web" || r.Findings[1].Resource != "v1/ConfigMap/default/app.settings" {
		t.Fatalf("expected resource identities, got %+v", r)
	}
	for _, f := range r.Findings {
		if strings.Contains(f.Resource, "private") {
			t.Fatal("diff values leaked")
		}
	}
}

func TestManifestSetSuppressionStillCoversDetailedFindings(t *testing.T) {
	output := "diff /tmp/LIVE-1/v1.ConfigMap.default.demo /tmp/MERGED-2/v1.ConfigMap.default.demo\n< old\n> new\n"
	run := func(context.Context, string, string, []string, []string) ([]byte, int, error) {
		return []byte(output), 1, nil
	}
	target := config.Target{ID: "k", Adapter: "kubernetes", Context: "test", Manifests: []string{"demo.yaml"}, Settings: config.Settings{Timeout: "1m"}, Suppressions: []config.Suppression{{Resource: "manifest-set", Reason: "review", Until: "2999-01-01"}}}
	result := (Scanner{Run: run}).Scan(context.Background(), target)
	if result.Status != "incomplete" || result.Suppressed != 1 {
		t.Fatalf("set suppression did not cover resource findings: %+v", result)
	}
}
