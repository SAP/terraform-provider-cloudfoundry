# Terraform Provider for Cloud Foundry

![Golang](https://img.shields.io/badge/Go-1.22-informational)
[![REUSE status](https://api.reuse.software/badge/github.com/SAP/terraform-provider-cloudfoundry)](https://api.reuse.software/info/github.com/SAP/terraform-provider-cloudfoundry)

## About This Project

The Terraform provider for [Cloud Foundry](https://www.cloudfoundry.org/) allows the management of resources via [Terraform](https://terraform.io/).

This provider makes use of the [go-cfclient](https://github.com/cloudfoundry/go-cfclient) to interact with the Cloud Foundry Cloud Controller [v3 APIs](https://v3-apidocs.cloudfoundry.org/version/3.159.0/index.html) and take advantages of the same. Additionally, the [v2 APIs are deprecated](https://apidocs.cloudfoundry.org/16.22.0/).

You can find usage examples in the [examples folder](examples/) of this repository.

Check the [Authentication](/Authentication.md) documentation for supported approaches.

## Developing & Contributing to the Provider

The [developer documentation](DEVELOPER.md) file is a basic outline on how to build and develop the provider. 

For more information about how to contribute, the project structure, and additional contribution information, see our [Contribution Guidelines](CONTRIBUTING.md).

## Prerequisites and Usage of the Provider

For the best experience using the Terraform Provider for Cloud Foundry, we recommend applying the common best practices for Terraform adoption as described in the [Hashicorp documentation](https://developer.hashicorp.com/well-architected-framework/operational-excellence/operational-excellence-terraform-maturity).

## Support and Feedback

❓ - If you have a *question* you can ask it here in the [GitHub Discussions](https://github.com/SAP/terraform-provider-cloudfoundry/discussions/).

🐞 - If you find a bug, feel free to create a [bug report](https://github.com/SAP/terraform-provider-cloudfoundry/issues/new?assignees=&labels=bug%2Cneeds-triage&projects=&template=bug_report.yml&title=%5BBUG%5D).

💡 - If you have an idea for improvement or a feature request, please open a [feature request](https://github.com/SAP/terraform-provider-cloudfoundry/issues/new?assignees=&labels=enhancement%2Cneeds-triage&projects=&template=feature_request.yml&title=%5BFEATURE%5D).

## Security / Disclosure

If you find any bug that may be a security problem, please follow our instructions at [in our security policy](https://github.com/SAP/terraform-provider-cloudfoundry/security/policy) on how to report it. Please do not create GitHub issues for security-related doubts or problems.

## Code of Conduct

Members, contributors, and leaders pledge to make participation in our community a harassment-free experience. By participating in this project, you agree to always abide by its [Code of Conduct](https://github.com/SAP/.github/blob/main/CODE_OF_CONDUCT.md).

## Licensing

Copyright 2024 SAP SE or an SAP affiliate company and `terraform-provider-cloudfoundry` contributors. See our [LICENSE](LICENSE) for copyright and license information. Detailed information, including third-party components and their licensing/copyright information, is available [via the REUSE tool](https://api.reuse.software/info/github.com/SAP/terraform-provider-cloudfoundry).
