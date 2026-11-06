# C4 architecture views

Each page holds one Mermaid diagram and the explanation that goes with it. Goals,
principles and risks are in [Design and risks](../diataxis/explanation/design-and-risks.md),
and the reasons for each choice are in the [ADRs](../adr/index.md).

- [System context](context.md): people, GitHub Actions and the external systems
- [Containers and components](building-blocks.md): the binary, its receipt store and internal packages
- [Runtime](runtime.md): the order of one scan, failure handling and dry-run
- [Deployment](deployment.md): GitHub-hosted runners, OIDC, caches and the Docker image

Update the matching view in the same PR when integrations, package responsibilities,
scan order or hosting change.
