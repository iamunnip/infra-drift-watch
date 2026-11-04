# Project verification

## Purpose

Provide explicit, reproducible validation for pull requests and manual runs while preserving contributor-managed tools.

## Requirements

### Requirement: Pull request and manual validation
CI SHALL run on pull requests targeting develop or main and workflow_dispatch. Workflows SHALL NOT use push triggers.

#### Scenario: A merge request is raised
- **WHEN** a pull request is opened or updated against develop or main
- **THEN** Go, documentation, and integration validation jobs run

### Requirement: Complete validation coverage
Validation SHALL check Go formatting, lint, vet, module integrity, tests and build; documentation links; OKF knowledge metadata; Mermaid syntax; and strict OpenSpec conformance.

#### Scenario: Invalid documentation or specification
- **WHEN** a required docs, knowledge, diagram, or OpenSpec check fails
- **THEN** the corresponding validation job fails

### Requirement: Contributor-managed tools
Local validation SHALL resolve tools from PATH without installing host binaries or prepending ignored local tool copies.

#### Scenario: Brew tools are installed
- **WHEN** a contributor invokes verification or integration targets
- **THEN** the installed tools are used without executing installer targets

### Requirement: Manual draft release
Release publishing SHALL require workflow_dispatch with an existing version tag whose commit belongs to main.

#### Scenario: A tag is pushed
- **WHEN** a release tag is pushed without a manual dispatch
- **THEN** no release workflow starts
