# First run without cloud credentials

Install Go 1.27.3 and Make. From the repository root:

```sh
make check
bin/infra-drift-watch --help
bin/infra-drift-watch version
bin/infra-drift-watch validate --config examples/config.yaml --inventory examples/inventory.yaml
```

Validation should report one target. It creates no cloud resources, contacts no
delivery services, and writes no receipts. Continue with the
[AWS scan guide](../how-to/scan-aws.md) when an initialized workspace is available.
