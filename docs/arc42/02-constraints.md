# Constraints

Go implementation; Linux runtime; development through develop. Existing initialized state workspaces are required. AWS hosted-runner connectivity and least-privilege credentials must be configured outside the scanner. Kubectl dry-run requires RBAC and may invoke admission webhooks. Only trusted IaC and executable plugins are in scope.

Go 1.27.3 is required. PR/manual CI performs full validation; no workflow runs on push. Local tool resolution uses PATH.
