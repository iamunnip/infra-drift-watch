GO ?= go
GOVULNCHECK ?= govulncheck
ACTIONLINT ?= actionlint
GOLANGCI_LINT ?= golangci-lint
OPENSPEC ?= openspec
.DEFAULT_GOAL := help
.PHONY: help build test test-race fmt fmt-check vet check docker-build
help:
	@echo 'Targets: build test test-race test-integration test-kubernetes fmt fmt-check vet lint lint-docker mod-check check docs-check mermaid-check knowledge-check openspec-check tools kind-tools dev-tools vuln workflow-check docker-build'
build:
	$(GO) build -trimpath -o bin/infra-drift-watch ./cmd/infra-drift-watch
test:
	$(GO) test ./...
test-race:
	$(GO) test -race ./...
fmt:
	gofmt -w cmd internal scripts/*.go tests
fmt-check:
	@test -z "$$(gofmt -l cmd internal scripts/*.go tests)" || { gofmt -l cmd internal scripts/*.go tests; exit 1; }
vet:
	$(GO) vet ./...
check: fmt-check vet mod-check test-race build docs-check knowledge-check
docker-build:
	docker build -t infra-drift-watch:dev .

.PHONY: tools test-integration
tools:
	./scripts/install-tools.sh .tools/bin
test-integration:
	$(GO) test -tags=integration ./tests/integration -v

.PHONY: docs-check mermaid-check lint vuln dev-tools workflow-check test-kubernetes
docs-check:
	$(GO) run ./scripts/check-docs.go --scope docs
mermaid-check:
	npm --prefix scripts/mermaid run check
lint:
	$(GOLANGCI_LINT) run ./...
dev-tools:
	GOBIN="$(CURDIR)/.tools/bin" $(GO) install golang.org/x/vuln/cmd/govulncheck@v1.9.0
	GOBIN="$(CURDIR)/.tools/bin" $(GO) install github.com/rhysd/actionlint/cmd/actionlint@v1.7.12
vuln:
	$(GOVULNCHECK) ./...
workflow-check:
	$(ACTIONLINT) -shellcheck=
test-kubernetes:
	$(GO) test -tags=kubernetes ./tests/integration -v -timeout 7m
.PHONY: kind-tools
kind-tools:
	./scripts/install-kind.sh .tools/bin

.PHONY: knowledge-check openspec-check
knowledge-check:
	$(GO) run ./scripts/check-docs.go --scope knowledge
openspec-check:
	$(OPENSPEC) validate --all --strict

.PHONY: mod-check
mod-check:
	$(GO) mod verify
	$(GO) mod tidy -diff

.PHONY: lint-docker
lint-docker:
	docker run --rm --mount type=bind,src="$(CURDIR)",dst=/app,readonly --workdir /app golangci/golangci-lint:v2.15.0 golangci-lint run ./...
