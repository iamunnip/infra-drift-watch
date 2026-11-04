# Develop a change

Branch from `develop`. Write a behavioral test beside the package, run it and
confirm its expected failure, implement the change, and refactor.

```sh
go test ./internal/cli -run TestRun -v
make check
```

Update relevant specifications and documentation. See [contributing](../../../CONTRIBUTING.md).

Pull requests and manual dispatch run CI; pushes alone do not. make check includes
module verification and separate docs/knowledge validation. Run make lint or make
lint-docker for Go linting, make mermaid-check and make openspec-check for diagrams
and specifications. Installed tools are resolved from PATH, with no automatic
host installation or .tools/bin override.
