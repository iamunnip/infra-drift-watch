package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHistoryReadsWithoutInventoryOrCredentials(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "receipts.json"), []byte(`{"version":2,"entries":{},"history":{"aws":[{"at":"2026-11-04T00:00:00Z","result":{"target":"aws","status":"clean","findings":[]}}]}}`), 0600)
	var out, errOut bytes.Buffer
	code := Run([]string{"history", "--state-dir", dir, "--target", "aws", "--format", "json"}, &out, &errOut)
	if code != 0 || !strings.Contains(out.String(), `"status": "clean"`) {
		t.Fatalf("code=%d out=%s err=%s", code, &out, &errOut)
	}
}
