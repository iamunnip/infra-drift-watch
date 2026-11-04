# Configuration and inventory reference

Both files are strict, single-document YAML. Unknown keys and duplicate YAML keys
are errors. No environment substitution is performed inside YAML strings.

## Global file

| Key | Default / behavior |
| --- | --- |
| state_dir | `.state`, relative to global config; receipts and history directory |
| format | `terminal` or `json`; default terminal |
| defaults | Settings inherited by all targets |

## Settings (global defaults or target settings)

| Key | Default / behavior |
| --- | --- |
| aws.profile | Omitted inherits runtime AWS selection; `default` selects your shared default profile; explicit empty clears inherited profile selectors |
| aws.region | Omitted inherits runtime AWS region; e.g. `us-east-1`; explicit empty clears inherited region selectors |
| timeout | `5m`, positive per-target tool budget |
| http_timeout | `30s`, positive total budget per destination including retries |
| retries | `2`, range 0–5; explicit 429 responses, or GET 5xx only |
| executable | Adapter's default `tofu`, `terraform`, or `kubectl`; override executable path (use an absolute path or a name on PATH) |
| history.enabled | true; persist normal scan records, including failed/incomplete results |
| history.limit | 100 records per target; integer 1–10000 |
| github.enabled | false |
| github.auto_close | false; close the marked issue only on complete clean scans; reopen on recurrence |
| github.repository | Required owner/repository when enabled |
| github.api_url | `https://api.github.com`; HTTPS, no credentials/query/fragment |
| github.token_env | `GITHUB_TOKEN` |
| github.labels | Optional replacement list of issue labels |
| github.assignees | Optional replacement list of GitHub logins |
| slack.enabled | false |
| slack.webhook_env | `SLACK_WEBHOOK_URL`; variable contains HTTPS webhook URL |
| slack.repeat | `24h`, positive unchanged-drift reminder interval |
| jira.enabled | false |
| jira.auto_close | false; opt-in resolve/reopen transitions on existing marked issues |
| jira.resolve_transition | Project-specific transition ID string; required with auto_close |
| jira.reopen_transition | Project-specific transition ID string; required with auto_close |
| jira.resolve_fields | Optional replacement map of fields needed by the resolve transition, such as resolution |
| jira.reopen_fields | Optional replacement map of fields needed by the reopen transition |
| jira.url | Required HTTPS Jira Cloud site when enabled |
| jira.project | Required uppercase project key |
| jira.issue_type | `Task` |
| jira.email_env | `JIRA_EMAIL` |
| jira.token_env | `JIRA_TOKEN` |
| jira.fields | Replacement map of `customfield_*` fields, applied at creation |

Nonempty scalars override inherited values. Explicit booleans and retry zero
are respected. Nested settings merge field by field. Lists/maps replace entirely;
`[]` or `{}` clears inherited collections. Null and ordinary empty strings inherit. AWS profile/region are the exception:
an explicit empty string clears the selection, while omission or null inherits.
Closure remains manual unless explicitly enabled. Jira transition destinations
are checked against done (resolution) or non-done (recurrence) status categories.

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

Kubernetes identities use apiVersion/Kind/namespace/name, such as
`v1/ConfigMap/default/demo` or `apps/v1/Deployment/production/web`. Cluster-scoped
resources use `_cluster` instead of namespace.
Use `manifest-set` to suppress the complete Kubernetes comparison, including
detailed resource findings, preserving older configurations. Suppressions
are exact resource identities, not globs. Any active matching suppression produces
an incomplete result and prevents delivery/resolution for the target. Expired
suppressions no longer hide findings. This conservative policy avoids announcing
resolution merely because an exception hid remaining drift.

Targets scan serially; delivery is also serial. Timeouts and retry counts are
configurable; concurrency is intentionally one in this release. Raw values,
arbitrary shell hooks, automatic remediation, and unsafe TLS are not options.

## AWS selector precedence

Built-in AWS selectors are unspecified. Global defaults merge into target settings;
explicit CLI --aws-profile/--aws-region overrides then apply to selected targets.
Profile and region are scoped to tool subprocesses. Profile populates AWS_PROFILE
and AWS_DEFAULT_PROFILE; region populates AWS_REGION and AWS_DEFAULT_REGION.
The parent environment and existing credential variables are never rewritten.
Kubernetes credential plugins inherit these selectors too, which supports EKS.

AWS credentials are not accepted in this configuration. Shared profiles, SSO,
environment credentials, web identity and role providers remain the responsibility
of the invoked tool's AWS SDK. Explicit provider/backend configuration can take
precedence over environment selectors. Validation checks selector syntax only;
it does not inspect credentials, verify a profile exists, or contact AWS.
