# Infrastructure inspection and coverage

## Purpose

Infrastructure inspection and coverage. Define testable behavior, safe failure handling, and documented limits.

## Requirements

### Requirement: Observed drift
State adapters SHALL use refresh-only plans and resource_drift JSON, without applying changes or saving refreshed state.

#### Scenario: Pending configuration
- **WHEN** configuration changed but remote objects did not
- **THEN** the scan reports clean for observed drift

### Requirement: Private artifacts
Temporary plans SHALL be created in private directories and removed after normal execution, including failure.

#### Scenario: Plan interpretation fails
- **WHEN** show returns invalid JSON
- **THEN** the scan fails and its temporary directory is removed

### Requirement: Kubernetes scope
The Kubernetes adapter SHALL compare explicit local manifests using server-side diff and an explicit context and field manager.

#### Scenario: Owned field changes
- **WHEN** a supplied manifest differs in an owned field
- **THEN** the result reports manifest-set drift without raw values

### Requirement: Incomplete coverage
Failed, incomplete or suppressed scans SHALL NOT establish clean status or resolve previous findings.

#### Scenario: Active suppression
- **WHEN** a matching unexpired suppression hides a finding
- **THEN** the result is incomplete and resolution delivery is skipped
