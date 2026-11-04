// Package cli implements command dispatch without terminating the process.
package cli

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/iamunnip/infra-drift-watch/internal/app"
)

// Version is set at release build time using -ldflags.
var Version = "dev"

// Run writes command output and returns a process exit code.
func Run(args []string, out, errOut io.Writer) int {
	return RunContext(context.Background(), args, out, errOut)
}

// RunContext allows shutdown signals to cancel scans and HTTP requests.
func RunContext(ctx context.Context, args []string, out, errOut io.Writer) int {
	if len(args) == 0 {
		return help(out)
	}
	switch args[0] {
	case "history":
		return history(args[1:], out, errOut)
	case "scan", "validate":
		return configured(ctx, args, out, errOut)
	case "help", "--help", "-h":
		if len(args) != 1 {
			return failure(errOut, "unexpected arguments after help")
		}
		return help(out)
	case "version":
		if len(args) != 1 {
			return failure(errOut, "unexpected arguments after version")
		}
		if _, err := fmt.Fprintln(out, "infra-drift-watch "+Version); err != nil {
			return failure(errOut, "cannot write version output")
		}
		return 0
	default:
		return failure(errOut, fmt.Sprintf("unknown command %q; use --help", args[0]))
	}
}

func help(out io.Writer) int {
	_, err := fmt.Fprintln(out, `infra-drift-watch: infrastructure drift detection

Usage:
  infra-drift-watch --help
  infra-drift-watch version
  infra-drift-watch history --state-dir .state [--target ID] [--format json]
  infra-drift-watch validate --config config.yaml --inventory inventory.yaml
  infra-drift-watch scan --config config.yaml --inventory inventory.yaml [--dry-run]

Use scan --help for configuration, selection, timeout, and output options.
Dry-run inspects infrastructure but skips delivery and state writes.`)
	if err != nil {
		return 1
	}
	return 0
}

func failure(out io.Writer, message string) int {
	fmt.Fprintln(out, "error:", message)
	return 1
}

func configured(ctx context.Context, args []string, out, errOut io.Writer) int {
	var o app.Options
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(errOut)
	flags.StringVar(&o.Config, "config", "config.yaml", "global YAML configuration")
	flags.StringVar(&o.Inventory, "inventory", "inventory.yaml", "inventory YAML")
	flags.StringVar(&o.Target, "target", "", "comma-separated target IDs")
	flags.StringVar(&o.Format, "format", "", "terminal or json")
	flags.StringVar(&o.Timeout, "timeout", "", "per-target scan timeout override")
	flags.StringVar(&o.StateDir, "state-dir", "", "delivery state directory override")
	flags.BoolVar(&o.DryRun, "dry-run", false, "inspect infrastructure without delivery or receipt writes")
	if err := flags.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if flags.NArg() != 0 {
		return failure(errOut, "unexpected positional arguments")
	}
	if args[0] == "validate" {
		c, err := app.Load(o)
		if err != nil {
			return failure(errOut, err.Error())
		}
		if _, err = fmt.Fprintf(out, "Configuration valid: %d target(s)\n", len(c.Targets)); err != nil {
			return failure(errOut, "cannot write validation output")
		}
		return 0
	}
	code, err := app.Scan(ctx, o, out)
	if err != nil {
		return failure(errOut, err.Error())
	}
	return code
}

func history(args []string, out, errOut io.Writer) int {
	flags := flag.NewFlagSet("history", flag.ContinueOnError)
	flags.SetOutput(errOut)
	dir := flags.String("state-dir", ".state", "existing state directory")
	target := flags.String("target", "", "exact target ID; all targets by default")
	format := flags.String("format", "terminal", "terminal or json")
	limit := flags.Int("limit", 20, "latest records to return, 1-10000")
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}
	if flags.NArg() != 0 {
		return failure(errOut, "unexpected positional arguments")
	}
	if err := app.History(*dir, *target, *format, *limit, out); err != nil {
		return failure(errOut, err.Error())
	}
	return 0
}
