# 0004: Detailed findings, bounded history and optional closure

Status: Accepted
Date: 2026-11-04
Supersedes: ADR 0002's coarse Kubernetes findings and manual-only tracker status decisions.

## Context

Operators need resource/field detail, recent scan history, and configurable tracker
resolution without interpreting missing coverage as success.

## Decision

Request full-prefix unified Kubernetes diffs, derive resource identities and field
paths when parseable, and publish no values. Hash ephemeral changed-value pairs;
ignore unchanged context in fingerprints. Preserve resource/set-level fallback
and manifest-set suppressions when detail is unavailable.

Keep bounded per-target history in state format 2, sharing the writer lock and
atomic replacement with independent receipts. Preserve format-1 receipts during
migration and keep read-only history inspection independent of inventory/secrets.
Record failed/incomplete outcomes; dry-run records nothing.

Keep closure off by default. GitHub closure and Jira resolve/reopen transitions
are opt-in per destination/target. Jira requires both transition IDs, validates
status categories and never acknowledges a failed transition. Incomplete scans
cannot resolve. Policy changes can reconcile known resolved issues.

## Consequences

State upgrades are not readable by old binaries. Cache eviction can lose local
history. Jira workflows still require deployment-specific setup and verification.
Raw plan/diff values remain ephemeral. No cloud resources are remediated and no
UI, persistent daemon, or distributed storage is introduced.

Workflow status is read from the current Jira issue before transition decisions;
search results are used only to locate the issue, avoiding stale-index recurrence decisions.
