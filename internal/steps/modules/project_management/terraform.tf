terraform {
  required_providers {
    octopusdeploy = { source = "OctopusDeployLabs/octopusdeploy", version = "0.22.0" }
  }
}

provider "octopusdeploy" {
  address  = var.octopus_server
  api_key  = var.octopus_apikey
  space_id = var.octopus_space_id
}

variable "octopus_server" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The URL of the Octopus server e.g. https://myinstance.octopus.app."
}
variable "octopus_apikey" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The API key used to access the Octopus server. See https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key for details on creating an API key."
}
variable "octopus_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The space ID to populate"
}
variable "octopus_project_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The project to add the runbook to"
}

resource "octopusdeploy_runbook" "runbook" {
  project_id         = var.octopus_project_id
  name               = "__ 1. Serialize Project"
  description        = "Serialize the project to a Terraform module"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "Specified"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "runbook" {
  runbook_id = octopusdeploy_runbook.runbook.id

  step {
    condition           = "Success"
    name                = "Hello world (using PowerShell)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Hello world (using PowerShell)"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = true
      worker_pool_id                     = ""
      properties                         = {
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.ScriptBody"   = "Write-Host 'Hello world, using PowerShell'\n\n#TODO: Experiment with steps of your own :)\n\nWrite-Host '[Learn more about the types of steps available in Octopus](https://oc.to/OnboardingAddStepsLearnMore)'"
        "Octopus.Action.Script.Syntax"       = "PowerShell"
      }
      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []
      features              = []
    }

    properties   = {}
    target_roles = []
  }
}