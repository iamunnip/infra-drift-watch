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
		return e.request(ctx, http.MethodPut, base+"/issue/"+url.PathEscape(key), auth, map[string]any{"fields": fields}, nil, retries)
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
