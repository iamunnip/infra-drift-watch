# Develop a change

Branch from `develop`. Write a behavioral test beside the package, run it and
confirm its expected failure, implement the change, and refactor.

```sh
go test ./internal/cli -run TestRun -v
make check
```

Update relevant specifications and documentation. See [contributing](../../../CONTRIBUTING.md).
