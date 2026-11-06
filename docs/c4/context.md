# System context

C4 level 1: the scanner, the people and systems that start it, and the systems it inspects and writes to.

Operators supply inventory and config; GitHub Actions schedules scans. External infrastructure tools inspect AWS/state or Kubernetes. Delivery writes GitHub/Jira issues and Slack messages.

```mermaid
flowchart LR
    Operator[Operator] --> Watch[infra-drift-watch]
    Actions[GitHub Actions] --> Watch
    Watch --> StateTools[OpenTofu and Terraform]
    StateTools --> AWS[AWS and state backend]
    Watch --> Kube[Kubernetes API through kubectl]
    Watch --> GitHub[GitHub Issues]
    Watch --> Slack[Slack]
    Watch --> Jira[Jira Cloud]
```

All infrastructure and delivery destinations are external systems.

[C4 views](index.md)
