# Step 0: project information collection

Collected on 2026-11-04. This records requirements and observations; it does not
claim that planned features are implemented.

## Confirmed requirements

- Open source Go CLI, licensed under MIT; no web UI for the initial scope.
- Development through `develop`, with feature/fix PRs into `develop` and releases
  promoted to `main`.
- Test driven development: demonstrate a meaningful failing test, implement,
  then refactor. Bug fixes start with a reproducing regression test.
- Detect and report infrastructure drift for OpenTofu, Terraform, and Kubernetes
  manifests, delivered incrementally with OpenTofu first.
- AWS is the confirmed provider for the first working OpenTofu example.
- Inventory with stable target identities, global YAML configuration, and
  individual target overrides. Operational settings must be configurable and
  validated; core safety guarantees remain enforced.
- Terminal and JSON reports, Slack notifications, GitHub issue creation and
  updates, and optional Jira integration.
- Group issues by target by default. Notify and mark resolved drift; automatic
  issue closure is disabled by default.
- Initial scheduled execution through GitHub Actions.
- Use GitHub-hosted runners for the initial AWS scans.
- Persistent local delivery state, initially with a single writer. Workflow
  runs need an explicit strategy to preserve state between executions.
- Dockerfile and `.dockerignore`, a non-root runtime, and documented runtime
  dependencies for infrastructure adapters.
- Beginner-friendly code and documentation, root `AGENTS.md`, and a Makefile
  exposing documented build and verification commands.
- A comprehensive root README covering status, setup, examples, architecture,
  configuration, limitations, development, and documentation navigation.

## Documentation responsibilities

| Location | Responsibility |
| --- | --- |
| `openspec/` | Behavior specifications, acceptance scenarios, and changes |
| `knowledge/` | OKF domain concepts and policies, indexed with relative links and an update log |
| `docs/adr/` | Significant decisions and consequences |
| `docs/diataxis/` | Tutorials, how-to guides, reference, and explanation |
| `docs/arc42/` | Architecture narrative, including C4 views |

Use Mermaid diagrams in Markdown. C4 does not need a separate directory.
The local `infra-vuln-watch/knowledge/` bundle is the reference for OKF layout;
its root uses `okf_version: "0.2"`, topic indexes, typed content, and `log.md`.
Do not copy vulnerability-specific facts into this project.

## Repository and environment observations

- Repository: `github.com/iamunnip/infra-drift-watch`.
- Current branch: `develop`; `main` also exists locally.
- Starting commit: `a69214a` (`Initial commit`).
- The only tracked file before this intake was `README.md`.
- No repository `AGENTS.md`, Go module, tests, or implementation existed.
- Local host: Linux, amd64.
- Installed Go reports `go1.27.3-X:nodwarf5`; this is an environment observation,
  not a selected release or CI toolchain requirement.
- Installed tools: Docker 29.9.2, GNU Make 4.4.1, GitHub CLI 2.103.3,
  OpenSpec 1.16.0.
- Terraform, OpenTofu (`tofu`), and `kubectl` were not found on `PATH`.
- Docker daemon availability and external service access have not been tested.
- No credentials were inspected and no infrastructure or delivery APIs called.

## Information still needed before live execution

- Representative AWS resource and target for the first OpenTofu example.
- Confirm that the selected AWS target and state backend are accessible from
  GitHub-hosted runners.
- Destination repository for drift issues; it need not be the source repository.
- Scan schedule and time zone.
- Slack destination and authentication method.
- Optional Jira project, issue type, field mapping, and authentication method.
- Storage and retention for delivery state between workflow runs.
- Runtime executable versions, supported OS/architectures, and container contents.

Credentials should be configured as runtime secrets, never entered in this file.
These deployment details do not block foundation work.

## Next step

Step 1 establishes repository guidance, MIT licensing, initial specifications,
documentation indexes, and a small Go scaffold with meaningful tests. Toolchain
and dependency versions must be selected and verified during that work.
No live scans, issue creation, or notifications are needed for foundation checks.
