# 0005: Drop arc42 for C4 views and Diataxis

Status: Accepted
Date: 2026-11-06
Supersedes: the `docs/arc42/` location for architecture documentation.

## Context

The architecture narrative lived in 12 short arc42 sections with C4 diagrams embedded in
them. For one CLI binary, most sections were a few lines, one only linked to the ADR index,
and the same facts were also kept in the README, Diataxis pages and ADRs.

## Decision

Remove `docs/arc42/`. Keep each C4 view in `docs/c4/` with its diagram and explanation on
one page: system context, containers and components, runtime and deployment. Move goals,
constraints, approach, principles, quality requirements and risks to the Diataxis
explanation `docs/diataxis/explanation/design-and-risks.md`, and the glossary to
`docs/diataxis/reference/glossary.md`. OpenSpec, OKF and ADRs are unchanged.

## Consequences

There is one fewer documentation standard to maintain and each fact has one home. Readers
lose arc42's fixed section checklist, so reviewers must check that design and risk notes
stay current. Older links to `docs/arc42/` no longer resolve; archived OpenSpec changes
keep their historical references.
