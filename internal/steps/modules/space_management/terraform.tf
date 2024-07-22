terraform {
  required_providers {
    octopusdeploy = { source = "OctopusDeployLabs/octopusdeploy", version = "0.22.0" }
    // Use the option below when debugging
    // octopusdeploy = { source = "octopus.com/com/octopusdeploy" }
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
variable "octopus_destination_server" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The URL of the Octopus server e.g. https://myinstance.octopus.app."
}
variable "octopus_destination_apikey" {
  type        = string
  nullable    = false
  sensitive   = true
  description = "The API key used to access the Octopus server. See https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key for details on creating an API key."
}
variable "octopus_destination_space_id" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The space ID to populate"
}


data "octopusdeploy_lifecycles" "lifecycle_default_lifecycle" {
  ids          = null
  partial_name = "Default Lifecycle"
  skip         = 0
  take         = 1
}

data "octopusdeploy_feeds" "built_in_feed" {
  feed_type    = "BuiltIn"
  ids          = null
  partial_name = ""
  skip         = 0
  take         = 1
}

resource "octopusdeploy_project_group" "octoterra" {
  name        = "Octoterra"
}

resource "octopusdeploy_library_variable_set" "octopus_library_variable_set" {
  name = "Octoterra"
  description = "Common variables used by Octoterra to deploy Octopus resources"
}

resource "octopusdeploy_variable" "destination_server" {
  name = "Octopus.Destination.Server"
  type = "String"
  description = "Octoterra destination server"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = var.octopus_destination_server
}

resource "octopusdeploy_variable" "destination_spaceid" {
  name = "Octopus.Destination.SpaceID"
  type = "String"
  description = "Octoterra destination server space ID"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = var.octopus_destination_space_id
}

resource "octopusdeploy_variable" "destination_api_key" {
  name = "Octopus.Destination.ApiKey"
  type = "Sensitive"
  description = "Octoterra destination server API key"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  sensitive_value = var.octopus_destination_apikey
}

resource "octopusdeploy_project" "deploy_frontend_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Test project"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Octoterra Space Management"
  project_group_id                     = octopusdeploy_project_group.octoterra.id
  tenanted_deployment_participation    = "Untenanted"
  included_library_variable_sets       = [octopusdeploy_library_variable_set.octopus_library_variable_set.id]
  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_runbook" "runbook" {
  project_id         = octopusdeploy_project.deploy_frontend_project.id
  name               = "Runbook"
  description        = "Test Runbook"
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
  environments                = null
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
      features              = ["Octopus.Features.JsonConfigurationVariables"]

      package {
        name                      = "package1"
        package_id                = "package1"
        acquisition_location      = "Server"
        extract_during_deployment = false
        feed_id                   = data.octopusdeploy_feeds.built_in_feed.feeds[0].id
        properties                = { Extract = "True", Purpose = "", SelectionMode = "immediate" }
      }
    }

    properties   = {}
    target_roles = []
  }
}