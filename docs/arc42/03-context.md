# Context and scope

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

This is the C4 system context view. All infrastructure and delivery destinations are external systems.
