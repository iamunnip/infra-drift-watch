# Working in infra-drift-watch

Read README.md, docs/index.md, knowledge/index.md, and relevant OpenSpec specs
and ADRs before changing behavior. The project has an initial scanning and delivery implementation; live deployments remain unverified.

## Workflow

Develop through `develop`; feature and fix PRs target `develop`. Releases move
through a PR to `main`. Do not publish, push, or contact external destinations
merely to verify a change.

Use TDD: write a meaningful failing test, confirm the reason for failure,
implement the smallest solution, then refactor. Bug fixes require regression
tests. Keep unit tests next to code and fixtures in package-local `testdata/`.
Use mocked delivery and disposable infrastructure for integration tests.

## Go structure

Keep cmd/infra-drift-watch/main.go small. Implementation belongs in internal/.
Create responsibility packages as features arrive, not speculative empty layers.
Use clear names, explicit errors, context cancellation for external work, and
interfaces at consumer boundaries. Avoid generic utils packages and unnecessary
dependencies. No public pkg/ API is planned.

## Safety and configuration

Never apply or destroy infrastructure or persist refreshed infrastructure state.
Never print secrets or commit credentials, state files, or plan artifacts.
Incomplete scans cannot report clean or resolve findings. Destination failures
must remain independent; stable target IDs drive issue identity.
Global YAML defaults and target overrides must have documented merge semantics.
Validate unknown settings and invalid combinations. Operational configuration
must not disable core safety guarantees.

## Documentation

- openspec/: behavior and acceptance scenarios; keep specs aligned with changes.
- knowledge/: OKF domain knowledge, typed frontmatter, relative links, indexes.
  Update knowledge/log.md when changing domain facts or policies.
- docs/diataxis/: tutorials, how-to, reference, explanation.
- docs/adr/: lasting decisions; preserve accepted history and supersede explicitly.
- docs/arc42/: architecture narrative and C4 views using Mermaid in Markdown.
- README.md: comprehensive entry point, with implemented and planned features
  clearly distinguished.

Update affected documentation with behavior changes. Review knowledge impacts;
do not invent domain changes solely to modify a file.

## Verification

Run `make check` (format check, vet, module integrity, race tests, build, local docs/OKF checks). Use `make test` or focused
`go test` commands during TDD. `make fmt` changes formatting. `make docker-build`
requires a Docker daemon and network access for its builder image.
Runtime/developer tools require explicit opt-in `make tools` / `make dev-tools`.
Run Mermaid parsing with `make mermaid-check` after the documented npm install.
Integration tests may provision disposable local fixtures; never real AWS resources.
Only list checks as passed when run successfully. Missing prerequisites and
unverified remote behavior must be reported. Do not install tools silently.

Validation triggers are pull_request and workflow_dispatch, never push. Releases
are manual dispatch against an existing version tag on main. Local checks resolve
tools from PATH; do not prepend ignored .tools/bin copies. Do not run installer
targets on the user's host. Use the container linter when golangci-lint is absent.
Also verify make openspec-check, mermaid-check, workflow-check and vuln when their
prerequisites are available. CI checks all of these with pinned current tools.

Keep THIRD_PARTY_NOTICES.md, licenses/ provenance, container copies and release
notices aligned when changing dependencies or bundled tools. Distinguish project
MIT, OpenTofu MPL-2.0 and the pinned Terraform BSL-1.1; providers have separate terms.
