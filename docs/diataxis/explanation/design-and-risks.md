# Design and risks

Why infra-drift-watch is shaped the way it is, and what it accepts. Diagrams are in the
[C4 views](../../c4/index.md); decisions are in the [ADRs](../../adr/index.md).

## Goals

Detect observed infrastructure changes with explicit scope, safe output, and dependable delivery accounting. Prioritize honest coverage, beginner-friendly code, and reproducible TDD. No initial UI or automatic remediation.

## Constraints

Go implementation; Linux runtime; development through develop. Existing initialized state workspaces are required. AWS hosted-runner connectivity and least-privilege credentials must be configured outside the scanner. Kubectl dry-run requires RBAC and may invoke admission webhooks. Only trusted IaC and executable plugins are in scope.

Go 1.27.3 is required. PR/manual CI performs full validation; no workflow runs on push. Local tool resolution uses PATH.

## Approach

Delegate infrastructure interpretation to native tools. Parse refresh-only resource_drift rather than treating pending actions as drift. Bound process output and execution time. Omit raw values; persist only delivery fingerprints and available issue links. Run targets and delivery serially initially.

## Principles

Stable target IDs define issue identity. Each destination has a separate fingerprint and acknowledgement. Secret values are omitted; raw process stderr is discarded. SHA-256 fingerprints detect changes without retaining raw diff values. TLS and redirect restrictions protect credentials. Context cancellation bounds external work. Configuration is strict YAML with documented override semantics.

State format 2 preserves migrated receipts and bounded per-target history. History reads use atomic file snapshots without writes. Optional tracker closure is explicit; Jira transition destinations are validated. Kubernetes reports field paths and hashes ephemeral value differences without persisting values.

AWS profile and region selectors resolve globally, per target, and through CLI overrides. They are applied only to native-tool subprocesses; SDK credential precedence and explicit provider/backend settings remain in force. No credential files are read or rewritten by the application.

## Quality requirements

Tests must establish clean, drift, failure and incomplete behavior. Prove no delivery/state writes during dry-run, independent retries, pending-change separation, deterministic ordering, and no raw secret output. Local fixture integration covers native tools; live accounts require deployment validation. TDD is required for behavior and regression fixes.

## Risks and technical debt

Remote acceptance may precede local acknowledgement. Receipt cache eviction, Jira indexing delay, concurrent external writers, or manual marker removal can cause duplicates. Kubernetes field detail depends on parseable native diff hunks; unknown formats retain coarser findings. Suppressions pause all delivery for that target. Supported provider behaviors and cloud permissions require deployment checks. A crash can leave temporary plan artifacts. History is bounded local state and can be lost with cache eviction. No distributed locking, persistent daemon, or per-field policy engine exists. Optional closure requires live permission/workflow verification.
