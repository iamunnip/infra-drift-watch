# Configure GitHub, Slack, and optional Jira

Set destination settings under defaults, or override them per target. The shipped
example disables all delivery. Configure runtime environment variables before enabling.

```yaml
defaults:
  github:
    enabled: true
    repository: your-org/drift-issues
    token_env: GITHUB_TOKEN
    labels: [drift]
  slack:
    enabled: true
    webhook_env: SLACK_WEBHOOK_URL
    repeat: 24h
  jira:
    enabled: false
    url: https://your-team.atlassian.net
    project: INFRA
    issue_type: Task
    email_env: JIRA_EMAIL
    token_env: JIRA_TOKEN
```

GitHub needs issue read/write access to the destination. The workflow token only
has access within its repository; use an appropriately scoped secret token for
another repository. Ensure configured labels exist. The client searches repository
issues (including closed) for its stable marker, then creates or updates one issue
per target. It reopens a matching closed issue for drift and keeps resolution
updates open. Issue bodies are scanner-owned; keep human notes in comments.

Slack uses an incoming webhook; the destination is chosen when provisioning the
webhook. Messages use plain text and contain available tracker links. Store the
URL in a secret. A failed destination does not erase another destination's receipt.

Jira Cloud uses email/API-token Basic authentication, enhanced JQL search,
create/update issue APIs, a stable label, and ADF plain-text descriptions. Configure
customfield_* values in jira.fields if required by the project's create screen.
Jira workflow status transitions are manual in this release; resolved drift changes
the summary/description, not Jira status. Eventual search indexing can cause duplicates.

Unchanged issues are not rewritten; unchanged Slack reminders follow repeat.
Changed findings and resolved drift send new events. Failed, incomplete, and
suppressed scans never resolve previous findings. Notification acknowledgement is
not evidence that infrastructure was fixed.

Use `--dry-run` for live inspection without any destination requests or receipts.
Normal tests use local TLS servers; they never send messages to real destinations.
