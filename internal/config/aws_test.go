package config

import (
	"strings"
	"testing"
)

func TestAWSSettingsInheritOverrideAndClear(t *testing.T) {
	c, err := Parse(strings.NewReader(`defaults:
  aws:
    profile: default
    region: us-east-1
`), strings.NewReader(`targets:
 - {id: inherited, adapter: opentofu, path: .}
 - id: override
   adapter: terraform
   path: .
   settings:
     aws:
       profile: production
       region: eu-west-1
 - id: oidc
   adapter: opentofu
   path: .
   settings:
     aws:
       profile: ""
`))
	if err != nil {
		t.Fatal(err)
	}
	if *c.Targets[0].Settings.AWS.Profile != "default" || *c.Targets[0].Settings.AWS.Region != "us-east-1" {
		t.Fatal("defaults lost")
	}
	if *c.Targets[1].Settings.AWS.Profile != "production" || *c.Targets[1].Settings.AWS.Region != "eu-west-1" {
		t.Fatal("target override lost")
	}
	if c.Targets[2].Settings.AWS.Profile == nil || *c.Targets[2].Settings.AWS.Profile != "" || *c.Targets[2].Settings.AWS.Region != "us-east-1" {
		t.Fatal("explicit clear lost")
	}
}

func TestAWSSelectionValidation(t *testing.T) {
	for _, value := range []string{" bad", "bad\nprofile", "bad\x00profile", "bad\x1bprofile"} {
		if err := ValidateAWS(AWS{Profile: &value}); err == nil {
			t.Fatalf("invalid profile accepted: %q", value)
		}
	}
	for _, value := range []string{"us east 1", "us-east", "https://example.com"} {
		if err := ValidateAWS(AWS{Region: &value}); err == nil {
			t.Fatalf("invalid region accepted: %q", value)
		}
	}
	for _, value := range []string{"us-east-1", "ap-south-1", "us-gov-west-1", "eusc-de-east-1", ""} {
		if err := ValidateAWS(AWS{Region: &value}); err != nil {
			t.Fatalf("valid region rejected: %q", value)
		}
	}
}
