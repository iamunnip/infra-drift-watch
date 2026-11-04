// Package delivery coordinates independent tracker and notification receipts.
package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
	"github.com/iamunnip/infra-drift-watch/internal/state"
)

type Engine struct {
	Client *http.Client
	Store  *state.Store
	Now    func() time.Time
	links  map[string]string
}

func Marker(id string) string { return "<!-- infra-drift-watch:" + drift.Hash(id) + " -->" }
func summary(t config.Target, r drift.Result) string {
	status := "Drift detected"
	if r.Status == "clean" {
		status = "Drift resolved"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%s for target %s.\n\n", status, t.ID)
	if t.Environment != "" {
		fmt.Fprintf(&b, "Environment: %s\n", safe(t.Environment))
	}
	if t.Owner != "" {
		fmt.Fprintf(&b, "Owner: %s\n", safe(t.Owner))
	}
	for _, f := range r.Findings {
		fmt.Fprintf(&b, "- %s: %s\n", safe(f.Resource), strings.Join(f.Actions, ", "))
	}
	b.WriteString("\nValues are omitted. Resolution describes the successfully scanned scope; issue closure remains manual.\n")
	return b.String()
}

// Keep untrusted names from creating mentions, links, or HTML in delivered text.
func safe(s string) string {
	s = strings.Map(func(r rune) rune {
		if r < 32 || strings.ContainsRune("@<>[]`*_", r) {
			return ' '
		}
		return r
	}, s)
	if len(s) > 300 {
		s = s[:300]
	}
	return s
}
func (e *Engine) Deliver(ctx context.Context, t config.Target, r drift.Result) error {
	if r.Status != "clean" && r.Status != "drift" {
		return nil
	}
	if e.Store == nil {
		return errors.New("delivery requires persistent state")
	}
	now := time.Now()
	if e.Now != nil {
		now = e.Now()
	}
	fingerprint := r.Fingerprint()
	e.links = map[string]string{}
	s := t.Settings
	retries := 0
	if s.Retries != nil {
		retries = *s.Retries
	}
	timeout, err := time.ParseDuration(s.HTTPTimeout)
	if err != nil {
		return errors.New("invalid delivery timeout")
	}
	type destination struct {
		name, identity string
		enabled        bool
		send           func(context.Context) error
		repeat         time.Duration
	}
	repeat, _ := time.ParseDuration(s.Slack.Repeat)
	destinations := []destination{
		{name: "github", identity: s.GitHub.APIURL + "/" + s.GitHub.Repository, enabled: config.Enabled(s.GitHub.Enabled), send: func(ctx context.Context) error { return e.github(ctx, t, r, retries) }},
		{name: "jira", identity: s.Jira.URL + "/" + s.Jira.Project, enabled: config.Enabled(s.Jira.Enabled), send: func(ctx context.Context) error { return e.jira(ctx, t, r, retries) }},
		{name: "slack", identity: os.Getenv(s.Slack.WebhookEnv), enabled: config.Enabled(s.Slack.Enabled), repeat: repeat, send: func(ctx context.Context) error {
			text := summary(t, r)
			for _, name := range []string{"github", "jira"} {
				if link := e.links[name]; link != "" {
					text += "\n" + name + ": " + link
				}
			}
			if config.Enabled(s.GitHub.Enabled) {
				text += "\nGitHub issue repository: " + safe(s.GitHub.Repository)
			}
			return e.request(ctx, http.MethodPost, os.Getenv(s.Slack.WebhookEnv), "", map[string]any{"text": text, "blocks": []any{map[string]any{"type": "section", "text": map[string]any{"type": "plain_text", "text": truncate(text, 2800)}}}}, nil, retries)
		}},
	}
	var problems []error
	for _, d := range destinations {
		if !d.enabled {
			continue
		}
		key := drift.Hash(t.ID + "\x00" + d.name + "\x00" + d.identity)
		entry := e.Store.Entries[key]
		if entry.Link != "" {
			e.links[d.name] = entry.Link
		}
		if r.Status == "clean" && !entry.Active {
			continue
		}
		if entry.Fingerprint == fingerprint && !(r.Status == "drift" && d.repeat > 0 && now.Sub(entry.Updated) >= d.repeat) {
			continue
		}
		sendCtx, cancel := context.WithTimeout(ctx, timeout)
		err := d.send(sendCtx)
		cancel()
		if err != nil {
			problems = append(problems, fmt.Errorf("%s: %w", d.name, err))
			continue
		}
		old, had := e.Store.Entries[key]
		e.Store.Entries[key] = state.Entry{Fingerprint: fingerprint, Active: r.Status == "drift", Updated: now, Link: e.links[d.name]}
		if err = e.Store.Save(); err != nil {
			if had {
				e.Store.Entries[key] = old
			} else {
				delete(e.Store.Entries, key)
			}
			problems = append(problems, fmt.Errorf("%s: cannot save acknowledgement", d.name))
		}
	}
	return errors.Join(problems...)
}
func truncate(s string, n int) string {
	if len(s) > n {
		return s[:n] + "\n[truncated]"
	}
	return s
}
