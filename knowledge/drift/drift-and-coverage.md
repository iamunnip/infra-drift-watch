---
type: Concept
title: Drift and scan coverage
description: Distinguishing observed drift, pending changes, and incomplete inspection.
tags: [drift, coverage, infrastructure]
---

# Drift and scan coverage

For state-based infrastructure tools, observed drift describes remote changes
relative to recorded state. Pending configuration changes describe intended
changes that have not been deployed. A planned action alone is not proof of drift.

For Kubernetes, interpretation requires an explicit manifest scope and rules for
server defaults and controller-owned fields. An omitted resource is not automatically
an unwanted resource.

A clean scan means no drift was found within successfully inspected scope.
Failed or incomplete inspection cannot establish a clean result or resolve an
existing finding. Scan coverage must remain visible alongside findings.
