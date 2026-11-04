# Third-party software

The project's own code is MIT licensed. Bundled tools retain their upstream
licenses; this project's license does not relicense them.

- YAML Go module: https://github.com/yaml/go-yaml (MIT and Apache-2.0).
- OpenTofu: https://github.com/opentofu/opentofu (MPL-2.0).
- Terraform: https://github.com/hashicorp/terraform (Business Source License 1.1
  for the pinned release; consult its license for permitted use).
- kubectl: https://github.com/kubernetes/kubernetes (Apache-2.0).
- Alpine packages retain the license information distributed by Alpine.

Pinned runtime versions are declared in scripts/install-tools.sh. Downloads
verify published SHA-256 checksums over HTTPS. This is not independent signature
verification. Reassess upstream releases and licenses when updating bundled tools.
