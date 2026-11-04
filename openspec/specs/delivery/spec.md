# Independent delivery and persistence

## Purpose

Independent delivery and persistence. Define testable behavior, safe failure handling, and documented limits.

## Requirements

### Requirement: GitHub issue grouping
GitHub delivery SHALL reconcile one marked issue per target and leave resolution updates open.

#### Scenario: Recover a previous issue
- **WHEN** a matching issue marker exists without a local acknowledgement
- **THEN** the client updates that issue instead of creating a new one

### Requirement: Independent acknowledgement
Each enabled destination SHALL acknowledge separately and preserve successful receipts when another destination fails.

#### Scenario: Slack failure
- **WHEN** GitHub succeeds and Slack fails
- **THEN** the next scan retries Slack without rewriting unchanged GitHub delivery

### Requirement: Optional Jira
Jira Cloud delivery SHALL use stable labels, enhanced JQL lookup and plain-text ADF descriptions.

#### Scenario: Jira enabled
- **WHEN** drift is found and no matching Jira issue exists
- **THEN** the client creates an issue with configured project and type

### Requirement: Bounded retry
Delivery SHALL retry explicit rate-limit rejections within configured budgets and SHALL NOT retry ambiguous mutation failures.

#### Scenario: Server error on create
- **WHEN** an issue creation returns HTTP 500
- **THEN** the error is reported and creation is not immediately repeated

### Requirement: Single writer
Persistent receipts SHALL use an exclusive lock and atomic replacement.

#### Scenario: Concurrent invocation
- **WHEN** a second process opens the same locked state directory
- **THEN** the second process fails before scanning
