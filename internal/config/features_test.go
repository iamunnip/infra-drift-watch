package config

import (
	"strings"
	"testing"
)

func TestHistoryAndClosureConfiguration(t *testing.T) {
	c, err := Parse(strings.NewReader(`defaults:
  history:
    enabled: true
    limit: 12
  github:
    auto_close: true
`), strings.NewReader(`targets:
  - id: aws
    adapter: opentofu
    path: .
    settings:
      history:
        enabled: false
      github:
        auto_close: false
`))
	if err != nil {
		t.Fatal(err)
	}
	if Enabled(c.Targets[0].Settings.History.Enabled) || *c.Targets[0].Settings.History.Limit != 12 || Enabled(c.Targets[0].Settings.GitHub.AutoClose) {
		t.Fatal("feature overrides lost")
	}
	for _, global := range []string{
		"defaults:\n  history:\n    limit: 0",
		"defaults:\n  jira:\n    auto_close: true",
	} {
		if _, err := Parse(strings.NewReader(global), strings.NewReader("targets:\n - {id: a, adapter: terraform, path: .}")); err == nil {
			t.Fatal("invalid retention or Jira closure accepted")
		}
	}
}
