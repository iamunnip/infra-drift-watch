// Package app coordinates scanning, reporting, and optional delivery.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/adapters"
	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/delivery"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
	"github.com/iamunnip/infra-drift-watch/internal/report"
	"github.com/iamunnip/infra-drift-watch/internal/state"
)

type Options struct {
	Config, Inventory, Target, Format, Timeout, StateDir string
	DryRun                                               bool
}

func Load(o Options) (config.Config, error) {
	c, e := config.Load(o.Config, o.Inventory)
	if e != nil {
		return c, e
	}
	if o.Format != "" {
		c.Format = o.Format
	}
	if c.Format != "terminal" && c.Format != "json" {
		return c, errors.New("format must be terminal or json")
	}
	if o.StateDir != "" {
		c.StateDir = o.StateDir
	}
	if o.Timeout != "" {
		d, e := time.ParseDuration(o.Timeout)
		if e != nil || d <= 0 {
			return c, errors.New("timeout must be a positive duration")
		}
		for i := range c.Targets {
			c.Targets[i].Settings.Timeout = o.Timeout
		}
	}
	if o.Target != "" {
		selected := map[string]bool{}
		for _, id := range strings.Split(o.Target, ",") {
			selected[id] = true
		}
		filtered := []config.Target{}
		for _, t := range c.Targets {
			if selected[t.ID] {
				filtered = append(filtered, t)
				delete(selected, t.ID)
			}
		}
		if len(selected) != 0 || len(filtered) == 0 {
			return c, errors.New("target selection contains unknown or empty IDs")
		}
		c.Targets = filtered
	}
	sort.Slice(c.Targets, func(i, j int) bool { return c.Targets[i].ID < c.Targets[j].ID })
	return c, nil
}
func Scan(ctx context.Context, o Options, out io.Writer) (int, error) {
	c, e := Load(o)
	if e != nil {
		return 1, e
	}
	var store *state.Store
	enabled := false
	for _, t := range c.Targets {
		s := t.Settings
		enabled = enabled || config.Enabled(s.GitHub.Enabled) || config.Enabled(s.Slack.Enabled) || config.Enabled(s.Jira.Enabled)
	}
	if !o.DryRun && enabled {
		if e = config.ValidateSecrets(c.Targets); e != nil {
			return 1, e
		}
		store, e = state.Open(c.StateDir)
		if e != nil {
			return 1, e
		}
		defer store.Close()
	}
	engine := delivery.Engine{Store: store}
	r := report.Report{SchemaVersion: "1", Started: time.Now().UTC(), Results: []drift.Result{}, DeliveryErrors: []string{}}
	code := 0
	for _, t := range c.Targets {
		if ctx.Err() != nil {
			code = 1
			r.Results = append(r.Results, drift.Result{Target: t.ID, Adapter: t.Adapter, Status: "incomplete", Findings: []drift.Finding{}, Error: "scan cancelled"})
			continue
		}
		result := (adapters.Scanner{}).Scan(ctx, t)
		if result.Status == "failed" || result.Status == "incomplete" {
			code = 1
		} else if result.Status == "drift" && code == 0 {
			code = 2
		}
		r.Results = append(r.Results, result)
		if !o.DryRun && enabled {
			if e = engine.Deliver(ctx, t, result); e != nil {
				code = 1
				r.DeliveryErrors = append(r.DeliveryErrors, fmt.Sprintf("target %s: %s", t.ID, e))
			}
		}
	}
	if e = report.Write(out, c.Format, r); e != nil {
		return 1, errors.New("cannot write report")
	}
	return code, nil
}
