package delivery

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

func (e *Engine) github(ctx context.Context, t config.Target, r drift.Result, retries int) error {
	s := t.Settings.GitHub
	auth := "Bearer " + os.Getenv(s.TokenEnv)
	base := strings.TrimRight(s.APIURL, "/") + "/repos/" + s.Repository + "/issues"
	type issue struct {
		Number      int    `json:"number"`
		Body        string `json:"body"`
		PullRequest any    `json:"pull_request"`
	}
	number := 0
	// List the repository directly rather than relying on delayed search indexes.
	for page := 1; page <= 1000; page++ {
		var list []issue
		if err := e.request(ctx, http.MethodGet, fmt.Sprintf("%s?state=all&per_page=100&page=%d", base, page), auth, nil, &list, retries); err != nil {
			return err
		}
		for _, i := range list {
			if i.PullRequest == nil && strings.Contains(i.Body, Marker(t.ID)) {
				number = i.Number
				break
			}
		}
		if number != 0 || len(list) < 100 {
			break
		}
		if page == 1000 {
			return errors.New("GitHub issue lookup exceeded pagination limit")
		}
	}
	if number == 0 && r.Status == "clean" {
		return nil
	}
	status := "Drift"
	if r.Status == "clean" {
		status = "Resolved drift"
	}
	payload := map[string]any{"title": status + ": " + t.ID, "body": truncate(summary(t, r), 50000) + "\n" + Marker(t.ID), "state": "open"}
	if s.Labels != nil {
		payload["labels"] = s.Labels
	}
	if s.Assignees != nil {
		payload["assignees"] = s.Assignees
	}
	if number == 0 {
		delete(payload, "state")
		var created struct {
			Number int `json:"number"`
		}
		if err := e.request(ctx, http.MethodPost, base, auth, payload, &created, retries); err != nil {
			return err
		}
		if created.Number <= 0 {
			return errors.New("GitHub creation response lacks issue number")
		}
		number = created.Number
	} else if err := e.request(ctx, http.MethodPatch, fmt.Sprintf("%s/%d", base, number), auth, payload, nil, retries); err != nil {
		return err
	}
	web := strings.TrimSuffix(strings.TrimRight(s.APIURL, "/"), "/api/v3")
	if web == "https://api.github.com" {
		web = "https://github.com"
	}
	e.links["github"] = fmt.Sprintf("%s/%s/issues/%d", web, s.Repository, number)
	return nil
}
