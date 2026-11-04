// Run from the repository root with go run ./scripts/check-docs.go.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"go.yaml.in/yaml/v3"
)

func main() {
	scope := flag.String("scope", "all", "all, docs, or knowledge")
	flag.Parse()
	if *scope != "all" && *scope != "docs" && *scope != "knowledge" {
		fmt.Fprintln(os.Stderr, "invalid documentation scope")
		os.Exit(1)
	}
	var failures []string
	links := regexp.MustCompile(`\]\(([^)]+)\)`)
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" || d.Name() == "node_modules" || d.Name() == ".tools" {
				return filepath.SkipDir
			}
			return nil
		}
		if *scope != "knowledge" && strings.HasPrefix(path, ".github/workflows/") && strings.HasSuffix(path, ".yml") {
			data, e := os.ReadFile(path)
			if e != nil {
				return e
			}
			var v any
			if e = yaml.Unmarshal(data, &v); e != nil {
				failures = append(failures, path+": invalid workflow YAML")
			}
		}
		if *scope == "knowledge" && !strings.HasPrefix(path, "knowledge/") || *scope == "docs" && strings.HasPrefix(path, "knowledge/") {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		data, e := os.ReadFile(path)
		if e != nil {
			return e
		}
		for _, m := range links.FindAllSubmatch(data, -1) {
			target := string(m[1])
			if strings.Contains(target, "://") || strings.HasPrefix(target, "#") {
				continue
			}
			target = strings.SplitN(target, "#", 2)[0]
			if filepath.IsAbs(target) {
				failures = append(failures, path+": absolute documentation link")
				continue
			}
			if _, e := os.Stat(filepath.Join(filepath.Dir(path), target)); e != nil {
				failures = append(failures, path+": broken link "+target)
			}
		}
		if strings.HasPrefix(path, "knowledge/") {
			name := filepath.Base(path)
			if name == "log.md" || name == "index.md" && path != "knowledge/index.md" {
				return nil
			}
			parts := bytes.SplitN(data, []byte("---"), 3)
			metadata := map[string]any{}
			if len(parts) != 3 || len(bytes.TrimSpace(parts[0])) != 0 || yaml.Unmarshal(parts[1], &metadata) != nil {
				failures = append(failures, path+": invalid OKF frontmatter")
				return nil
			}
			if path == "knowledge/index.md" {
				if metadata["okf_version"] != "0.2" {
					failures = append(failures, path+": expected okf_version 0.2")
				}
			} else {
				for _, key := range []string{"type", "title", "description", "tags"} {
					if metadata[key] == nil {
						failures = append(failures, path+": missing "+key)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		failures = append(failures, err.Error())
	}
	if len(failures) > 0 {
		for _, f := range failures {
			fmt.Fprintln(os.Stderr, f)
		}
		os.Exit(1)
	}
	fmt.Printf("Documentation checks pass (scope: %s)\n", *scope)
}
