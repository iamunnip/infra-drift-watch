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
suppressed. An issue's closure is a human workflow action, not proof of clean infrastructure.

See [delivery setup](../../docs/diataxis/how-to/configure-delivery.md).
