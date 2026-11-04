# Configure scheduled GitHub Actions scans

The workflow exists locally; enabling it is a deployment step you perform after
review. Schedules run from the repository default branch. Use release PRs from
develop to main before expecting the scheduled workflow to execute.

1. Commit your own global config, inventory and existing AWS backend configuration.
2. Set DRIFT_CONFIG_PATH and DRIFT_INVENTORY_PATH repository variables to those paths.
3. Set DRIFT_INIT_PATH to the initialized AWS workspace directory (one OpenTofu
   workspace is initialized by the initial workflow; extend initialization explicitly
   if inventory includes more workspaces or Terraform).
4. Configure AWS OIDC trust scoped to this repository/ref and set the
   DRIFT_AWS_ROLE_ARN secret; set DRIFT_AWS_REGION (default us-east-1).
5. For the example, set DRIFT_AWS_BUCKET_NAME and ensure its state already exists.
6. Set SLACK_WEBHOOK_URL, optional JIRA_EMAIL/JIRA_TOKEN, and optional
   DRIFT_GITHUB_TOKEN for a separate issue repository.
7. Set DRIFT_SCAN_ENABLED to true after reviewing permissions and destinations.

The default schedule is every six hours, UTC; edit the cron in scan.yml to change it.
The job's id-token permission allows AWS OIDC. Issues write permission applies to
its GitHub token. No static AWS credentials are committed.

Receipts restore from GitHub Actions cache and save under a unique run key,
including after failed delivery. A concurrency group serializes scheduled/manual
runs. Caches can expire or be evicted: tracker markers recover issues where possible,
but Slack reminders can repeat after lost receipts. Cache misses do not prove a
previous drift resolved; resolution notifications require prior active local state.
For stronger retention use a persistent external volume or explicit durable state
transfer under your own runner workflow; this release has no remote state backend.

A JSON report is retained for 14 days. Drift exit code 2 is a successful workflow
outcome; scan/coverage/delivery failures fail the workflow after state/report saving.
GitHub scheduling can be delayed. Private backend or API connectivity must be
verified from a hosted runner before enabling.

## Validation and releases

CI runs on pull_request targeting develop/main and workflow_dispatch, not push.
Go checks cover formatting, golangci-lint, vet, module verification, race tests,
and builds. Documentation checks cover local links, OKF knowledge metadata,
Mermaid parsing and strict OpenSpec validation. Integration checks cover native
adapters and Docker. These jobs require no live infrastructure or destination secrets.

Release binaries are produced only through workflow_dispatch with an existing
version tag on main. The workflow validates the tag and creates a draft release;
pushing a tag alone triggers nothing. Dependabot submits weekly dependency updates
targeting develop. Pin upgrades are reviewed and tested rather than using floating
latest tags at execution time.
