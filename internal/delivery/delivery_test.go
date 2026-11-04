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

func TestIndependentDeliveryAndRetry(t *testing.T) {
	githubCalls, slackCalls := 0, 0
	failSlack := true
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasSuffix(r.URL.Path, "/issues") && r.Method == "GET":
			fmt.Fprint(w, `[]`)
		case strings.HasSuffix(r.URL.Path, "/issues") && r.Method == "POST":
			githubCalls++
			json.NewEncoder(w).Encode(map[string]any{"number": 1})
		case r.URL.Path == "/slack":
			slackCalls++
			if failSlack {
				w.WriteHeader(500)
			} else {
				fmt.Fprint(w, "ok")
			}
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	t.Setenv("TEST_GH", "token")
	t.Setenv("TEST_SLACK", server.URL+"/slack")
	yes := true
	retries := 0
	target := config.Target{ID: "aws-prod", Settings: config.Settings{HTTPTimeout: "1s", Retries: &retries, GitHub: config.GitHub{Enabled: &yes, Repository: "a/b", APIURL: server.URL, TokenEnv: "TEST_GH"}, Slack: config.Slack{Enabled: &yes, WebhookEnv: "TEST_SLACK", Repeat: "24h"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	engine := Engine{Client: server.Client(), Store: store, Now: time.Now}
	result := drift.Result{Target: "aws-prod", Status: "drift", Findings: []drift.Finding{{Resource: "bucket", Actions: []string{"update"}}}}
	if err := engine.Deliver(context.Background(), target, result); err == nil {
		t.Fatal("Slack failure hidden")
	}
	failSlack = false
	if err := engine.Deliver(context.Background(), target, result); err != nil {
		t.Fatal(err)
	}
	if githubCalls != 1 || slackCalls != 2 {
		t.Fatalf("github=%d slack=%d", githubCalls, slackCalls)
	}
	if err := engine.Deliver(context.Background(), target, result); err != nil {
		t.Fatal(err)
	}
	if slackCalls != 2 {
		t.Fatal("duplicate unchanged notification")
	}
	result.Status = "failed"
	if err := engine.Deliver(context.Background(), target, result); err != nil {
		t.Fatal(err)
	}
	if slackCalls != 2 {
		t.Fatal("failed scan resolved finding")
	}
}

func TestGitHubRecoveryAndResolution(t *testing.T) {
	patches := 0
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			json.NewEncoder(w).Encode([]map[string]any{{"number": 42, "body": Marker("target"), "state": "closed"}})
			return
		}
		if r.Method == "PATCH" {
			patches++
			var payload map[string]any
			json.NewDecoder(r.Body).Decode(&payload)
			if payload["state"] != "open" {
				t.Error("default must retain open issue")
			}
			fmt.Fprint(w, `{}`)
			return
		}
		t.Errorf("unexpected request %s", r.Method)
	}))
	defer server.Close()
	t.Setenv("GH_TEST", "token")
	yes := true
	retries := 0
	target := config.Target{ID: "target", Settings: config.Settings{HTTPTimeout: "1s", Retries: &retries, GitHub: config.GitHub{Enabled: &yes, Repository: "a/b", APIURL: server.URL, TokenEnv: "GH_TEST"}}}
	store, err := state.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	engine := Engine{Client: server.Client(), Store: store, Now: time.Now}
	for _, status := range []string{"drift", "clean"} {
		if err := engine.Deliver(context.Background(), target, drift.Result{Target: "target", Status: status}); err != nil {
			t.Fatal(err)
		}
	}
	if patches != 2 {
		t.Fatalf("patches=%d", patches)
	}
}
