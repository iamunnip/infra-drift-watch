package config

import (
	"strings"
	"testing"
)

func TestStrictConfigurationAndOverrides(t *testing.T) {
	global := `defaults:
  timeout: 2m
  github:
    enabled: true
    repository: owner/issues
    token_env: GH_TOKEN
`
	inv := `targets:
  - id: aws-prod
    adapter: opentofu
    path: .
    settings:
      timeout: 5m
      github:
        enabled: false
`
	c, err := Parse(strings.NewReader(global), strings.NewReader(inv))
	if err != nil {
		t.Fatal(err)
	}
	if c.Targets[0].Settings.Timeout != "5m" || *c.Targets[0].Settings.GitHub.Enabled {
		t.Fatalf("override not applied: %+v", c.Targets[0])
	}
	if c.Targets[0].Settings.GitHub.Repository != "owner/issues" {
		t.Fatal("partial object lost inherited setting")
	}
}

func TestRejectInvalidConfiguration(t *testing.T) {
	for _, inv := range []string{
		"targets: []",
		"unknown: true",
		"targets:\n - {id: a, adapter: invalid, path: .}",
		"targets:\n - {id: a, adapter: opentofu, path: ., settings: {timeout: bad}}",
		"targets:\n - {id: a, adapter: opentofu, path: .}\n - {id: a, adapter: terraform, path: .}",
	} {
		if _, err := Parse(strings.NewReader("{}"), strings.NewReader(inv)); err == nil {
			t.Fatalf("accepted %q", inv)
		}
	}
}
