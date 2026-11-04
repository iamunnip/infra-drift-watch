# Project, engine, and provider licenses

infra-drift-watch's original code uses [MIT](../../../LICENSE). The downloadable
CLI links the YAML parser and Go runtime; the Docker image additionally bundles
independent infrastructure executables and OS packages. Their licenses remain
separate. [Third-party notices](../../../THIRD_PARTY_NOTICES.md) record the pinned
versions and link preserved license texts.

## OpenTofu versus Terraform

| Aspect | OpenTofu 1.13.5 | Terraform 1.17.1 |
| --- | --- | --- |
| Current license | MPL-2.0 | Business Source License 1.1, SPDX BUSL-1.1 |
| License model | Open-source, file-level copyleft | Source-available with production-use conditions |
| Version-specific terms | MPL notices and covered-source requirements | Additional use grant, competitive-offering terms, and change-date provisions |
| Future change specified in this license | No BSL-style timed conversion | Change license MPL-2.0 under the stated four-year provisions |
| Application integration here | Separate executable invoked for refresh-only plan/show | Separate executable invoked for refresh-only plan/show |

These summaries do not replace the exact [OpenTofu license](../../../licenses/opentofu-1.13.5-LICENSE.txt)
or [Terraform license](../../../licenses/terraform-1.17.1-LICENSE.txt).
The [Terraform licensing FAQ](https://www.hashicorp.com/en/license-faq) provides
upstream interpretation of its additional use grant.

Terraform 1.5.7 had an [MPL-2.0 license](https://github.com/hashicorp/terraform/blob/v1.5.7/LICENSE).
The current pinned license describes Terraform 1.6.0 or later as the licensed work.
A historical Terraform license does not describe the release bundled today.
The scanner's support for both engines does not make their license terms identical.

## Providers and other components

Engine licenses do not determine every provider's license. The AWS example uses
a provider constraint ~> 6.71.0; the local integration test pins provider 2.9.1.
Both cited provider versions carry MPL-2.0. Other providers require their own
version-specific notice review. Workspace initialization downloads providers;
they are not part of the default application image.

The YAML parser has MIT terms for its libyaml-derived files and Apache-2.0 for
its other files; it is not an either/or selection. kubectl uses Apache-2.0.
The image's GNU diffutils package uses GPL-3.0-or-later. Development tools such
as Mermaid/jsdom and golangci-lint are distinguished from shipped runtime tools
in the root notices; their dependencies have their own licenses too.

## Where notices ship

The repository stores unmodified license texts and source/checksum provenance
in [licenses](../../../licenses/README.md). The image includes these at
/usr/share/licenses/infra-drift-watch/, plus /LICENSE and /THIRD_PARTY_NOTICES.md.
The manual release workflow attaches third-party-notices.tar.gz, and includes its
checksum alongside binary checksums. This named-component inventory is not a
complete transitive dependency SBOM.

When a dependency, tool, or provider baseline changes, update the relevant notice,
license copy/provenance, container references, and release documentation together.
