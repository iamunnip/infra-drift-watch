# 0002: Delegate comparison and preserve independent delivery

Status: Accepted
Date: 2026-11-04

## Context

Three infrastructure systems have different comparison semantics. Delivery can
partially fail, and notifications need repeat suppression across scheduled runs.

## Decision

Use installed tools with explicit arguments: refresh-only plan plus JSON inspection
for state tools, and server-side kubectl diff for selected manifests. Never apply.
Report values nowhere; keep only resource/action summaries and internal fingerprints.
Use per-target issue grouping, remote markers, independent receipts, a locked local
store, and sequential execution. Do not retry ambiguous mutation failures. Retry
explicit rate-limit rejections within a bounded time budget. Keep issues open on
resolution; no automatic closure is implemented.

## Consequences

Workspaces must be initialized outside the scanner. Kubernetes findings are coarse.
Delivery caches are best effort and cannot guarantee exactly-once notifications.
Jira status transitions remain manual. Suppressed targets cannot automatically
resolve. These limits favor honest coverage over silently declaring success.
