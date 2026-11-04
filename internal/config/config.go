// Package config loads strict YAML configuration and resolves target overrides.
package config

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"go.yaml.in/yaml/v3"
)

type GitHub struct {
	Enabled    *bool    `yaml:"enabled"`
	Repository string   `yaml:"repository"`
	TokenEnv   string   `yaml:"token_env"`
	APIURL     string   `yaml:"api_url"`
	Labels     []string `yaml:"labels"`
	Assignees  []string `yaml:"assignees"`
}
type Slack struct {
	Enabled    *bool  `yaml:"enabled"`
	WebhookEnv string `yaml:"webhook_env"`
	Repeat     string `yaml:"repeat"`
}
type Jira struct {
	Enabled   *bool  `yaml:"enabled"`
	URL       string `yaml:"url"`
	Project   string `yaml:"project"`
	IssueType string `yaml:"issue_type"`
	EmailEnv  string `yaml:"email_env"`
	TokenEnv  string `yaml:"token_env"`
	// Fields supports installation-specific custom fields on creation.
	Fields map[string]any `yaml:"fields"`
}
type Settings struct {
	Timeout     string `yaml:"timeout"`
	HTTPTimeout string `yaml:"http_timeout"`
	Retries     *int   `yaml:"retries"`
	Executable  string `yaml:"executable"`
	GitHub      GitHub `yaml:"github"`
	Slack       Slack  `yaml:"slack"`
	Jira        Jira   `yaml:"jira"`
}
type Suppression struct {
	Resource string `yaml:"resource"`
	Reason   string `yaml:"reason"`
	Until    string `yaml:"until"`
}
type Target struct {
	FieldManager string        `yaml:"field_manager"`
	ID           string        `yaml:"id"`
	Adapter      string        `yaml:"adapter"`
	Description  string        `yaml:"description"`
	Owner        string        `yaml:"owner"`
	Environment  string        `yaml:"environment"`
	Path         string        `yaml:"path"`
	Workspace    string        `yaml:"workspace"`
	VarFiles     []string      `yaml:"var_files"`
	Context      string        `yaml:"context"`
	Namespace    string        `yaml:"namespace"`
	Manifests    []string      `yaml:"manifests"`
	Suppressions []Suppression `yaml:"suppressions"`
	Settings     Settings      `yaml:"settings"`
}
type File struct {
	StateDir string   `yaml:"state_dir"`
	Format   string   `yaml:"format"`
	Defaults Settings `yaml:"defaults"`
}
type Inventory struct {
	Targets []Target `yaml:"targets"`
}
type Config struct {
	File
	Targets []Target
}

func Enabled(v *bool) bool { return v != nil && *v }
func decode(r io.Reader, v any) error {
	d := yaml.NewDecoder(io.LimitReader(r, 4<<20))
	d.KnownFields(true)
	if err := d.Decode(v); err != nil {
		return errors.New("invalid YAML or unknown configuration field")
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return errors.New("configuration must contain exactly one YAML document")
	}
	return nil
}

// merge applies nonzero scalars, explicit pointers, and replacement lists/maps.
// Empty lists/maps deliberately clear inherited values.
func merge(dst, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		a, b := dst.Field(i), src.Field(i)
		if b.Kind() == reflect.Struct {
			merge(a, b)
		} else if !b.IsZero() {
			a.Set(b)
		}
	}
}
func Parse(global, inventory io.Reader) (Config, error) {
	retries := 2
	c := Config{File: File{StateDir: ".state", Format: "terminal", Defaults: Settings{Timeout: "5m", HTTPTimeout: "30s", Retries: &retries, GitHub: GitHub{APIURL: "https://api.github.com", TokenEnv: "GITHUB_TOKEN"}, Slack: Slack{WebhookEnv: "SLACK_WEBHOOK_URL", Repeat: "24h"}, Jira: Jira{IssueType: "Task", EmailEnv: "JIRA_EMAIL", TokenEnv: "JIRA_TOKEN"}}}}
	var raw File
	if err := decode(global, &raw); err != nil {
		return c, err
	}
	merge(reflect.ValueOf(&c.File).Elem(), reflect.ValueOf(raw))
	var inv Inventory
	if err := decode(inventory, &inv); err != nil {
		return c, err
	}
	if len(inv.Targets) == 0 {
		return c, errors.New("inventory needs at least one target")
	}
	if c.Format != "terminal" && c.Format != "json" {
		return c, errors.New("format must be terminal or json")
	}
	seen := map[string]bool{}
	id := regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]{0,99}$`)
	for _, t := range inv.Targets {
		if !id.MatchString(t.ID) || seen[t.ID] {
			return c, errors.New("target IDs must be unique and contain only letters, numbers, dots, underscores or hyphens")
		}
		seen[t.ID] = true
		settings := c.Defaults
		merge(reflect.ValueOf(&settings).Elem(), reflect.ValueOf(t.Settings))
		t.Settings = settings
		if t.Adapter != "opentofu" && t.Adapter != "terraform" && t.Adapter != "kubernetes" {
			return c, fmt.Errorf("target %s: unsupported adapter", t.ID)
		}
		if t.Adapter == "kubernetes" {
			if t.Context == "" || len(t.Manifests) == 0 {
				return c, fmt.Errorf("target %s: context and manifests required", t.ID)
			}
			for _, p := range t.Manifests {
				if strings.HasPrefix(p, "-") || strings.Contains(p, "://") {
					return c, errors.New("manifest paths must be local paths")
				}
			}
		} else if t.Path == "" {
			return c, fmt.Errorf("target %s: path required", t.ID)
		}
		for _, d := range []string{settings.Timeout, settings.HTTPTimeout, settings.Slack.Repeat} {
			n, e := time.ParseDuration(d)
			if e != nil || n <= 0 {
				return c, fmt.Errorf("target %s: durations must be positive", t.ID)
			}
		}
		if *settings.Retries < 0 || *settings.Retries > 5 {
			return c, errors.New("retries must be between 0 and 5")
		}
		if Enabled(settings.GitHub.Enabled) {
			if !regexp.MustCompile(`^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$`).MatchString(settings.GitHub.Repository) {
				return c, errors.New("GitHub repository must be owner/name")
			}
			if err := secureURL(settings.GitHub.APIURL); err != nil {
				return c, err
			}
		}
		if Enabled(settings.Jira.Enabled) {
			if secureURL(settings.Jira.URL) != nil || !regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`).MatchString(settings.Jira.Project) || settings.Jira.IssueType == "" {
				return c, errors.New("Jira requires an HTTPS site, project key and issue type")
			}
			for k := range settings.Jira.Fields {
				if !strings.HasPrefix(k, "customfield_") {
					return c, errors.New("Jira fields only accepts customfield_ keys")
				}
			}
		}
		for _, s := range t.Suppressions {
			if s.Resource == "" || s.Reason == "" {
				return c, errors.New("suppression needs resource, reason and until")
			}
			if _, err := time.Parse("2006-01-02", s.Until); err != nil {
				return c, errors.New("suppression until must be YYYY-MM-DD")
			}
		}
		c.Targets = append(c.Targets, t)
	}
	return c, nil
}
func secureURL(s string) error {
	u, e := url.Parse(s)
	if e != nil || u.Scheme != "https" || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return errors.New("API URLs must be HTTPS without credentials, query or fragment")
	}
	return nil
}
func Load(global, inventory string) (Config, error) {
	g, e := os.Open(global)
	if e != nil {
		return Config{}, errors.New("cannot open configuration")
	}
	defer g.Close()
	i, e := os.Open(inventory)
	if e != nil {
		return Config{}, errors.New("cannot open inventory")
	}
	defer i.Close()
	c, e := Parse(g, i)
	if e != nil {
		return c, e
	}
	base, e := filepath.Abs(filepath.Dir(inventory))
	if e != nil {
		return c, e
	}
	for n := range c.Targets {
		t := &c.Targets[n]
		if t.Path != "" && !filepath.IsAbs(t.Path) {
			t.Path = filepath.Join(base, t.Path)
		}
		for j, p := range t.Manifests {
			if !filepath.IsAbs(p) {
				t.Manifests[j] = filepath.Join(base, p)
			}
		}
	}
	if !filepath.IsAbs(c.StateDir) {
		c.StateDir = filepath.Join(filepath.Dir(global), c.StateDir)
	}
	return c, nil
}

// ValidateSecrets checks all enabled destinations before scans or delivery begin.
func ValidateSecrets(targets []Target) error {
	for _, t := range targets {
		s := t.Settings
		for _, pair := range []struct {
			on   bool
			name string
		}{{Enabled(s.GitHub.Enabled), s.GitHub.TokenEnv}, {Enabled(s.Slack.Enabled), s.Slack.WebhookEnv}, {Enabled(s.Jira.Enabled), s.Jira.EmailEnv}, {Enabled(s.Jira.Enabled), s.Jira.TokenEnv}} {
			if pair.on && (pair.name == "" || os.Getenv(pair.name) == "") {
				return fmt.Errorf("target %s: an enabled destination credential is missing", t.ID)
			}
		}
		if Enabled(s.Slack.Enabled) && secureURL(os.Getenv(s.Slack.WebhookEnv)) != nil {
			return errors.New("Slack webhook must be a valid HTTPS URL")
		}
	}
	return nil
}
