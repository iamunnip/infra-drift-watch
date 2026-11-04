package config

import (
	"strings"
	"testing"
)

func TestRejectAdditionalDocumentAndDuplicateKeys(t *testing.T) {
	for _, global := range []string{"{}\n---\n{}", "defaults:\n  timeout: 1m\n  timeout: 2m", "defaults:\n  github:\n    enabled: true\n    repository: a/b\n    api_url: http://unsafe"} {
		if _, err := Parse(strings.NewReader(global), strings.NewReader("targets:\n - {id: a, adapter: terraform, path: .}")); err == nil {
			t.Fatalf("accepted %q", global)
		}
	}
}
func TestExplicitEmptyListClearsDefault(t *testing.T) {
	c, err := Parse(strings.NewReader("defaults:\n  github:\n    labels: [drift]"), strings.NewReader("targets:\n - id: a\n   adapter: terraform\n   path: .\n   settings:\n     github:\n       labels: []"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Targets[0].Settings.GitHub.Labels == nil || len(c.Targets[0].Settings.GitHub.Labels) != 0 {
		t.Fatal("empty list did not clear default")
	}
}
