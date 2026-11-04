package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	for _, tt := range []struct {
		name     string
		args     []string
		code     int
		out, err string
	}{
		{"default help", nil, 0, "Usage:", ""},
		{"help", []string{"--help"}, 0, "Usage:", ""},
		{"version", []string{"version"}, 0, "infra-drift-watch dev", ""},
		{"unknown command", []string{"unknown"}, 1, "", "unknown command"},
		{"extra arguments", []string{"version", "extra"}, 1, "", "unexpected arguments"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var out, err bytes.Buffer
			if got := Run(tt.args, &out, &err); got != tt.code {
				t.Fatalf("exit = %d, want %d", got, tt.code)
			}
			if !strings.Contains(out.String(), tt.out) || !strings.Contains(err.String(), tt.err) {
				t.Fatalf("stdout=%q stderr=%q", out.String(), err.String())
			}
			if tt.code != 0 && out.Len() != 0 {
				t.Fatal("failure wrote to stdout")
			}
			if tt.code == 0 && err.Len() != 0 {
				t.Fatal("success wrote to stderr")
			}
		})
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

func TestRunOutputFailure(t *testing.T) {
	for _, args := range [][]string{nil, {"version"}} {
		if code := Run(args, failingWriter{}, io.Discard); code != 1 {
			t.Fatalf("Run(%v) = %d, want 1 for failed output", args, code)
		}
	}
}
