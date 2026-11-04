---
type: Process
title: Delivery and resolution
description: Independent acknowledgements, stable identity, exceptions, and resolution.
tags: [delivery, identity, resolution, suppression]
---

# Delivery and resolution

A target ID identifies a deployment scope and must stay stable when its owner,
environment label, or configuration changes. Reusing an ID for another deployment
can attach findings to the wrong issue. Group issues by target, not by mutable title.

Each destination acknowledges independently. Success in one tracker does not prove
another tracker or notification channel received an event. Remote acceptance and
local receipt persistence can diverge on timeout or interruption. Exactly-once
delivery is not guaranteed; markers help recover issues, while notifications can repeat.

Only a successful complete scan can establish resolution in selected scope. Failed,
incomplete, or suppressed scans cannot establish that earlier drift disappeared.
An expired suppression no longer hides findings. Suppressions need a reason and
an end date. This release conservatively pauses target delivery when findings are
suppressed. Issue closure is manual by default and can be explicitly enabled for verified
resolution. Closure alone is not proof of clean infrastructure. Jira transitions
need project-specific IDs and destination-category checks; failed transitions
remain pending. Recurring drift reopens trackers according to configured policy.

See [delivery setup](../../docs/diataxis/how-to/configure-delivery.md).


# Scan history

A historical record describes what was observed within a target's selected scope
at one time. An incomplete result cannot establish resolution, and a notification
receipt cannot establish that infrastructure was fixed. Keep history bounded and
value-free; retain it outside ephemeral runner caches when long-term history is
required. Disabling collection does not erase past observations.
