GO ?= go
.DEFAULT_GOAL := help
.PHONY: help build test test-race fmt fmt-check vet check docker-build
help:
	@echo 'Targets: build test test-race test-integration test-kubernetes fmt fmt-check vet lint check docs-check mermaid-check tools dev-tools vuln workflow-check docker-build'
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
check: fmt-check vet test-race build docs-check
docker-build:
	docker build -t infra-drift-watch:dev .

.PHONY: tools test-integration
tools:
	./scripts/install-tools.sh .tools/bin
test-integration:
	PATH="$(CURDIR)/.tools/bin:$$PATH" $(GO) test -tags=integration ./tests/integration -v

.PHONY: docs-check mermaid-check lint vuln dev-tools workflow-check test-kubernetes
docs-check:
	$(GO) run ./scripts/check-docs.go
mermaid-check:
	npm --prefix scripts/mermaid run check
lint: vet
dev-tools:
	GOBIN="$(CURDIR)/.tools/bin" $(GO) install golang.org/x/vuln/cmd/govulncheck@v1.9.0
	GOBIN="$(CURDIR)/.tools/bin" $(GO) install github.com/rhysd/actionlint/cmd/actionlint@v1.7.7
vuln:
	.tools/bin/govulncheck ./...
workflow-check:
	.tools/bin/actionlint -shellcheck=
test-kubernetes:
	PATH="$(CURDIR)/.tools/bin:$$PATH" $(GO) test -tags=kubernetes ./tests/integration -v -timeout 7m
.PHONY: kind-tools
kind-tools:
	./scripts/install-kind.sh .tools/bin
