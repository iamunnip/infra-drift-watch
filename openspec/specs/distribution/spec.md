# Distribution notices

## Purpose

Preserve the distinction between the project license and version-specific third-party terms in source, container, and release artifacts.

## Requirements

### Requirement: Named third-party notices
The repository SHALL use THIRD_PARTY_NOTICES.md and preserve the named component license texts with version and source provenance.

#### Scenario: Inspect project licenses
- **WHEN** a contributor reads the notices
- **THEN** project MIT, OpenTofu MPL-2.0, the pinned Terraform BSL-1.1 and separate provider/library terms are distinguished

### Requirement: Include notices in artifacts
The runtime image SHALL include the project license, third-party notices and named license texts. Manual releases SHALL attach a notices archive with its checksum.

#### Scenario: Inspect the runtime image
- **WHEN** the Docker image is built
- **THEN** /LICENSE, /THIRD_PARTY_NOTICES.md and /usr/share/licenses/infra-drift-watch/ are present
