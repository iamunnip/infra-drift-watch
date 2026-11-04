---
type: Concept
title: Infrastructure engines and provider licenses
description: Why engine, provider, and scanner licenses must be distinguished by version.
tags: [licensing, opentofu, terraform, providers]
---

# Infrastructure engines and provider licenses

A scanner's license does not relicense the engines or providers it invokes.
OpenTofu uses MPL-2.0. The pinned Terraform 1.17.1 uses BSL-1.1 with an additional
use grant and future MPL change provisions; historical Terraform 1.5.7 used MPL-2.0.
Engine compatibility does not imply license equivalence.

Providers are independent components. A provider may use a different license from
its engine, so review the selected provider version rather than inheriting the
engine's label. The pinned AWS/local provider sources are MPL-2.0, while other
provider licenses may differ. Redistributing an image also includes independent
OS utilities and libraries, not just the scanner.

See [version-specific third-party notices](../../THIRD_PARTY_NOTICES.md) and the
[licensing explanation](../../docs/diataxis/explanation/licensing.md) for source links
and the exact pinned license texts.
