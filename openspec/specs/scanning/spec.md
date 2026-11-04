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
- **THEN** the result reports resource identity and changed field paths when parseable, without raw values; unfamiliar output retains a resource/set-level fallback

### Requirement: Incomplete coverage
Failed, incomplete or suppressed scans SHALL NOT establish clean status or resolve previous findings.

#### Scenario: Active suppression
- **WHEN** a matching unexpired suppression hides a finding
- **THEN** the result is incomplete and resolution delivery is skipped

### Requirement: Stable detailed findings
The Kubernetes adapter SHALL derive slash-separated resource identities and value-free field paths from supported native diff headers/hunks. Unchanged metadata and temporary filenames SHALL NOT change the fingerprint of unchanged findings.

#### Scenario: Repeat a changed resource scan
- **WHEN** the resource drift is unchanged but native temporary filenames or unchanged resource metadata differ
- **THEN** the finding fingerprint stays stable

### Requirement: Preserve whole-set suppressions
The manifest-set suppression SHALL continue to cover all Kubernetes resource findings.

#### Scenario: Existing set suppression
- **WHEN** an unexpired manifest-set suppression is configured for a detailed Kubernetes comparison
- **THEN** matching findings are suppressed and target coverage remains incomplete
