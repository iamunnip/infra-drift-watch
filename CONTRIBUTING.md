# Contributing

Start with [README](README.md) and [agent guidance](AGENTS.md).
Create a feature/fix branch from `develop`, and target `develop` in your PR.

Use Go 1.27.3 for the foundation toolchain; the module minimum is Go 1.26.
Run `make test` while iterating and `make check` before submitting.
Follow red, green, refactor: confirm a failing behavioral test before implementation.
Bug fixes begin with a reproducing test. Explain the behavior change, verification,
and documentation impact in your PR. No infrastructure credentials are needed
for foundation tests.

[Architecture](docs/arc42/index.md) and [specifications](openspec/specs/cli/spec.md)
provide the current boundaries. Scanning and delivery clients are implemented; live account behavior remains unverified.
See [verification](docs/diataxis/reference/verification.md) for integration prerequisites.
