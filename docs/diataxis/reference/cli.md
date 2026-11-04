# CLI reference

```text
infra-drift-watch --help
infra-drift-watch version
infra-drift-watch validate [flags]
infra-drift-watch scan [flags]
```

| Flag | Default / meaning |
| --- | --- |
| --config | config.yaml |
| --inventory | inventory.yaml |
| --target | All targets; comma-separated exact IDs when supplied |
| --format | Inherit config; terminal or json |
| --timeout | Inherit settings; positive duration |
| --state-dir | Inherit config; explicit flag resolves from current directory |
| --dry-run | False; scan infrastructure but skip delivery and receipt writes |

Flags follow the subcommand. Validate checks file structure, settings, and target
selection; it does not contact APIs or validate credentials, executables, or backend
access. `scan --help` lists flags. Extra positional arguments are rejected.

Exit codes: 0 clean/success, 1 failure or incomplete coverage, 2 drift. Across
targets, failure wins over drift. SIGINT/SIGTERM cancel ongoing tools and HTTP work.

JSON reports use schema_version `1`, started UTC timestamp, results array, and
delivery_errors array. Each result has target, adapter, status, findings, suppressed,
and an optional safe error. Findings include resource and actions; values and
internal fingerprints are omitted. Results sort by target ID; IaC findings sort
by resource address. Reports describe selected scope, not the whole account/cluster.
