// Package drift defines value-free findings and scan coverage.
package drift

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
)

type Finding struct {
	Fields   []string `json:"fields,omitempty"`
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}
type Result struct {
	Target     string    `json:"target"`
	Adapter    string    `json:"adapter"`
	Status     string    `json:"status"`
	Findings   []Finding `json:"findings"`
	Suppressed int       `json:"suppressed"`
	Error      string    `json:"error,omitempty"`
	// Digest tracks changes without persisting raw plan/diff values.
	Digest string `json:"-"`
}

func (r *Result) Normalize() {
	if r.Findings == nil {
		r.Findings = []Finding{}
	}
	for i := range r.Findings {
		sort.Strings(r.Findings[i].Actions)
		sort.Strings(r.Findings[i].Fields)
	}
	sort.Slice(r.Findings, func(i, j int) bool { return r.Findings[i].Resource < r.Findings[j].Resource })
}
func Hash(s string) string { h := sha256.Sum256([]byte(s)); return hex.EncodeToString(h[:]) }
func (r Result) Fingerprint() string {
	r.Normalize()
	b, _ := json.Marshal(struct {
		Status   string
		Findings []Finding
		Digest   string
	}{r.Status, r.Findings, r.Digest})
	return Hash(string(b))
}
