package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
	"github.com/iamunnip/infra-drift-watch/internal/state"
)

func TestGitHubOptInCloseAndReopen(t *testing.T) {
	states := []string{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode([]map[string]any{{"number": 42, "body": Marker("aws")}})
			return
		}
		var payload map[string]any
		json.NewDecoder(r.Body).Decode(&payload)
		states = append(states, payload["state"].(string))
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	yes := true
	retry := 0
	t.Setenv("GH_CLOSE_TOKEN", "private")
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "2s", Retries: &retry, GitHub: config.GitHub{Enabled: &yes, AutoClose: &yes, Repository: "a/b", APIURL: server.URL, TokenEnv: "GH_CLOSE_TOKEN"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	e := Engine{Client: server.Client(), Store: store, Now: time.Now}
	for _, status := range []string{"drift", "failed", "incomplete", "clean", "drift"} {
		if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: status}); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(states, []string{"open", "closed", "open"}) {
		t.Fatalf("states=%v", states)
	}
}

func TestJiraOptInTransitionsAndRecurrence(t *testing.T) {
	category := "indeterminate"
	transitions := []string{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/api/3/search/jql":
			json.NewEncoder(w).Encode(map[string]any{"issues": []any{map[string]any{"key": "INFRA-1", "fields": map[string]any{"status": map[string]any{"statusCategory": map[string]any{"key": category}}}}}})
		case r.URL.Path == "/rest/api/3/issue/INFRA-1" && r.Method == "GET":
			json.NewEncoder(w).Encode(map[string]any{"fields": map[string]any{"status": map[string]any{"statusCategory": map[string]any{"key": category}}}})
		case r.URL.Path == "/rest/api/3/issue/INFRA-1" && r.Method == "PUT":
			w.WriteHeader(204)
		case r.URL.Path == "/rest/api/3/issue/INFRA-1/transitions" && r.Method == "GET":
			fmt.Fprint(w, `{"transitions":[{"id":"31","to":{"statusCategory":{"key":"done"}}},{"id":"21","to":{"statusCategory":{"key":"indeterminate"}}}]}`)
		case r.URL.Path == "/rest/api/3/issue/INFRA-1/transitions" && r.Method == "POST":
			var payload struct {
				Transition struct {
					ID string `json:"id"`
				} `json:"transition"`
			}
			json.NewDecoder(r.Body).Decode(&payload)
			transitions = append(transitions, payload.Transition.ID)
			if payload.Transition.ID == "31" {
				category = "done"
			} else {
				category = "indeterminate"
			}
			w.WriteHeader(204)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	yes := true
	retry := 0
	t.Setenv("JIRA_CLOSE_EMAIL", "user@example.com")
	t.Setenv("JIRA_CLOSE_TOKEN", "private")
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "2s", Retries: &retry, Jira: config.Jira{Enabled: &yes, AutoClose: &yes, URL: server.URL, Project: "INFRA", IssueType: "Task", EmailEnv: "JIRA_CLOSE_EMAIL", TokenEnv: "JIRA_CLOSE_TOKEN", ResolveTransition: "31", ReopenTransition: "21"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	e := Engine{Client: server.Client(), Store: store, Now: time.Now}
	for _, status := range []string{"drift", "clean", "failed", "clean", "drift"} {
		if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: status}); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(transitions, []string{"31", "21"}) {
		t.Fatalf("transitions=%v", transitions)
	}
}

func TestEnablingClosureReconcilesAlreadyResolvedIssue(t *testing.T) {
	states := []string{}
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode([]map[string]any{{"number": 42, "body": Marker("aws")}})
			return
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		states = append(states, body["state"].(string))
		fmt.Fprint(w, `{}`)
	}))
	defer server.Close()
	yes, no := true, false
	retry := 0
	t.Setenv("GH_POLICY_TOKEN", "private")
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "2s", Retries: &retry, GitHub: config.GitHub{Enabled: &yes, AutoClose: &no, Repository: "a/b", APIURL: server.URL, TokenEnv: "GH_POLICY_TOKEN"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	e := Engine{Client: server.Client(), Store: store}
	for _, status := range []string{"drift", "clean"} {
		if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: status}); err != nil {
			t.Fatal(err)
		}
	}
	target.Settings.GitHub.AutoClose = &yes
	if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: "clean"}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(states, []string{"open", "open", "closed"}) {
		t.Fatalf("states=%v", states)
	}
}

func TestJiraUnavailableTransitionPreservesPendingReceipt(t *testing.T) {
	available := false
	posts := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/api/3/search/jql":
			fmt.Fprint(w, `{"issues":[{"key":"INFRA-1","fields":{"status":{"statusCategory":{"key":"indeterminate"}}}}]}`)
		case r.URL.Path == "/rest/api/3/issue/INFRA-1" && r.Method == "GET":
			fmt.Fprint(w, `{"fields":{"status":{"statusCategory":{"key":"indeterminate"}}}}`)
		case r.Method == "PUT":
			w.WriteHeader(204)
		case r.Method == "GET":
			if available {
				fmt.Fprint(w, `{"transitions":[{"id":"31","to":{"statusCategory":{"key":"done"}}}]}`)
			} else {
				fmt.Fprint(w, `{"transitions":[]}`)
			}
		case r.Method == "POST":
			posts++
			w.WriteHeader(204)
		}
	}))
	defer server.Close()
	yes := true
	retry := 0
	t.Setenv("JIRA_PENDING_EMAIL", "user@example.com")
	t.Setenv("JIRA_PENDING_TOKEN", "private")
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "2s", Retries: &retry, Jira: config.Jira{Enabled: &yes, AutoClose: &yes, URL: server.URL, Project: "INFRA", EmailEnv: "JIRA_PENDING_EMAIL", TokenEnv: "JIRA_PENDING_TOKEN", ResolveTransition: "31", ReopenTransition: "21"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	e := Engine{Client: server.Client(), Store: store}
	if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: "drift"}); err != nil {
		t.Fatal(err)
	}
	if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: "clean"}); err == nil {
		t.Fatal("unavailable transition hidden")
	}
	for _, entry := range store.Entries {
		if !entry.Active {
			t.Fatal("failed transition acknowledged resolution")
		}
	}
	available = true
	if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: "clean"}); err != nil {
		t.Fatal(err)
	}
	if posts != 1 {
		t.Fatalf("transition posts=%d", posts)
	}
}

func TestJiraRecurrenceUsesFreshStatusRatherThanSearchIndex(t *testing.T) {
	reopened := false
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/rest/api/3/search/jql":
			fmt.Fprint(w, `{"issues":[{"key":"INFRA-1","fields":{"status":{"statusCategory":{"key":"indeterminate"}}}}]}`)
		case r.URL.Path == "/rest/api/3/issue/INFRA-1" && r.Method == "GET":
			fmt.Fprint(w, `{"fields":{"status":{"statusCategory":{"key":"done"}}}}`)
		case r.Method == "PUT":
			w.WriteHeader(204)
		case r.Method == "GET":
			fmt.Fprint(w, `{"transitions":[{"id":"21","to":{"statusCategory":{"key":"indeterminate"}}}]}`)
		case r.Method == "POST":
			reopened = true
			w.WriteHeader(204)
		}
	}))
	defer server.Close()
	yes := true
	retry := 0
	t.Setenv("JIRA_FRESH_EMAIL", "user@example.com")
	t.Setenv("JIRA_FRESH_TOKEN", "private")
	target := config.Target{ID: "aws", Settings: config.Settings{HTTPTimeout: "2s", Retries: &retry, Jira: config.Jira{Enabled: &yes, AutoClose: &yes, URL: server.URL, Project: "INFRA", EmailEnv: "JIRA_FRESH_EMAIL", TokenEnv: "JIRA_FRESH_TOKEN", ResolveTransition: "31", ReopenTransition: "21"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	e := Engine{Client: server.Client(), Store: store}
	if err := e.Deliver(context.Background(), target, drift.Result{Target: "aws", Status: "drift"}); err != nil {
		t.Fatal(err)
	}
	if !reopened {
		t.Fatal("stale search status prevented reopening")
	}
}
