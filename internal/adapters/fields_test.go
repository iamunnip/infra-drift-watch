package adapters

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestKubernetesChangedFieldsOmitValuesAndIgnoreUnchangedMetadata(t *testing.T) {
	output := `diff -U 1000000 -N /tmp/LIVE-1/v1.ConfigMap.default.demo /tmp/MERGED-2/v1.ConfigMap.default.demo
--- /tmp/LIVE-1/v1.ConfigMap.default.demo	2026-11-04 10:00:00
+++ /tmp/MERGED-2/v1.ConfigMap.default.demo	2026-11-04 10:00:01
@@ -1,9 +1,9 @@
 apiVersion: v1
 kind: ConfigMap
 metadata:
   name: demo
   namespace: default
   resourceVersion: "10"
 data:
-  greeting: private-before
+  greeting: private-after
   stable: unchanged
`
	findings, digest := kubernetesDetails([]byte(output))
	if len(findings) != 1 || !reflect.DeepEqual(findings[0].Fields, []string{"/data/greeting"}) {
		t.Fatalf("%+v", findings)
	}
	data, _ := json.Marshal(findings)
	if strings.Contains(string(data), "private-") {
		t.Fatal("changed values leaked")
	}
	updated := strings.ReplaceAll(output, `resourceVersion: "10"`, `resourceVersion: "11"`)
	updated = strings.ReplaceAll(strings.ReplaceAll(updated, "LIVE-1/", "LIVE-99/"), "MERGED-2/", "MERGED-77/")
	_, next := kubernetesDetails([]byte(updated))
	if digest != next {
		t.Fatal("unchanged metadata or temporary names changed finding fingerprint")
	}
}

func TestChangedFieldPointersHandleKeysArraysAndMissingNull(t *testing.T) {
	changes := []fieldChange{}
	before := map[string]any{"a/b": map[string]any{"~key": []any{"first", "old"}}, "present": nil}
	after := map[string]any{"a/b": map[string]any{"~key": []any{"first", "new"}}}
	compareFields(before, after, "", &changes)
	paths := []string{}
	for _, change := range changes {
		paths = append(paths, change.Path)
	}
	if !reflect.DeepEqual(paths, []string{"/a~1b/~0key/1", "/present"}) {
		t.Fatalf("paths=%v", paths)
	}
}

func TestPartialKubernetesHunksKeepHonestResourceFallback(t *testing.T) {
	output := "--- /tmp/LIVE-1/v1.Secret.default.example\n+++ /tmp/MERGED-2/v1.Secret.default.example\n@@ -10 +10 @@\n- old\n+ new\n"
	findings, _ := kubernetesDetails([]byte(output))
	if len(findings) != 1 || findings[0].Resource != "v1/Secret/default/example" || len(findings[0].Fields) != 0 {
		t.Fatalf("invented partial field details: %+v", findings)
	}
}
