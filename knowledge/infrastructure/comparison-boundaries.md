---
type: Concept
title: State tools and manifest comparison
description: What refresh-only plans and server-side manifest comparisons establish.
tags: [terraform, opentofu, kubernetes, coverage]
---

# Comparison boundaries

Terraform/OpenTofu state is the recorded baseline for observed infrastructure drift.
A refresh-only plan inspects changes to managed objects without proposing deployment
of pending configuration. Empty state does not establish that desired infrastructure
exists. Refresh-only results do not discover all unmanaged account resources.

Kubernetes comparison uses supplied manifests and server-side dry-run semantics.
Defaults, admission, and field ownership affect the result. Ownership conflicts
are failures requiring review, never clean results. Resources omitted from supplied
manifests are outside coverage. This release reports a manifest-set finding rather
than attributing Kubernetes differences to individual fields or resources.

See [architecture constraints](../../docs/arc42/02-constraints.md).
