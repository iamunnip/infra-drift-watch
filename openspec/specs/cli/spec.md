# CLI entry point

## Purpose

CLI entry point. Define testable behavior, safe failure handling, and documented limits.

## Requirements

### Requirement: Discover commands
The CLI SHALL expose help, version, validate, scan and history commands.

#### Scenario: Request help
- **WHEN** the user invokes --help
- **THEN** usage lists validate and scan

### Requirement: Validate offline
The validate command SHALL validate settings and selection without credentials or external calls.

#### Scenario: Validate an inventory
- **WHEN** a valid global file and inventory are supplied
- **THEN** validation succeeds without writing receipts

### Requirement: Scan exit status
The scan command SHALL return 0 for clean, 2 for drift, and 1 for failure or incomplete coverage; failures SHALL take precedence.

#### Scenario: Mixed target outcomes
- **WHEN** one target has drift and another fails
- **THEN** the report preserves both results and exit code is 1

### Requirement: Dry run
The scan command SHALL skip all delivery requests and receipt and history writes when --dry-run is supplied.

#### Scenario: Inspect without delivery
- **WHEN** a user runs scan --dry-run
- **THEN** infrastructure is inspected and no destination is called
