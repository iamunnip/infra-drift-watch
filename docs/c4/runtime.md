# Runtime

A C4 dynamic view of one scan.

Validation performs no remote calls. Scan preflights enabled destination credentials and locks state before inspection. A failed target remains visible while later targets continue.

```mermaid
sequenceDiagram
    participant CLI
    participant App
    participant Tool
    participant Delivery
    participant State
    CLI->>App: Scan selected inventory
    App->>State: Lock state if delivery or history enabled
    App->>Tool: Inspect under timeout
    Tool-->>App: Result and coverage
    alt Complete clean or drift result
        App->>Delivery: Reconcile destination events
        Delivery->>State: Save each successful receipt
    else Failed or incomplete
        App->>App: Preserve unresolved findings
    end
    App->>State: Save bounded scan record if enabled
    App-->>CLI: Report and exit code
```

Explicit dry-run skips delivery, credential preflight and receipt and history writes, but still inspects infrastructure.

[C4 views](index.md)
