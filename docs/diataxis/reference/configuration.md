# Configuration and inventory reference

Both files are strict, single-document YAML. Unknown keys and duplicate YAML keys
are errors. No environment substitution is performed inside YAML strings.

## Global file

| Key | Default / behavior |
| --- | --- |
| state_dir | `.state`, relative to global config; receipt directory |
| format | `terminal` or `json`; default terminal |
| defaults | Settings inherited by all targets |

## Settings (global defaults or target settings)

| Key | Default / behavior |
| --- | --- |
| timeout | `5m`, positive per-target tool budget |
| http_timeout | `30s`, positive total budget per destination including retries |
| retries | `2`, range 0–5; explicit 429 responses, or GET 5xx only |
| executable | Adapter's default `tofu`, `terraform`, or `kubectl`; override executable path (use an absolute path or a name on PATH) |
| github.enabled | false |
| github.repository | Required owner/repository when enabled |
| github.api_url | `https://api.github.com`; HTTPS, no credentials/query/fragment |
| github.token_env | `GITHUB_TOKEN` |
| github.labels | Optional replacement list of issue labels |
| github.assignees | Optional replacement list of GitHub logins |
| slack.enabled | false |
| slack.webhook_env | `SLACK_WEBHOOK_URL`; variable contains HTTPS webhook URL |
| slack.repeat | `24h`, positive unchanged-drift reminder interval |
| jira.enabled | false |
| jira.url | Required HTTPS Jira Cloud site when enabled |
| jira.project | Required uppercase project key |
| jira.issue_type | `Task` |
| jira.email_env | `JIRA_EMAIL` |
| jira.token_env | `JIRA_TOKEN` |
| jira.fields | Replacement map of `customfield_*` fields, applied at creation |

Nonempty scalars override inherited values. Explicit booleans and retry zero
are respected. Nested settings merge field by field. Lists/maps replace entirely;
`[]` or `{}` clears inherited collections. Null and empty strings inherit.
There is no automatic issue closure setting in this release.

## Inventory targets

| Key | Meaning |
| --- | --- |
| id | Required stable unique identity: letters, numbers, dot, underscore, hyphen; max 100 |
| adapter | Required opentofu, terraform, or kubernetes |
| description, owner, environment | Optional metadata; owner/environment appear in delivery |
| path | Required initialized workspace directory for state tools |
| workspace | Optional TF_WORKSPACE selection; otherwise tool default |
| var_files | Variable files resolved by the infrastructure executable relative to path |
| context | Explicit Kubernetes context, required for Kubernetes |
| field_manager | Kubernetes ownership identity; default kubectl-client-side-apply; align with the deploying manager |
| namespace | Optional namespace; manifest namespaces and kubectl scope rules still apply |
| manifests | Required local paths for Kubernetes; relative to inventory |
| suppressions | List of resource, reason, until (YYYY-MM-DD, inclusive UTC) |
| settings | Overrides for this target |

Use `manifest-set` to suppress the complete Kubernetes comparison. Suppressions
are exact resource identities, not globs. Any active matching suppression produces
an incomplete result and prevents delivery/resolution for the target. Expired
suppressions no longer hide findings. This conservative policy avoids announcing
resolution merely because an exception hid remaining drift.

Targets scan serially; delivery is also serial. Timeouts and retry counts are
configurable; concurrency is intentionally one in this release. Raw values,
arbitrary shell hooks, automatic remediation, and unsafe TLS are not options.
