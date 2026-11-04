# Contributing

Start with [README](README.md) and [agent guidance](AGENTS.md).
Create a feature/fix branch from `develop`, and target `develop` in your PR.

Use Go 1.27.3 for the foundation toolchain; the module also requires Go 1.27.3.
Run `make test` while iterating and `make check` before submitting.
Follow red, green, refactor: confirm a failing behavioral test before implementation.
Bug fixes begin with a reproducing test. Explain the behavior change, verification,
and documentation impact in your PR. No infrastructure credentials are needed
for foundation tests.

[Architecture](docs/arc42/index.md) and [specifications](openspec/specs/cli/spec.md)
provide the current boundaries. Scanning and delivery clients are implemented; live account behavior remains unverified.
See [verification](docs/diataxis/reference/verification.md) for integration prerequisites.

Validation runs on PR creation/updates and workflow_dispatch, with no push trigger.
Run make lint (installed golangci-lint) or make lint-docker, make mermaid-check,
make openspec-check and make workflow-check in addition to make check.
Local checks use PATH, including Brew installations; no tools are installed by
check targets. make tools/dev-tools/kind-tools are explicit optional installers,
used only on ephemeral CI runners unless a contributor deliberately invokes them.

When updating tool/dependency versions, update [third-party notices](THIRD_PARTY_NOTICES.md)
and preserved license texts/provenance. Project MIT does not replace upstream terms.
