# OctoTerra Wizard

This tool prepares an Octopus space to allow space level resources and projects to serialized to a Terraform module and reapply them in another space.

* [Linux](https://github.com/mcasperson/OctoterraWizard/releases/latest/download/octoterrawiz_linux_amd64)
* [macOS](https://github.com/mcasperson/OctoterraWizard/releases/latest/download/octoterrawiz_macos_arm64)
* [Windows](https://github.com/mcasperson/OctoterraWizard/releases/latest/download/octoterrawiz_windows_amd64.exe)

## Environment variables

* `OCTOTERRAWIZ_BACKEND_TYPE`: Either `AWS S3` or `Azure Storage`
* `AWS_ACCESS_KEY_ID`: [AWS environment variable](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* `AWS_SECRET_ACCESS_KEY`: [AWS environment variable](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* `AWS_DEFAULT_REGION`: [AWS environment variable](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
* `AWS_DEFAULT_BUCKET`: The name of the S3 bucket holding the Terraform state
* `OCTOTERRAWIZ_PROMPT_FOR_DELETE`: If set to `true`, the tool will prompt for confirmation before deleting resources
* `OCTOTERRAWIZ_USE_CONTAINER_IMAGES`: If set to `true`, the tool will use container images to run Terraform steps
* `OCTOTERRAWIZ_AZURE_RESOURCE_GROUP`: The name of the Azure resource group holding the Terraform state
* `OCTOTERRAWIZ_AZURE_STORAGE_ACCOUNT`: The name of the Azure storage account holding the Terraform state
* `OCTOTERRAWIZ_AZURE_CONTAINER`: The name of the Azure storage container holding the Terraform state
* `AZURE_SUBSCRIPTION_ID`: [Azure environment variable](https://azure.github.io/static-web-apps-cli/docs/cli/env-vars/)
* `AZURE_TENANT_ID`: [Azure environment variable](https://azure.github.io/static-web-apps-cli/docs/cli/env-vars/)
* `AZURE_CLIENT_ID`: [Azure environment variable](https://azure.github.io/static-web-apps-cli/docs/cli/env-vars/)
* `AZURE_CLIENT_SECRET`: [Azure environment variable](https://azure.github.io/static-web-apps-cli/docs/cli/env-vars/)

## Screenshot

![](screenshot.png)