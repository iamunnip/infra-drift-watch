# Crosscutting concepts

Stable target IDs define issue identity. Each destination has a separate fingerprint and acknowledgement. Secret values are omitted; raw process stderr is discarded. SHA-256 fingerprints detect changes without retaining raw diff values. TLS and redirect restrictions protect credentials. Context cancellation bounds external work. Configuration is strict YAML with documented override semantics.

State format 2 preserves migrated receipts and bounded per-target history. History reads use atomic file snapshots without writes. Optional tracker closure is explicit; Jira transition destinations are validated. Kubernetes reports field paths and hashes ephemeral value differences without persisting values.

AWS profile and region selectors resolve globally, per target, and through CLI overrides. They are applied only to native-tool subprocesses; SDK credential precedence and explicit provider/backend settings remain in force. No credential files are read or rewritten by the application.
