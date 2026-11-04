// Package report renders stable, value-free scan reports.
package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

type Report struct {
	SchemaVersion  string         `json:"schema_version"`
	Started        time.Time      `json:"started"`
	Results        []drift.Result `json:"results"`
	DeliveryErrors []string       `json:"delivery_errors"`
}

func Write(out io.Writer, format string, r Report) error {
	if format == "json" {
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		return enc.Encode(r)
	}
	for _, v := range r.Results {
		if _, e := fmt.Fprintf(out, "%s [%s]: %s (%d findings, %d suppressed)\n", v.Target, v.Adapter, v.Status, len(v.Findings), v.Suppressed); e != nil {
			return e
		}
		for _, f := range v.Findings {
			if _, e := fmt.Fprintf(out, "  %s: %v\n", f.Resource, f.Actions); e != nil {
				return e
			}
		}
		if v.Error != "" {
			if _, e := fmt.Fprintln(out, "  "+v.Error); e != nil {
				return e
			}
		}
	}
	for _, e := range r.DeliveryErrors {
		if _, err := fmt.Fprintln(out, "delivery failed: "+e); err != nil {
			return err
		}
	}
	return nil
}
