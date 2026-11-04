# Verification and compatibility

| Component | Pinned test/runtime baseline | Verification scope |
| --- | --- | --- |
| Go | 1.27.3 | Unit, race, vet and binary/container builds |
| OpenTofu | 1.13.1 | Real disposable local workspace; plan fixtures |
| Terraform | 1.13.5 | Real disposable local workspace; plan fixtures |
| kubectl | 1.36.6 | Command boundary tests and disposable cluster test |
| kind | 0.33.0 | Optional Kubernetes test prerequisite |
| Kubernetes | kind node 1.36.4 | Disposable ConfigMap clean/change/resolution |
| GitHub, Slack, Jira | Documented REST endpoints | Mocked TLS-server tests; no live account verification |
| AWS | AWS provider ~> 6.0 example | Configuration example only; no live AWS scan verified |

Pinned baselines are compatibility fixtures, not a promise that older binaries
receive upstream security support. Review upstream updates before production use.
Only Linux amd64 execution is verified locally; arm64 release artifacts are built
but need deployment testing. Plan format major 1 is accepted; malformed or errored
plans fail, and explicitly incomplete/deferred plans produce incomplete results. Version combinations outside this table are unverified.

`make check` runs local checks after dependencies are present. `make tools` installs
runtime executables explicitly. `make test-integration` needs those tools. Kubernetes
tests need kind/kubectl/Docker and create/delete their own named cluster using a
temporary kubeconfig. Tests apply only disposable fixtures, never your AWS resources.

For Mermaid run npm ci --prefix scripts/mermaid --ignore-scripts then make
mermaid-check. `make dev-tools` installs pinned govulncheck/actionlint; `make vuln`
and `make workflow-check` execute them. CI invokes these explicitly. Live cloud,
backend, permissions, webhooks and Jira field availability need deployment checks.

## Local verification record (2026-11-04)

Passed: make check, race-enabled tests, OpenSpec validation, actionlint, five
Mermaid diagrams parsed with the real parser, and govulncheck (no vulnerabilities
reported). Real OpenTofu and Terraform tests used the local provider 2.5.3 to prove
clean/deleted/resolved behavior without persisting refreshed state; built-in
terraform_data tests proved pending configuration is not observed drift.

The disposable Kubernetes test passed clean, changed, repeated unchanged drift,
resolved, and conflicting-owner cases. Its cluster and kubeconfig are temporary.
GitHub/Jira/Slack tests use local TLS servers and test authentication shape, issue
recovery, ADF, independent retries, redaction and receipt deduplication.

Docker builds and non-root offline version/validation smoke checks are part of
local verification and CI. GitHub-hosted execution, arm64 runtime, AWS permissions,
and real service delivery have not been verified. No remote workflows were enabled
and no releases, messages, or issues were published during implementation.
