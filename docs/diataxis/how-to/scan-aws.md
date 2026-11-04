# Scan an existing AWS workspace

Use your installed OpenTofu/Terraform from PATH (Homebrew is supported). Optional
installer targets are for deliberately provisioned environments such as CI. Configure AWS
credentials through your normal credential chain. In GitHub Actions prefer OIDC.
If your credentials are already in the shared default profile, select profile
`default` and its region in configuration; credentials themselves stay outside Git.

```yaml
defaults:
  aws:
    profile: default
    region: us-east-1
```

Individual targets can override either setting under settings.aws, and a run can
override selected targets with --aws-profile/--aws-region. To opt out of an inherited
named profile (for example with environment/OIDC credentials), set profile to an
explicit empty string. Omitting it or using null inherits the global selection.

The [example](../../../examples/aws/main.tf) models an S3 bucket. Select an existing
workspace with matching state, or provision a disposable example separately.
This scanner never creates or imports infrastructure. An empty state is not a
meaningful drift baseline and a refresh-only scan does not report unprovisioned configuration.

```sh
export TF_VAR_bucket_name=your-existing-demo-bucket
# Configure an existing backend in the workspace before initialization.
tofu -chdir=examples/aws init -input=false
bin/infra-drift-watch validate --config examples/config.yaml --inventory examples/inventory.yaml
bin/infra-drift-watch scan --config examples/config.yaml --inventory examples/inventory.yaml --aws-profile default --aws-region us-east-1 --dry-run
```

Scan invokes `plan -refresh-only -out=...` and `show -json`, never apply or refresh.
Grant required provider read permissions, backend read access and backend locking
permissions as needed. Do not grant deployment permissions just for scanning.
Keep backend credentials and variable secrets out of committed configuration.
Planning can run provider code and external data sources: scan trusted workspaces.

For Terraform, change adapter to terraform and initialize with Terraform itself.
Variable files are passed as dedicated arguments; no shell interpolation occurs.

AWS selection is passed to both plan and show, and to kubectl credential plugins.
It does not change your shell environment or shared AWS files. Explicit provider
or backend settings, including credentials or region variables set in the workspace,
retain their own native precedence. The example's region variable defaults to null
so configured AWS region selectors are not hidden by a hardcoded provider region;
an explicit TF_VAR_region still overrides it.

Workspace initialization happens separately. When initializing a non-default
profile/region, supply the same selectors in that command's environment:

```sh
AWS_PROFILE=default AWS_DEFAULT_PROFILE=default AWS_REGION=us-east-1 AWS_DEFAULT_REGION=us-east-1 tofu -chdir=examples/aws init -input=false
```

See [AWS environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
and [SDK settings precedence](https://docs.aws.amazon.com/sdkref/latest/guide/settings-reference.html).
