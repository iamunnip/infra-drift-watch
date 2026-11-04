package adapters

import (
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
)

// kubernetesFindings uses only tool-generated file headers, never changed values.
// kubectl names temporary files with group.version.kind.namespace.name.
func kubernetesFindings(output []byte) []drift.Finding {
	names := map[string]bool{}
	pathPattern := regexp.MustCompile(`/(?:LIVE|MERGED)-[0-9]+/([^\s]+)`)
	identity := regexp.MustCompile(`^(?:(.+)\.)?(v[0-9]+(?:(?:alpha|beta)[0-9]+)?)\.([A-Za-z][A-Za-z0-9]*)\.([^.]*)\.(.+)$`)
	for _, line := range strings.Split(string(output), "\n") {
		if !strings.HasPrefix(line, "diff ") && !strings.HasPrefix(line, "--- /") && !strings.HasPrefix(line, "+++ /") {
			continue
		}
		for _, match := range pathPattern.FindAllStringSubmatch(line, -1) {
			filename := filepath.Base(match[1])
			parts := identity.FindStringSubmatch(filename)
			if parts == nil {
				continue
			}
			api := parts[2]
			if parts[1] != "" {
				api = parts[1] + "/" + api
			}
			namespace := parts[4]
			if namespace == "" {
				namespace = "_cluster"
			}
			names[api+"/"+parts[3]+"/"+namespace+"/"+parts[5]] = true
		}
	}
	findings := make([]drift.Finding, 0, len(names))
	for resource := range names {
		findings = append(findings, drift.Finding{Resource: resource, Actions: []string{"difference"}})
	}
	sort.Slice(findings, func(i, j int) bool { return findings[i].Resource < findings[j].Resource })
	// Preserve honest set-level coverage if an unfamiliar tool format lacks headers.
	if len(findings) == 0 {
		findings = append(findings, drift.Finding{Resource: "manifest-set", Actions: []string{"difference"}})
	}
	return findings
}
