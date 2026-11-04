# Third-party notices

The original infra-drift-watch code is licensed under [MIT](LICENSE).
That license does not relicense third-party source, libraries, executables,
providers, documentation tooling, or operating-system packages.

This is a version-specific inventory of the named components, not a complete
software bill of materials or a replacement for upstream license terms.
[Verbatim license texts and provenance](licenses/README.md) accompany these notices.

## Linked and bundled software

| Component | Version | License | How it is used |
| --- | --- | --- | --- |
| Go runtime/standard library | 1.27.3 | BSD-3-Clause | Compiled into the CLI; [license](licenses/go-1.27.3-LICENSE.txt) |
| go.yaml.in/yaml/v3 | v3.0.6 | MIT and Apache-2.0, applying to different files | Linked configuration parser; [combined upstream notice](licenses/go-yaml-v3.0.6-LICENSE.txt) |
| OpenTofu | 1.13.5 | MPL-2.0 | Separate executable copied from the official minimal image; [license](licenses/opentofu-1.13.5-LICENSE.txt) |
| Terraform | 1.17.1 | BUSL-1.1 (Business Source License 1.1) | Separate executable in the runtime image; [version-specific license and parameters](licenses/terraform-1.17.1-LICENSE.txt) |
| kubectl | 1.37.2 | Apache-2.0 | Separate executable in the runtime image; [license](licenses/kubectl-1.37.2-LICENSE.txt) |
| GNU diffutils | 3.12 (Alpine package 3.12-r0) | GPL-3.0-or-later | Separate diff executable used by kubectl; [GPL v3 text](licenses/diffutils-3.12-GPL-3.0.txt) |
| Alpine Linux | 3.24.3 | Component-specific licenses | Base image, including musl, BusyBox and certificate data; inspect the installed APK package metadata |

The YAML dependency is not a choice between two licenses: its libyaml-derived
files use MIT, while its remaining files use Apache-2.0. Preserve both notices.
OpenTofu, Terraform and kubectl have their own dependencies and upstream notices.
The Go application calls their executables; it does not embed their source code.

## OpenTofu and Terraform are licensed differently

OpenTofu 1.13.5 uses MPL-2.0. Terraform 1.17.1 uses BSL-1.1, whose additional use
grant permits production use subject to the version's competitive-offering terms.
The Terraform license specifies a change to MPL-2.0 four years from publication,
subject to its full change-date provisions. Do not infer a calendar date from
this summary or describe the current Terraform release as MIT/MPL licensed.

Historical versions differ: Terraform 1.5.7 used MPL-2.0. The pinned current
Terraform license identifies Terraform 1.6.0 or later as its licensed work.
Consult the exact version's license rather than applying one label to all releases.
See the [licensing explanation](docs/diataxis/explanation/licensing.md) and
[upstream licensing FAQ](https://www.hashicorp.com/en/license-faq).

## Example and test providers

Providers are installed by the selected workspace or disposable test fixtures;
they are not shipped in the CLI or the default runtime image. Their licenses are
separate from both infrastructure engines.

| Provider | Selected version | License and source |
| --- | --- | --- |
| hashicorp/aws | ~> 6.71.0 example constraint | MPL-2.0 at [v6.71.0](https://github.com/hashicorp/terraform-provider-aws/blob/v6.71.0/LICENSE) |
| hashicorp/local | 2.9.1 test fixture | MPL-2.0 at [v2.9.1](https://github.com/hashicorp/terraform-provider-local/blob/v2.9.1/LICENSE) |

## Development and documentation tools

These tools are used for building, validating, or testing the repository rather
than shipped as CLI runtime dependencies.

| Component | Version | License/source |
| --- | --- | --- |
| Mermaid | 12.3.0 | MIT; [license](licenses/mermaid-12.3.0-LICENSE.txt) |
| jsdom | 30.1.3 | MIT; [license](licenses/jsdom-30.1.3-LICENSE.txt) |
| kind | 0.33.0 | Apache-2.0; [source](https://github.com/kubernetes-sigs/kind/tree/v0.33.0) |
| golangci-lint | 2.15.0 | GPL-3.0; [source](https://github.com/golangci/golangci-lint/tree/v2.15.0) |
| actionlint | 1.7.12 | MIT; [source](https://github.com/rhysd/actionlint/tree/v1.7.12) |
| govulncheck | 1.9.0 | BSD-3-Clause; [source](https://github.com/golang/vuln/tree/v1.9.0) |
| OpenSpec | 1.16.0 | MIT; [source](https://github.com/Fission-AI/OpenSpec/tree/v1.16.0) |

Transitive npm/Go dependencies, build images, and GitHub Actions retain their own
licenses. Version information is recorded in go.mod/go.sum, package-lock.json,
Dockerfile, installer scripts, and workflow files. AWS CLI is an optional
host authentication/troubleshooting tool, not bundled in this image.

## Source and distribution references

- [OpenTofu 1.13.5 source](https://github.com/opentofu/opentofu/tree/v1.13.5).
- [Terraform 1.17.1 source and terms](https://github.com/hashicorp/terraform/tree/v1.17.1).
- [Kubernetes 1.37.2 source](https://github.com/kubernetes/kubernetes/tree/v1.37.2).
- [Go YAML v3.0.6 source](https://github.com/yaml/go-yaml/tree/v3.0.6).
- [GNU diffutils 3.12 source archive](https://ftp.gnu.org/gnu/diffutils/diffutils-3.12.tar.xz)
  and [Alpine packaging/patches](https://gitlab.alpinelinux.org/alpine/aports/-/tree/3.24-stable/main/diffutils).
- [Alpine package sources](https://gitlab.alpinelinux.org/alpine/aports).

The image includes /LICENSE, /THIRD_PARTY_NOTICES.md, and the named license texts
under /usr/share/licenses/infra-drift-watch/. Release workflows attach a notices
archive alongside binaries and checksums. Preserve applicable notices and review
upstream terms when redistributing, modifying, or changing a component version.
SHA-256 download checks establish artifact integrity, not license equivalence or
independent signature verification.
