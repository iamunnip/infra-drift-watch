# Scan an existing AWS workspace

Use your installed OpenTofu/Terraform from PATH (Homebrew is supported). Optional
installer targets are for deliberately provisioned environments such as CI. Configure AWS
credentials through your normal credential chain. In GitHub Actions prefer OIDC.

The [example](../../../examples/aws/main.tf) models an S3 bucket. Select an existing
workspace with matching state, or provision a disposable example separately.
This scanner never creates or imports infrastructure. An empty state is not a
meaningful drift baseline and a refresh-only scan does not report unprovisioned configuration.

```sh
export TF_VAR_bucket_name=your-existing-demo-bucket
# Configure an existing backend in the workspace before initialization.
tofu -chdir=examples/aws init -input=false
bin/infra-drift-watch validate --config examples/config.yaml --inventory examples/inventory.yaml
bin/infra-drift-watch scan --config examples/config.yaml --inventory examples/inventory.yaml --dry-run
```

Scan invokes `plan -refresh-only -out=...` and `show -json`, never apply or refresh.
Grant required provider read permissions, backend read access and backend locking
permissions as needed. Do not grant deployment permissions just for scanning.
Keep backend credentials and variable secrets out of committed configuration.
Planning can run provider code and external data sources: scan trusted workspaces.

For Terraform, change adapter to terraform and initialize with Terraform itself.
Variable files are passed as dedicated arguments; no shell interpolation occurs.
