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

### Requirement: AWS selectors without credential storage
AWS profile and region SHALL be configurable globally, per target and through explicit CLI overrides. Omitted selectors SHALL preserve runtime AWS selection; explicit empty strings SHALL clear inherited selector environment values. Credentials SHALL remain outside scanner configuration.

#### Scenario: Select the default profile
- **WHEN** a target resolves aws.profile default and aws.region us-east-1
- **THEN** those selectors reach all native tool commands without changing the host environment

#### Scenario: OIDC target overrides local profile
- **WHEN** the CLI explicitly clears the AWS profile
- **THEN** inherited named-profile selectors are cleared in child processes while temporary credential variables remain available
