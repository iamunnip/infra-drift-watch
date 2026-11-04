# Strict inventory and layered configuration

## Purpose

Strict inventory and layered configuration. Define testable behavior, safe failure handling, and documented limits.

## Requirements

### Requirement: Strict documents
Configuration SHALL reject unknown keys, multiple YAML documents, duplicate IDs and invalid settings.

#### Scenario: Unknown setting
- **WHEN** an unknown YAML field is supplied
- **THEN** validation fails before scanning

### Requirement: Target overrides
Configuration SHALL merge global settings with target overrides; explicit false and zero SHALL override inherited values, and lists/maps SHALL replace.

#### Scenario: Disable inherited GitHub delivery
- **WHEN** a target sets github.enabled false
- **THEN** GitHub is disabled only for that target

### Requirement: Stable identity
Each inventory target SHALL have a stable unique ID independent of owner, environment and tool settings.

#### Scenario: Duplicate target IDs
- **WHEN** two targets share an ID
- **THEN** validation rejects the inventory

### Requirement: History and closure settings
History SHALL default to enabled with 100 records per target, with a configurable limit from 1 to 10000. GitHub and Jira auto_close SHALL default to false. Jira auto_close SHALL require resolve and reopen transition IDs.

#### Scenario: Target opts out
- **WHEN** a target explicitly disables inherited history or closure settings
- **THEN** false overrides the global setting for that target
