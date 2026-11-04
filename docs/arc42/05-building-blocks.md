# Building blocks

One binary is the sole C4 application container; the filesystem receipt store is its local data store. Internal packages form the component view:

```mermaid
flowchart TD
    Main[cmd entry point] --> CLI[cli]
    CLI --> App[app]
    App --> Config[config and inventory]
    App --> Adapters[adapters]
    Adapters --> Drift[drift model]
    App --> Report[report]
    App --> Delivery[delivery clients and coordinator]
    Delivery --> State[state receipts and history]
```

Configuration loading owns validation/override resolution. Adapters own comparison semantics. Delivery files separate HTTP, GitHub, Jira and coordination responsibilities without a plugin framework. State owns lock and atomic persistence.
