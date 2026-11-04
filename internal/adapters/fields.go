package adapters

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/iamunnip/infra-drift-watch/internal/drift"
	"go.yaml.in/yaml/v3"
)

// fieldChange is ephemeral. Values are hashed and never included in reports or state.
type fieldChange struct {
	Path                        string
	Before, After               any
	BeforePresent, AfterPresent bool
}

func kubernetesDetails(output []byte) ([]drift.Finding, string) {
	findings := kubernetesFindings(output)
	changes := map[string][]fieldChange{}
	lines := strings.Split(string(output), "\n")
	for i := 0; i+2 < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "--- /") || !strings.HasPrefix(lines[i+1], "+++ /") {
			continue
		}
		end := i + 2
		for end < len(lines) && !strings.HasPrefix(lines[end], "diff ") && !strings.HasPrefix(lines[end], "--- /") {
			end++
		}
		before, after, ok := unifiedObjects(lines[i+2 : end])
		if !ok {
			continue
		}
		resource := kubernetesFindings([]byte(lines[i] + "\n" + lines[i+1]))
		if len(resource) != 1 || resource[0].Resource == "manifest-set" {
			continue
		}
		diff := []fieldChange{}
		compareFields(before, after, "", &diff)
		if len(diff) > 0 {
			changes[resource[0].Resource] = diff
		}
		i = end - 1
	}
	complete := true
	for i := range findings {
		list, ok := changes[findings[i].Resource]
		if !ok {
			complete = false
			continue
		}
		for _, change := range list {
			findings[i].Fields = append(findings[i].Fields, change.Path)
		}
		sort.Strings(findings[i].Fields)
	}
	if complete {
		encoded, err := json.Marshal(changes)
		if err == nil {
			return findings, drift.Hash(string(encoded))
		}
	}
	// Unknown/redacted output still establishes drift, without invented field details.
	relevant := []string{}
	for _, line := range strings.Split(stableDiff(output), "\n") {
		if strings.HasPrefix(line, "diff ") || strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") || strings.HasPrefix(line, "< ") || strings.HasPrefix(line, "> ") {
			relevant = append(relevant, line)
		}
	}
	if len(relevant) == 0 {
		return findings, drift.Hash(stableDiff(output))
	}
	return findings, drift.Hash(strings.Join(relevant, "\n"))
}

// unifiedObjects accepts one full-prefix hunk. An unfamiliar/truncated hunk safely
// falls back to resource-level findings. Huge context is requested from GNU diff.
func unifiedObjects(lines []string) (any, any, bool) {
	header := regexp.MustCompile(`^@@ -([0-9]+)(?:,([0-9]+))? \+([0-9]+)(?:,([0-9]+))? @@`)
	var before, after strings.Builder
	oldCount, newCount := 0, 0
	wantOld, wantNew := 0, 0
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			match := header.FindStringSubmatch(line)
			if found || match == nil || match[1] != "0" && match[1] != "1" || match[3] != "0" && match[3] != "1" {
				return nil, nil, false
			}
			found = true
			wantOld, wantNew = 1, 1
			if match[2] != "" {
				wantOld, _ = strconv.Atoi(match[2])
			}
			if match[4] != "" {
				wantNew, _ = strconv.Atoi(match[4])
			}
			continue
		}
		if !found || line == "" || strings.HasPrefix(line, `\ No newline`) {
			continue
		}
		switch line[0] {
		case ' ':
			before.WriteString(line[1:] + "\n")
			after.WriteString(line[1:] + "\n")
			oldCount++
			newCount++
		case '-':
			before.WriteString(line[1:] + "\n")
			oldCount++
		case '+':
			after.WriteString(line[1:] + "\n")
			newCount++
		default:
			return nil, nil, false
		}
	}
	if !found || oldCount != wantOld || newCount != wantNew {
		return nil, nil, false
	}
	var oldObject, newObject any
	if yaml.Unmarshal([]byte(before.String()), &oldObject) != nil || yaml.Unmarshal([]byte(after.String()), &newObject) != nil {
		return nil, nil, false
	}
	return oldObject, newObject, true
}

func pointerKey(key string) string {
	return strings.ReplaceAll(strings.ReplaceAll(key, "~", "~0"), "/", "~1")
}
func compareFields(before, after any, path string, changes *[]fieldChange) {
	if reflect.DeepEqual(before, after) {
		return
	}
	left, lok := before.(map[string]any)
	right, rok := after.(map[string]any)
	if lok && rok {
		keys := map[string]bool{}
		for key := range left {
			keys[key] = true
		}
		for key := range right {
			keys[key] = true
		}
		ordered := []string{}
		for key := range keys {
			ordered = append(ordered, key)
		}
		sort.Strings(ordered)
		for _, key := range ordered {
			a, ap := left[key]
			b, bp := right[key]
			next := path + "/" + pointerKey(key)
			if ap != bp {
				*changes = append(*changes, fieldChange{next, a, b, ap, bp})
			} else {
				compareFields(a, b, next, changes)
			}
		}
		return
	}
	a, aok := before.([]any)
	b, bok := after.([]any)
	if aok && bok {
		for i := 0; i < len(a) || i < len(b); i++ {
			next := fmt.Sprintf("%s/%d", path, i)
			if i >= len(a) {
				*changes = append(*changes, fieldChange{next, nil, b[i], false, true})
			} else if i >= len(b) {
				*changes = append(*changes, fieldChange{next, a[i], nil, true, false})
			} else {
				compareFields(a[i], b[i], next, changes)
			}
		}
		return
	}
	if path == "" {
		path = "$"
	}
	*changes = append(*changes, fieldChange{path, before, after, true, true})
}
