package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
	"github.com/iamunnip/infra-drift-watch/internal/state"
)

func TestJiraADFAndResolutionUpdate(t *testing.T) {
	created := false
	updates := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
			t.Error("missing Jira authentication")
		}
		switch r.URL.Path {
		case "/rest/api/3/search/jql":
			if created {
				fmt.Fprint(w, `{"issues":[{"key":"INFRA-1"}]}`)
			} else {
				fmt.Fprint(w, `{"issues":[]}`)
			}
		case "/rest/api/3/issue":
			var payload struct {
				Fields map[string]any `json:"fields"`
			}
			json.NewDecoder(r.Body).Decode(&payload)
			desc := payload.Fields["description"].(map[string]any)
			if desc["type"] != "doc" || payload.Fields["customfield_123"] != "test" {
				t.Error("missing ADF or field mapping")
			}
			created = true
			fmt.Fprint(w, `{"key":"INFRA-1"}`)
		case "/rest/api/3/issue/INFRA-1":
			updates++
			w.WriteHeader(204)
		default:
			t.Error(r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	t.Setenv("JIRA_TEST_EMAIL", "test@example.com")
	t.Setenv("JIRA_TEST_TOKEN", "private")
	yes := true
	retries := 0
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "1s", Retries: &retries, Jira: config.Jira{Enabled: &yes, URL: server.URL, Project: "INFRA", IssueType: "Task", EmailEnv: "JIRA_TEST_EMAIL", TokenEnv: "JIRA_TEST_TOKEN", Fields: map[string]any{"customfield_123": "test"}}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	engine := Engine{Client: server.Client(), Store: store, Now: time.Now}
	for _, status := range []string{"drift", "clean"} {
		if err := engine.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: status}); err != nil {
			t.Fatal(err)
		}
	}
	if !created || updates != 1 {
		t.Fatal("Jira lifecycle missing")
	}
}

func TestAmbiguousCreateFailureNotRetriedAndErrorsRedacted(t *testing.T) {
	requests := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		w.WriteHeader(500)
		fmt.Fprint(w, "private-token-response")
	}))
	defer server.Close()
	engine := Engine{Client: server.Client()}
	err := engine.request(context.Background(), "POST", server.URL+"/secret-webhook", "", map[string]string{"text": "test"}, nil, 3)
	if err == nil || strings.Contains(err.Error(), "secret-webhook") || strings.Contains(err.Error(), "private-token-response") || requests != 1 {
		t.Fatalf("requests=%d err=%v", requests, err)
	}
}
func TestRateLimitRetry(t *testing.T) {
	requests := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if requests == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		fmt.Fprint(w, "ok")
	}))
	defer server.Close()
	engine := Engine{Client: server.Client()}
	if err := engine.request(context.Background(), "POST", server.URL, "", nil, nil, 1); err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("requests=%d", requests)
	}
}
