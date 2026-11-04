# Bounded scan history

## Purpose

Keep useful past scan outcomes without retaining raw infrastructure values or requiring inventory and credentials to inspect them.

## Requirements

### Requirement: Bounded history records
Normal scans SHALL persist the latest configured number of records per target, including successful, failed and incomplete outcomes, inspection duration and destination errors. Records SHALL omit raw plan/diff values.

#### Scenario: Retention limit reached
- **WHEN** a target records another scan beyond its configured retention count
- **THEN** only its newest configured number of records remain and other targets are preserved

### Requirement: Preserve existing receipts
State format 2 SHALL preserve format-1 delivery receipts during migration and save the updated state atomically using the existing single-writer lock.

#### Scenario: Upgrade an existing state directory
- **WHEN** an application opens format-1 receipts and saves a new record
- **THEN** existing delivery acknowledgements remain and the saved format is version 2

### Requirement: Read-only history command
The history command SHALL read an existing snapshot without credentials, inventory, remote calls or writes, support exact target selection and limits, and return newest records first in terminal or JSON format.

#### Scenario: Missing state
- **WHEN** history is requested for an absent state directory
- **THEN** the command fails without creating that directory

### Requirement: Dry-run and persistence errors
Dry-run SHALL NOT append history or change receipts. A history persistence failure during a normal scan SHALL fail the command and appear in the current report.

#### Scenario: Inspect without recording
- **WHEN** scan --dry-run is invoked with existing history
- **THEN** the history and receipt file remains byte-for-byte unchanged
