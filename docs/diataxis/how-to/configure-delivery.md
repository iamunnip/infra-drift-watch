# Configure GitHub, Slack, and optional Jira

Set destination settings under defaults, or override them per target. The shipped
example disables all delivery. Configure runtime environment variables before enabling.

```yaml
defaults:
  github:
    enabled: true
    auto_close: false
    repository: your-org/drift-issues
    token_env: GITHUB_TOKEN
    labels: [drift]
  slack:
    enabled: true
    webhook_env: SLACK_WEBHOOK_URL
    repeat: 24h
  jira:
    enabled: false
    auto_close: false
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
updates open by default. github.auto_close true closes it on a complete clean
result, then recurrence reopens it. Issue bodies are scanner-owned; keep human notes in comments.

Slack uses an incoming webhook; the destination is chosen when provisioning the
webhook. Messages use plain text and contain available tracker links. Store the
URL in a secret. A failed destination does not erase another destination's receipt.

Jira Cloud uses email/API-token Basic authentication, enhanced JQL search,
create/update issue APIs, a stable label, and ADF plain-text descriptions. Configure
customfield_* values in jira.fields if required by the project's create screen.
Jira workflow status transitions are manual by default. With jira.auto_close true,
configure project-specific resolve_transition and reopen_transition IDs. The client
checks transition availability and target status category before posting. An
unavailable transition is a delivery failure and remains pending. A previously
completed issue is reopened when drift recurs. Eventual search indexing can cause duplicates.

Unchanged issues are not rewritten; unchanged Slack reminders follow repeat.
Changed findings and resolved drift send new events. Failed, incomplete, and
suppressed scans never resolve previous findings. Notification acknowledgement is
not evidence that infrastructure was fixed.

Use `--dry-run` for live inspection without any destination requests or receipts.
Normal tests use local TLS servers; they never send messages to real destinations.

## Enable optional closure

Use the transition IDs available in your Jira workflow; these example IDs are
illustrative, not defaults:

```yaml
defaults:
  github:
    auto_close: true
  jira:
    auto_close: true
    resolve_transition: "31"
    reopen_transition: "21"
    resolve_fields:
      resolution:
        name: Done
    reopen_fields:
      resolution: null
```

A resolve transition must lead to Jira's done status category, and a reopen
transition to a non-done category. Required transition fields are configurable.
These settings can be overridden per target, including explicit false. No issue
is created merely because a clean target has no known issue. Enabling closure
can reconcile an already-resolved issue when a previous local receipt retains
its tracker link. Changed destination policy or owner/environment metadata can
update a known finding even if its resource diff is unchanged.

Do not enable Jira closure without verifying project workflow IDs, permissions
and required fields. This configuration changes tracker workflow only; the
scanner still never applies infrastructure changes.

Workflow status is read from the current Jira issue before transition decisions;
search results are used only to locate the issue, avoiding stale-index recurrence decisions.
