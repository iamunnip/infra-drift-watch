package report

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"testing"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }
func TestReportOutputAndFailures(t *testing.T) {
	r := Report{SchemaVersion: "1", Results: []drift.Result{{Target: "aws", Status: "failed", Error: "planning failed", Findings: []drift.Finding{}}}, DeliveryErrors: []string{"github: HTTP status 503"}}
	for _, format := range []string{"terminal", "json"} {
		var out bytes.Buffer
		if err := Write(&out, format, r); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out.String(), "planning failed") || !strings.Contains(out.String(), "HTTP status 503") {
			t.Fatal("report hid errors")
		}
		if format == "json" && !json.Valid(out.Bytes()) {
			t.Fatal("invalid JSON")
		}
		if Write(brokenWriter{}, format, r) == nil {
			t.Fatal("output error hidden")
		}
	}
}
