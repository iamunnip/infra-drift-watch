// Package adapters inspects infrastructure through explicit tool executables.
package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/config"
	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

type Runner func(context.Context, string, string, []string, []string) ([]byte, int, error)
type Scanner struct{ Run Runner }

const maxOutput = 32 << 20

// limitedBuffer stops a noisy process from filling memory. Values never enter reports.
type limitedBuffer struct {
	data     bytes.Buffer
	exceeded bool
}

func (b *limitedBuffer) Bytes() []byte { return b.data.Bytes() }

func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.data.Len()+len(p) > maxOutput {
		b.exceeded = true
		return 0, errors.New("tool output exceeded limit")
	}
	return b.data.Write(p)
}
func Execute(ctx context.Context, dir, exe string, args, env []string) ([]byte, int, error) {
	cmd := exec.CommandContext(ctx, exe, args...)
	cmd.Dir = dir
	// Remove implicit argument injection and logging settings from parent environment.
	for _, v := range os.Environ() {
		key := strings.SplitN(v, "=", 2)[0]
		if strings.HasPrefix(key, "TF_CLI_ARGS") || strings.HasPrefix(key, "TF_LOG") || strings.HasPrefix(key, "TOFU_LOG") || key == "KUBECTL_EXTERNAL_DIFF" || key == "TF_WORKSPACE" {
			continue
		}
		cmd.Env = append(cmd.Env, v)
	}
	cmd.Env = append(cmd.Env, env...)
	var out limitedBuffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	cmd.WaitDelay = 2 * time.Second
	err := cmd.Run()
	if out.exceeded {
		return nil, -1, errors.New("tool output exceeded limit")
	}
	if ctx.Err() != nil {
		return nil, -1, errors.New("tool timed out or was cancelled")
	}
	if err == nil {
		return out.Bytes(), 0, nil
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return out.Bytes(), exit.ExitCode(), nil
	}
	return nil, -1, errors.New("cannot execute tool or capture bounded output")
}
func failed(t config.Target, message string) drift.Result {
	return drift.Result{Target: t.ID, Adapter: t.Adapter, Status: "failed", Findings: []drift.Finding{}, Error: message}
}
func (s Scanner) Scan(ctx context.Context, t config.Target) drift.Result {
	timeout, err := time.ParseDuration(t.Settings.Timeout)
	if err != nil || timeout <= 0 {
		return failed(t, "invalid timeout")
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	run := s.Run
	if run == nil {
		run = Execute
	}
	exe := t.Settings.Executable
	if exe == "" {
		exe = map[string]string{"opentofu": "tofu", "terraform": "terraform", "kubernetes": "kubectl"}[t.Adapter]
	}
	awsEnv := awsEnvironment(t.Settings.AWS)
	if t.Adapter == "kubernetes" {
		args := []string{"--context", t.Context}
		if t.Namespace != "" {
			args = append(args, "--namespace", t.Namespace)
		}
		manager := t.FieldManager
		if manager == "" {
			manager = "kubectl-client-side-apply"
		}
		args = append(args, "diff", "--server-side=true", "--field-manager="+manager)
		for _, p := range t.Manifests {
			args = append(args, "-f", p)
		}
		output, code, e := run(ctx, "", exe, args, append(awsEnv, "KUBECTL_EXTERNAL_DIFF=diff -U 1000000 -N"))
		if e != nil || code > 1 || code < 0 {
			return failed(t, "kubectl diff failed; check context, permissions, manifests and timeout")
		}
		r := drift.Result{Target: t.ID, Adapter: t.Adapter, Status: "clean", Findings: []drift.Finding{}}
		if code == 1 {
			r.Status = "drift"
			r.Findings, r.Digest = kubernetesDetails(output)
		}
		return suppress(t, r)
	}
	if t.Adapter != "opentofu" && t.Adapter != "terraform" {
		return failed(t, "unsupported adapter")
	}
	temp, e := os.MkdirTemp("", "infra-drift-watch-")
	if e != nil {
		return failed(t, "cannot create private plan directory")
	}
	defer os.RemoveAll(temp)
	plan := filepath.Join(temp, "scan.tfplan")
	args := []string{"plan", "-refresh-only", "-input=false", "-no-color", "-detailed-exitcode", "-lock-timeout=30s", "-out=" + plan}
	for _, p := range t.VarFiles {
		args = append(args, "-var-file="+p)
	}
	env := append(awsEnv, "TF_IN_AUTOMATION=1")
	if t.Workspace != "" {
		env = append(env, "TF_WORKSPACE="+t.Workspace)
	}
	_, code, e := run(ctx, t.Path, exe, args, env)
	if e != nil || code != 0 && code != 2 {
		return failed(t, "planning failed; check initialized workspace, credentials, state lock and timeout")
	}
	output, code, e := run(ctx, t.Path, exe, []string{"show", "-json", plan}, env)
	if e != nil || code != 0 {
		return failed(t, "cannot inspect plan JSON")
	}
	r := ParsePlan(t.ID, output)
	r.Adapter = t.Adapter
	return suppress(t, r)
}
func ParsePlan(target string, data []byte) drift.Result {
	r := drift.Result{Target: target, Status: "clean", Findings: []drift.Finding{}}
	var plan struct {
		Format   string            `json:"format_version"`
		Errored  bool              `json:"errored"`
		Complete *bool             `json:"complete"`
		Deferred []json.RawMessage `json:"deferred_changes"`
		Drift    []struct {
			Address string `json:"address"`
			Change  struct {
				Actions []string        `json:"actions"`
				Before  json.RawMessage `json:"before"`
				After   json.RawMessage `json:"after"`
			} `json:"change"`
		} `json:"resource_drift"`
	}
	if err := json.Unmarshal(data, &plan); err != nil || !strings.HasPrefix(plan.Format, "1.") || plan.Errored {
		r.Status = "failed"
		r.Error = "invalid, errored or unsupported plan JSON"
		return r
	}
	if plan.Complete != nil && !*plan.Complete || len(plan.Deferred) > 0 {
		r.Status = "incomplete"
		r.Error = "plan contains incomplete or deferred inspection"
		return r
	}
	for _, v := range plan.Drift {
		if len(v.Change.Actions) == 1 && v.Change.Actions[0] == "no-op" {
			continue
		}
		if v.Address == "" || len(v.Change.Actions) == 0 {
			r.Status = "failed"
			r.Error = "invalid drift entry"
			r.Findings = []drift.Finding{}
			return r
		}
		for _, action := range v.Change.Actions {
			if action != "update" && action != "delete" && action != "create" && action != "read" && action != "forget" {
				r.Status = "failed"
				r.Error = "unsupported drift action"
				return r
			}
		}
		r.Findings = append(r.Findings, drift.Finding{Resource: v.Address, Actions: v.Change.Actions})
	}
	if len(r.Findings) > 0 {
		canonical, _ := json.Marshal(plan.Drift)
		r.Digest = drift.Hash(string(canonical))
		r.Status = "drift"
	}
	r.Normalize()
	return r
}
func suppress(t config.Target, r drift.Result) drift.Result {
	if r.Status != "drift" {
		return r
	}
	keep := []drift.Finding{}
	for _, f := range r.Findings {
		ignored := false
		for _, s := range t.Suppressions {
			until, e := time.Parse("2006-01-02", s.Until)
			if e == nil && time.Now().UTC().Before(until.Add(24*time.Hour)) && (s.Resource == f.Resource || t.Adapter == "kubernetes" && s.Resource == "manifest-set") {
				ignored = true
				break
			}
		}
		if ignored {
			r.Suppressed++
		} else {
			keep = append(keep, f)
		}
	}
	r.Findings = keep
	if r.Suppressed > 0 {
		r.Status = "incomplete"
		r.Error = fmt.Sprintf("%d drift findings suppressed; resolution is disabled", r.Suppressed)
	}
	r.Normalize()
	return r
}

// stableDiff strips only random temporary path components from diff headers.
// Resource names and changed content remain part of the fingerprint.
func stableDiff(output []byte) string {
	lines := strings.Split(string(output), "\n")
	paths := regexp.MustCompile(`/(LIVE|MERGED)-[0-9]+`)
	for i, line := range lines {
		if strings.HasPrefix(line, "diff ") || strings.HasPrefix(line, "--- /") || strings.HasPrefix(line, "+++ /") {
			lines[i] = paths.ReplaceAllString(strings.SplitN(line, "\t", 2)[0], "/$1")
		}
	}
	return strings.Join(lines, "\n")
}
