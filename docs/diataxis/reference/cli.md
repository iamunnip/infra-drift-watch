# CLI reference

```text
infra-drift-watch --help
infra-drift-watch version
infra-drift-watch validate [flags]
infra-drift-watch scan [flags]
infra-drift-watch history [--state-dir .state] [--target ID] [--limit 20] [--format terminal|json]
```

| Flag | Default / meaning |
| --- | --- |
| --config | config.yaml |
| --inventory | inventory.yaml |
| --target | All targets; comma-separated exact IDs when supplied |
| --format | Inherit config; terminal or json |
| --timeout | Inherit settings; positive duration |
| --state-dir | Inherit config; explicit flag resolves from current directory |
| --dry-run | False; scan infrastructure but skip delivery, receipt and history writes |

Flags follow the subcommand. Validate checks file structure, settings, and target
selection; it does not contact APIs or validate credentials, executables, or backend
access. `scan --help` lists flags. Extra positional arguments are rejected.

Exit codes: 0 clean/success, 1 failure or incomplete coverage, 2 drift. Across
targets, failure wins over drift. SIGINT/SIGTERM cancel ongoing tools and HTTP work.

JSON reports use schema_version `1`, started UTC timestamp, results array, and
delivery_errors array. Each result has target, adapter, status, findings, suppressed,
and an optional safe error. Findings include resource, actions, and optional changed field paths; values and
internal fingerprints are omitted. Kubernetes field paths use slash-separated
JSON Pointer escaping (~0 for tilde, ~1 for slash); `$` means the complete object.
Unfamiliar diff output falls back to resource or manifest-set detail without
inventing field paths. Optional history_errors report persistence failures. Results sort by target ID; IaC findings sort
by resource address. Reports describe selected scope, not the whole account/cluster.

## History command

History reads an existing state file without inventory, credentials, destination
calls, or writes. --state-dir defaults to .state relative to the current directory;
when your config uses another directory, supply it explicitly. --target selects
one exact ID, --limit is 1–10000 (default 20), and --format is terminal or json.
Records are newest first, with target ID breaking timestamp ties. Missing or
unsupported state files return 1 without creating a directory.

History JSON schema_version is 1. Records contain at (UTC), duration_ms (inspection
time), result, delivery_errors and an opaque fingerprint; no raw plan/diff values.
State storage format is version 2, which is distinct from the public report schema.
