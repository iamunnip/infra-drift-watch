---
type: Concept
title: AWS profile region and credential selection
description: How shared profiles and region selectors differ from credentials and deployment identity.
tags: [aws, profiles, regions, credentials]
---

# AWS profile, region and credential selection

A shared profile is a named configuration/credential chain, not a credential value.
A region selects an AWS service region; it does not identify an AWS account.
AWS environment credentials or explicit provider/backend settings can take
precedence over shared-profile values according to the invoked SDK's rules.
Selecting a profile alone is therefore not evidence of which account was accessed.

The scanner supplies optional profile and region selectors to its child tools;
it does not copy credentials into inventory or rewrite shared AWS files. Environment
and OIDC credentials remain available. EKS credential plugins inherit the same
subprocess selection when Kubernetes scanning is used.

A target ID represents a deployment scope and should remain stable across login
profile renaming, but must not be reused for a different account/region deployment.
See [AWS setup](../../docs/diataxis/how-to/scan-aws.md) and the
[configuration reference](../../docs/diataxis/reference/configuration.md).
