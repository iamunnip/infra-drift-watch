package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAWSCLIOverridesTargetSettings(t *testing.T) {
	dir := t.TempDir()
	global := filepath.Join(dir, "config.yaml")
	inventory := filepath.Join(dir, "inventory.yaml")
	os.WriteFile(global, []byte("defaults:\n  aws:\n    profile: default\n    region: us-east-1"), 0600)
	os.WriteFile(inventory, []byte("targets:\n - id: aws\n   adapter: opentofu\n   path: .\n   settings:\n     aws:\n       profile: target\n       region: eu-west-1"), 0600)
	profile, region := "", "ap-south-1"
	c, err := Load(Options{Config: global, Inventory: inventory, AWSProfile: &profile, AWSRegion: &region})
	if err != nil {
		t.Fatal(err)
	}
	if *c.Targets[0].Settings.AWS.Profile != "" || *c.Targets[0].Settings.AWS.Region != "ap-south-1" {
		t.Fatal("CLI selection did not override target")
	}
}
