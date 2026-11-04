# Verification and compatibility

| Component | Pinned test/runtime baseline | Verification scope |
| --- | --- | --- |
| Go | 1.27.3 | Unit, race, vet and binary/container builds |
| OpenTofu | 1.13.5 | Real disposable local workspace; plan fixtures |
| Terraform | 1.17.1 | Real disposable local workspace; plan fixtures |
| kubectl | 1.37.2 | Command boundary tests and disposable cluster test |
| kind | 0.33.0 | Optional Kubernetes test prerequisite |
| Kubernetes | kind node 1.37.0 | Disposable ConfigMap clean/change/resolution |
| GitHub, Slack, Jira | Documented REST endpoints | Mocked TLS-server tests; no live account verification |
| AWS | AWS provider ~> 6.71.0 example | Configuration example only; no live AWS scan verified |

Pinned baselines are compatibility fixtures, not a promise that older binaries
receive upstream security support. Review upstream updates before production use.
Only Linux amd64 execution is verified locally; arm64 release artifacts are built
but need deployment testing. Plan format major 1 is accepted; malformed or errored
plans fail, and explicitly incomplete/deferred plans produce incomplete results. Version combinations outside this table are unverified.

`make check` runs local checks after dependencies are present. `make tools` installs
runtime executables explicitly. `make test-integration` uses installed tools from PATH, including Brew. Kubernetes
tests need kind/kubectl/Docker and create/delete their own named cluster using a
temporary kubeconfig. Tests apply only disposable fixtures, never your AWS resources.

For Mermaid run npm ci --prefix scripts/mermaid --ignore-scripts then make
mermaid-check. `make dev-tools` installs pinned govulncheck/actionlint; `make vuln`
and `make workflow-check` execute the tools on PATH. Explicit installer targets
are not dependencies of local checks. CI installs pinned tools only on ephemeral runners. Live cloud,
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

## Updated tool and workflow baseline

The current dependency baseline is Go 1.27.3, YAML v3.0.6, local test provider
2.9.1, Mermaid 12.3.0, jsdom 30.1.3, OpenSpec 1.16.0, golangci-lint 2.15.0,
actionlint 1.7.12 and govulncheck 1.9.0. Node 26.13.0 runs documentation checks.
Alpine 3.24.3 is the runtime base. The OpenTofu container uses the official 1.13.5-minimal
image; the optional installer builds the checksum-verified source tag because
that tag has no published binary ZIPs. Updating its version requires updating
its source checksum too.

PR/manual CI has distinct Go, documentation, and integration jobs. Documentation
checks explicitly validate docs, knowledge/OKF, Mermaid and strict OpenSpec.
No workflow has a push trigger. The scan schedule remains separately opt-in.

Updated-baseline verification passed on the Brew tools: OpenTofu 1.13.5, Terraform
1.17.1, kubectl 1.37.2 and kind 0.33.0 with Kubernetes 1.37.0. The local provider
2.9.1 fixtures passed. Go/race/module checks, actionlint 1.7.12, strict OpenSpec,
Mermaid, separate docs/knowledge checks, and golangci-lint 2.15.0 in a read-only
container passed; govulncheck reported no vulnerabilities. npm dependency refresh
reported no vulnerabilities and no newer compatible Mermaid/jsdom versions.
No host tool installers were run for this update. Domain knowledge did not change.

License texts and notice paths are documented in [third-party notices](../../../THIRD_PARTY_NOTICES.md).
Runtime images and manual release archives include the named component licenses.

Feature verification adds resource/field identity and fingerprint checks on real
Kubernetes diffs; history retention, version-1 migration, read-only CLI, failure
records and dry-run preservation; and mocked GitHub/Jira closure, recurrence,
policy-change reconciliation and unavailable-transition retries. State format 2
preserves old receipts; public scan/history JSON schemas remain version 1.

The completed feature checks also cover unavailable Jira transitions, fresh issue
status instead of stale search status, closure-policy changes, manifest-set
suppression compatibility, and oversized process output rejection. All local Go
checks pass, the container linter reports zero issues, and the updated image can
read history with no network, inventory, credentials or writable mounts. Six
Mermaid diagrams and seven OpenSpec capabilities validate. Live AWS/destination
verification still needs deployment configuration; no live services were exercised.
