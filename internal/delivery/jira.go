package delivery

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

func (e *Engine) jira(ctx context.Context, t config.Target, r drift.Result, retries int) error {
	s := t.Settings.Jira
	base := strings.TrimRight(s.URL, "/") + "/rest/api/3"
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(os.Getenv(s.EmailEnv)+":"+os.Getenv(s.TokenEnv)))
	label := "infra-drift-watch-" + drift.Hash(t.ID)
	query := `project = "` + s.Project + `" AND labels = "` + label + `" ORDER BY created ASC`
	key := ""
	token := ""
	for page := 0; page < 1000; page++ {
		payload := map[string]any{"jql": query, "maxResults": 100, "fields": []string{"key"}}
		if token != "" {
			payload["nextPageToken"] = token
		}
		var result struct {
			Issues []struct {
				Key string `json:"key"`
			} `json:"issues"`
			Next string `json:"nextPageToken"`
		}
		if err := e.request(ctx, http.MethodPost, base+"/search/jql", auth, payload, &result, retries); err != nil {
			return err
		}
		if len(result.Issues) > 0 {
			key = result.Issues[0].Key
			break
		}
		if result.Next == "" {
			break
		}
		if result.Next == token {
			return errors.New("jira pagination did not advance")
		}
		token = result.Next
		if page == 999 {
			return errors.New("jira lookup exceeded pagination limit")
		}
	}
	if key == "" && r.Status == "clean" {
		return nil
	}
	status := "Drift"
	if r.Status == "clean" {
		status = "Resolved drift"
	}
	// Jira Cloud descriptions use Atlassian Document Format, with plain text only.
	fields := map[string]any{"summary": status + ": " + t.ID, "description": map[string]any{"type": "doc", "version": 1, "content": []any{map[string]any{"type": "paragraph", "content": []any{map[string]any{"type": "text", "text": truncate(summary(t, r), 25000)}}}}}}
	if key != "" {
		e.links["jira"] = strings.TrimRight(s.URL, "/") + "/browse/" + url.PathEscape(key)
		if err := e.request(ctx, http.MethodPut, base+"/issue/"+url.PathEscape(key), auth, map[string]any{"fields": fields}, nil, retries); err != nil {
			return err
		}
		if !config.Enabled(s.AutoClose) {
			return nil
		}
		category, err := e.jiraStatus(ctx, base+"/issue/"+url.PathEscape(key), auth, retries)
		if err != nil {
			return err
		}
		desiredDone := r.Status == "clean"
		if desiredDone && category == "done" || !desiredDone && category != "done" {
			return nil
		}
		transitionID, transitionFields := s.ResolveTransition, s.ResolveFields
		if !desiredDone {
			transitionID, transitionFields = s.ReopenTransition, s.ReopenFields
		}
		return e.jiraTransition(ctx, base, auth, key, transitionID, desiredDone, transitionFields, retries)
	}
	fields["project"] = map[string]string{"key": s.Project}
	fields["issuetype"] = map[string]string{"name": s.IssueType}
	fields["labels"] = []string{label}
	for k, v := range s.Fields {
		fields[k] = v
	}
	var created struct {
		Key string `json:"key"`
	}
	if err := e.request(ctx, http.MethodPost, base+"/issue", auth, map[string]any{"fields": fields}, &created, retries); err != nil {
		return err
	}
	if created.Key == "" {
		return errors.New("jira creation response lacks issue key")
	}
	e.links["jira"] = strings.TrimRight(s.URL, "/") + "/browse/" + url.PathEscape(created.Key)
	return nil
}

// jiraTransition checks availability and destination category before a workflow mutation.
func (e *Engine) jiraTransition(ctx context.Context, base, auth, key, id string, done bool, fields map[string]any, retries int) error {
	if id == "" {
		return errors.New("jira transition ID is required")
	}
	address := base + "/issue/" + url.PathEscape(key) + "/transitions"
	var response struct {
		Transitions []struct {
			ID string `json:"id"`
			To struct {
				Category struct {
					Key string `json:"key"`
				} `json:"statusCategory"`
			} `json:"to"`
		} `json:"transitions"`
	}
	if err := e.request(ctx, http.MethodGet, address, auth, nil, &response, retries); err != nil {
		return err
	}
	for _, transition := range response.Transitions {
		if transition.ID != id {
			continue
		}
		category := transition.To.Category.Key
		if category != "new" && category != "indeterminate" && category != "done" || (category == "done") != done {
			return errors.New("jira configured transition has an unexpected status category")
		}
		payload := map[string]any{"transition": map[string]string{"id": id}}
		if len(fields) > 0 {
			payload["fields"] = fields
		}
		return e.request(ctx, http.MethodPost, address, auth, payload, nil, retries)
	}
	return errors.New("jira configured transition is unavailable")
}

// jiraStatus avoids making workflow decisions from an eventually consistent search index.
func (e *Engine) jiraStatus(ctx context.Context, address, auth string, retries int) (string, error) {
	var issue struct {
		Fields struct {
			Status struct {
				Category struct {
					Key string `json:"key"`
				} `json:"statusCategory"`
			} `json:"status"`
		} `json:"fields"`
	}
	if err := e.request(ctx, http.MethodGet, address+"?fields=status", auth, nil, &issue, retries); err != nil {
		return "", err
	}
	category := issue.Fields.Status.Category.Key
	if category != "new" && category != "indeterminate" && category != "done" {
		return "", errors.New("jira issue status category is missing or unsupported")
	}
	return category, nil
}
