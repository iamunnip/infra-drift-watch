# 0003: Current tool baselines and PR/manual validation

Status: Accepted
Date: 2026-11-04
Supersedes: the Go minimum and foundation-only tooling choices in ADR 0001.

## Context

The project now has scanning and delivery, and contributors install tools with
Homebrew. Local ignored copies must not override those tools. Validation is
required on merge requests and manual runs rather than every push.

## Decision

Require Go 1.27.3 and pin currently verified stable tools/dependencies. Resolve
local verification tools from PATH. Install pinned tools only explicitly or on
ephemeral CI runners. Run Go lint, formatting, vet, module integrity, race tests,
builds, documentation links, OKF metadata, Mermaid parsing and strict OpenSpec
validation on PR/manual CI. Run native integrations/container checks separately.
Publish draft release artifacts only by manual dispatch against an existing tag
on main. No workflow uses a push trigger. Use weekly Dependabot updates into develop.

## Consequences

Local checks install no tools. Containers provide a linter without a host binary.
Workflows are reviewable and manually runnable; the opt-in scan schedule remains
independent. Version bumps require repeatable verification and updated documentation.
