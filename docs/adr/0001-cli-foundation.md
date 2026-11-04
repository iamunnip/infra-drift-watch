# 0001: Start with a Go CLI foundation

Status: Accepted
Date: 2026-11-04

## Context

The project needs incremental infrastructure detection with scheduled GitHub
Actions execution and no initial web UI.

## Decision

Use one Go module, a small cmd entry point, and internal implementation packages.
Use the Go standard library for foundation command handling. Target Go 1.27.3
for CI/container builds and declare Go 1.26 as the module minimum. Keep all
foundation tests offline. Follow TDD and add adapters only as their milestones begin.

## Consequences

There is no public Go API. The initial container supports help and version only;
scanning executables, CA certificates, writable state, and secret handling must
be designed before a scan-capable image is released. No UI or server is introduced.
