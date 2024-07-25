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
variable "octopus_space_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The space name to populate"
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
variable "octopus_serialize_actiontemplateid" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The ID of the step template used to serialize a space"
}
variable "octopus_deploys3_actiontemplateid" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The ID of the step template used to deploy a space"
}
variable "terraform_state_bucket" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The S3 bucket used to save Terraform state"
}
variable "terraform_state_bucket_region" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The S3 bucket region used to save Terraform state"
}
variable "terraform_state_aws_accesskey" {
  type        = string
  nullable    = true
  sensitive   = false
  description = "The access key used to access the S3 bucket"
}
variable "terraform_state_aws_secretkey" {
  type        = string
  nullable    = true
  sensitive   = true
  description = "The access key used to access the S3 bucket"
}

data "octopusdeploy_library_variable_sets" "octoterra" {
  partial_name = "Octoterra"
  skip         = 0
  take         = 1
}

data "octopusdeploy_accounts" "aws" {
  account_type = "AmazonWebServicesAccount"
  ids          = []
  partial_name = "Octoterra AWS Account"
  skip         = 0
  take         = 1
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

data "octopusdeploy_worker_pools" "ubuntu_worker_pool" {
  name = "Hosted Ubuntu"
  ids  = null
  skip = 0
  take = 1
}

resource "octopusdeploy_aws_account" "account_aws_account" {
  name                              = "Octoterra AWS Account"
  description                       = ""
  environments                      = null
  tenant_tags                       = []
  tenants                           = null
  tenanted_deployment_participation = "Untenanted"
  access_key                        = var.terraform_state_aws_accesskey
  secret_key                        = var.terraform_state_aws_secretkey
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}


resource "octopusdeploy_project_group" "octoterra" {
  name        = "Octoterra"
}

resource "octopusdeploy_library_variable_set" "octopus_library_variable_set" {
  name = "Octoterra"
  description = "Common variables used by Octoterra to deploy Octopus resources"
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "destination_server" {
  name = "Octopus.Destination.Server"
  type = "String"
  description = "Octoterra destination server"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  value = var.octopus_destination_server
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "destination_spaceid" {
  name = "Octopus.Destination.SpaceID"
  type = "String"
  description = "Octoterra destination server space ID"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  value = var.octopus_destination_space_id
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "destination_api_key" {
  name = "Octopus.Destination.ApiKey"
  type = "Sensitive"
  description = "Octoterra destination server API key"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  sensitive_value = var.octopus_destination_apikey
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "source_server" {
  name = "Octopus.Source.Server"
  type = "String"
  description = "Octoterra source server"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  value = var.octopus_server
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "source_space" {
  name = "Octopus.Source.SpaceID"
  type = "String"
  description = "Octoterra source server"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  value = var.octopus_space_id
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "source_api_key" {
  name = "Octopus.Source.ApiKey"
  type = "Sensitive"
  description = "Octoterra source server API key"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  sensitive_value = var.octopus_apikey
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_variable" "aws_account" {
  name = "Terraform.AWS.Account"
  type = "AmazonWebServicesAccount"
  description = "Octoterra source server API key"
  is_sensitive = false
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set[0].id
  value = octopusdeploy_aws_account.account_aws_account[0].id
  count = length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? 1 : 0
}

resource "octopusdeploy_docker_container_registry" "feed_docker" {
  name                                 = "Octoterra Docker Feed"
  password                             = null
  username                             = null
  api_version                          = "v2"
  feed_uri                             = "https://ghcr.io"
  package_acquisition_location_options = ["ExecutionTarget", "NotAcquired"]
}

data "octopusdeploy_library_variable_sets" "all_variable_sets" {
  skip         = 0
  take         = 10000
}

resource "octopusdeploy_project" "space_management_project" {
  auto_create_release                  = false
  default_guided_failure_mode          = "EnvironmentDefault"
  default_to_skip_if_already_installed = false
  description                          = "Runbooks used to migrate and synchronize Octopus resources"
  discrete_channel_release             = false
  is_disabled                          = false
  is_discrete_channel_release          = false
  is_version_controlled                = false
  lifecycle_id                         = data.octopusdeploy_lifecycles.lifecycle_default_lifecycle.lifecycles[0].id
  name                                 = "Octoterra Space Management"
  project_group_id                     = octopusdeploy_project_group.octoterra.id
  tenanted_deployment_participation    = "Untenanted"
  included_library_variable_sets       = concat(data.octopusdeploy_library_variable_sets.all_variable_sets.library_variable_sets[*].id, [length(data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets) == 0 ? octopusdeploy_library_variable_set.octopus_library_variable_set[0].id : data.octopusdeploy_library_variable_sets.octoterra.library_variable_sets[0].id])

  versioning_strategy {
    template = "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.LastPatch}.#{Octopus.Version.NextRevision}"
  }

  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
}

resource "octopusdeploy_runbook" "serialize_space" {
  project_id         = octopusdeploy_project.space_management_project.id
  name               = "__ 1. Serialize Space"
  description        = "Serialize the space to a Terraform module"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "All"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "runbook" {
  runbook_id = octopusdeploy_runbook.serialize_space.id

  step {
    condition           = "Success"
    name                = "Octopus - Serialize Space to Terraform"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Octopus - Serialize Space to Terraform"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      # Use the ubuntu worker pool if it is present, or use the default otherwise
      worker_pool_id = len(data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools) == 0 ? "" : data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools[0].id
      properties                         = {
        "Octopus.Action.RunOnServer" = "true"
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.Syntax" = "Python"
        "SerializeSpace.ThisInstance.Terraform.Backend" = "s3"
        "SerializeSpace.Exported.Space.IgnoreTargets" = "False"
        "SerializeSpace.Exported.Space.Id" = "#{Octopus.Space.Id}"
        "Octopus.Action.Template.Version" = "4"
        "SerializeSpace.ThisInstance.Api.Key" = "#{Octopus.Source.ApiKey}"
        "Octopus.Action.Script.ScriptBody" = "import argparse\nimport os\nimport re\nimport socket\nimport subprocess\nimport sys\nfrom datetime import datetime\nfrom urllib.parse import urlparse\nfrom itertools import chain\nimport platform\nfrom urllib.request import urlretrieve\nimport zipfile\nimport urllib.request\nimport urllib.parse\nimport json\n\n# If this script is not being run as part of an Octopus step, return variables from environment variables.\n# Periods are replaced with underscores, and the variable name is converted to uppercase\nif \"get_octopusvariable\" not in globals():\n    def get_octopusvariable(variable):\n        return os.environ[re.sub('\\\\.', '_', variable.upper())]\n\n# If this script is not being run as part of an Octopus step, print directly to std out.\nif \"printverbose\" not in globals():\n    def printverbose(msg):\n        print(msg)\n\n\ndef printverbose_noansi(output):\n    \"\"\"\n    Strip ANSI color codes and print the output as verbose\n    :param output: The output to print\n    \"\"\"\n    output_no_ansi = re.sub(r'\\x1b\\[[0-9;]*m', '', output)\n    printverbose(output_no_ansi)\n\n\ndef get_octopusvariable_quiet(variable):\n    \"\"\"\n    Gets an octopus variable, or an empty string if it does not exist.\n    :param variable: The variable name\n    :return: The variable value, or an empty string if the variable does not exist\n    \"\"\"\n    try:\n        return get_octopusvariable(variable)\n    except:\n        return ''\n\n\ndef execute(args, cwd=None, env=None, print_args=None, print_output=printverbose_noansi):\n    \"\"\"\n        The execute method provides the ability to execute external processes while capturing and returning the\n        output to std err and std out and exit code.\n    \"\"\"\n    process = subprocess.Popen(args,\n                               stdout=subprocess.PIPE,\n                               stderr=subprocess.PIPE,\n                               text=True,\n                               cwd=cwd,\n                               env=env)\n    stdout, stderr = process.communicate()\n    retcode = process.returncode\n\n    if print_args is not None:\n        print_output(' '.join(args))\n\n    if print_output is not None:\n        print_output(stdout)\n        print_output(stderr)\n\n    return stdout, stderr, retcode\n\n\ndef is_windows():\n    return platform.system() == 'Windows'\n\n\ndef init_argparse():\n    parser = argparse.ArgumentParser(\n        usage='%(prog)s [OPTION] [FILE]...',\n        description='Serialize an Octopus project to a Terraform module'\n    )\n    parser.add_argument('--terraform-backend',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.ThisInstance.Terraform.Backend') or get_octopusvariable_quiet(\n                            'ThisInstance.Terraform.Backend') or 'pg',\n                        help='Set this to the name of the Terraform backend to be included in the generated module.')\n    parser.add_argument('--server-url',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.ThisInstance.Server.Url') or get_octopusvariable_quiet(\n                            'ThisInstance.Server.Url'),\n                        help='Sets the server URL that holds the project to be serialized.')\n    parser.add_argument('--api-key',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.ThisInstance.Api.Key') or get_octopusvariable_quiet(\n                            'ThisInstance.Api.Key'),\n                        help='Sets the Octopus API key.')\n    parser.add_argument('--space-id',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.Id') or get_octopusvariable_quiet(\n                            'Exported.Space.Id') or get_octopusvariable_quiet('Octopus.Space.Id'),\n                        help='Set this to the space ID containing the project to be serialized.')\n    parser.add_argument('--upload-space-id',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Octopus.UploadSpace.Id') or get_octopusvariable_quiet(\n                            'Octopus.UploadSpace.Id') or get_octopusvariable_quiet('Octopus.Space.Id'),\n                        help='Set this to the space ID of the Octopus space where ' +\n                             'the resulting package will be uploaded to.')\n    parser.add_argument('--ignored-library-variable-sets',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IgnoredLibraryVariableSet') or get_octopusvariable_quiet(\n                            'Exported.Space.IgnoredLibraryVariableSet'),\n                        help='A comma separated list of library variable sets to ignore.')\n    parser.add_argument('--ignored-tenants',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IgnoredTenants') or get_octopusvariable_quiet(\n                            'Exported.Space.IgnoredTenants'),\n                        help='A comma separated list of tenants ignore.')\n\n    parser.add_argument('--ignored-tenants-with-tag',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IgnoredTenantTags') or get_octopusvariable_quiet(\n                            'Exported.Space.IgnoredTenants'),\n                        help='A comma separated list of tenant tags that identify tenants to ignore.')\n    parser.add_argument('--ignore-all-targets',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IgnoreTargets') or get_octopusvariable_quiet(\n                            'Exported.Space.IgnoreAllChanges') or 'false',\n                        help='Set to true to exlude targets from the exported module')\n\n    parser.add_argument('--dummy-secret-variables',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.DummySecrets') or get_octopusvariable_quiet(\n                            'Exported.Space.DummySecrets') or 'false',\n                        help='Set to true to set secret values, like account and feed passwords, to a dummy value by default')\n    parser.add_argument('--include-step-templates',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IncludeStepTemplates') or get_octopusvariable_quiet(\n                            'Exported.Space.IncludeStepTemplates') or 'false',\n                        help='Set this to true to include step templates in the exported module. ' +\n                             'This disables the default behaviour of detaching step templates.')\n    parser.add_argument('--ignored-accounts',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeSpace.Exported.Space.IgnoredAccounts') or get_octopusvariable_quiet(\n                            'Exported.Space.IgnoredAccounts'),\n                        help='A comma separated list of accounts to ignore.')\n\n    return parser.parse_known_args()\n\n\ndef ensure_octo_cli_exists():\n    if is_windows():\n        print(\"Checking for the Octopus CLI\")\n        try:\n            stdout, _, exit_code = execute(['octo', 'help'])\n            printverbose(stdout)\n            if not exit_code == 0:\n                raise \"Octo CLI not found\"\n        except:\n            print(\"Downloading the Octopus CLI\")\n            urlretrieve('https://download.octopusdeploy.com/octopus-tools/9.0.0/OctopusTools.9.0.0.win-x64.zip',\n                        'OctopusTools.zip')\n            with zipfile.ZipFile('OctopusTools.zip', 'r') as zip_ref:\n                zip_ref.extractall(os.getcwd())\n\n\ndef check_docker_exists():\n    try:\n        stdout, _, exit_code = execute(['docker', 'version'])\n        printverbose(stdout)\n        if not exit_code == 0:\n            raise \"Docker not found\"\n    except:\n        print('Docker must be installed: https://docs.docker.com/get-docker/')\n        sys.exit(1)\n\n\ncheck_docker_exists()\nensure_octo_cli_exists()\nparser, _ = init_argparse()\n\n# Variable precondition checks\nif len(parser.server_url) == 0:\n    print(\"--server-url, ThisInstance.Server.Url, or SerializeSpace.ThisInstance.Server.Url must be defined\")\n    sys.exit(1)\n\nif len(parser.api_key) == 0:\n    print(\"--api-key, ThisInstance.Api.Key, or SerializeSpace.ThisInstance.Api.Key must be defined\")\n    sys.exit(1)\n\noctoterra_image = 'ghcr.io/octopussolutionsengineering/octoterra-windows' if is_windows() else 'ghcr.io/octopussolutionsengineering/octoterra'\noctoterra_mount = 'C:/export' if is_windows() else '/export'\n\nprint(\"Pulling the Docker images\")\nexecute(['docker', 'pull', octoterra_image])\n\nif not is_windows():\n    execute(['docker', 'pull', 'ghcr.io/octopusdeploylabs/octo'])\n\n# Find out the IP address of the Octopus container\nparsed_url = urlparse(parser.server_url)\noctopus = socket.getaddrinfo(parsed_url.hostname, '80')[0][4][0]\n\nprint(\"Octopus hostname: \" + parsed_url.hostname)\nprint(\"Octopus IP: \" + octopus.strip())\n\n# Build the arguments to ignore library variable sets\nignores_library_variable_sets = parser.ignored_library_variable_sets.split(',')\nignores_library_variable_sets_args = [['-excludeLibraryVariableSet', x] for x in ignores_library_variable_sets if x.strip() != '']\n\n# Build the arguments to ignore tenants\nignores_tenants = parser.ignored_tenants.split(',')\nignores_tenants_args = [['-excludeTenants', x] for x in ignores_tenants if x.strip() != '']\n\n# Build the arguments to ignore tenants with tags\nignored_tenants_with_tag = parser.ignored_tenants_with_tag.split(',')\nignored_tenants_with_tag_args = [['-excludeTenantsWithTag', x] for x in ignored_tenants_with_tag if x.strip() != '']\n\n# Build the arguments to ignore accounts\nignored_accounts = parser.ignored_accounts.split(',')\nignored_accounts = [['-excludeAccounts', x] for x in ignored_accounts]\n\nos.mkdir(os.getcwd() + '/export')\n\nexport_args = ['docker', 'run',\n               '--rm',\n               '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n               '-v', os.getcwd() + '/export:' + octoterra_mount,\n               octoterra_image,\n               # the url of the instance\n               '-url', parser.server_url,\n               # the api key used to access the instance\n               '-apiKey', parser.api_key,\n               # add a postgres backend to the generated modules\n               '-terraformBackend', parser.terraform_backend,\n               # dump the generated HCL to the console\n               '-console',\n               # dump the project from the current space\n               '-space', parser.space_id,\n               # Use default dummy values for secrets (e.g. a feed password). These values can still be overridden if known,\n               # but allows the module to be deployed and have the secrets updated manually later.\n               '-dummySecretVariableValues=' + parser.dummy_secret_variables,\n               # Add support for experimental step templates\n               '-experimentalEnableStepTemplates=' + parser.include_step_templates,\n               # Don't export any projects\n               '-excludeAllProjects',\n               # Output variables allow the Octopus space and instance to be determined from the Terraform state file.\n               '-includeOctopusOutputVars',\n               # Provide an option to ignore targets.\n               '-excludeAllTargets=' + parser.ignore_all_targets,\n               # The directory where the exported files will be saved\n               '-dest', octoterra_mount] + list(chain(*ignores_library_variable_sets_args, *ignores_tenants_args, *ignored_tenants_with_tag_args, *ignored_accounts))\n\nprint(\"Exporting Terraform module\")\n_, _, octoterra_exit = execute(export_args)\n\nif not octoterra_exit == 0:\n    print(\"Octoterra failed. Please check the verbose logs for more information.\")\n    sys.exit(1)\n\ndate = datetime.now().strftime('%Y.%m.%d.%H%M%S')\n\nprint('Looking up space name')\nurl = parser.server_url + '/api/Spaces/' + parser.space_id\nheaders = {\n    'X-Octopus-ApiKey': parser.api_key,\n    'Accept': 'application/json'\n}\nrequest = urllib.request.Request(url, headers=headers)\n\n# Retry the request for up to a minute.\nresponse = None\nfor x in range(12):\n    response = urllib.request.urlopen(request)\n    if response.getcode() == 200:\n        break\n    time.sleep(5)\n\nif not response or not response.getcode() == 200:\n    print('The API query failed')\n    sys.exit(1)\n\ndata = json.loads(response.read().decode(\"utf-8\"))\nprint('Space name is ' + data['Name'])\n\nprint(\"Creating Terraform module package\")\nif is_windows():\n    execute(['octo',\n             'pack',\n             '--format', 'zip',\n             '--id', re.sub('[^0-9a-zA-Z]', '_', data['Name']),\n             '--version', date,\n             '--basePath', os.getcwd() + '\\\\export',\n             '--outFolder', 'C:\\\\export'])\nelse:\n    _, _, _ = execute(['docker', 'run',\n                       '--rm',\n                       '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n                       '-v', os.getcwd() + \"/export:/export\",\n                       'ghcr.io/octopusdeploylabs/octo',\n                       'pack',\n                       '--format', 'zip',\n                       '--id', re.sub('[^0-9a-zA-Z]', '_', data['Name']),\n                       '--version', date,\n                       '--basePath', '/export',\n                       '--outFolder', '/export'])\n\nprint(\"Uploading Terraform module package\")\nif is_windows():\n    _, _, _ = execute(['octo',\n                       'push',\n                       '--apiKey', parser.api_key,\n                       '--server', parser.server_url,\n                       '--space', parser.upload_space_id,\n                       '--package', 'C:\\\\export\\\\' +\n                       re.sub('[^0-9a-zA-Z]', '_', data['Name']) + '.' + date + '.zip',\n                       '--replace-existing'])\nelse:\n    _, _, _ = execute(['docker', 'run',\n                       '--rm',\n                       '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n                       '-v', os.getcwd() + \"/export:/export\",\n                       'ghcr.io/octopusdeploylabs/octo',\n                       'push',\n                       '--apiKey', parser.api_key,\n                       '--server', parser.server_url,\n                       '--space', parser.upload_space_id,\n                       '--package', '/export/' +\n                       re.sub('[^0-9a-zA-Z]', '_', data['Name']) + '.' + date + '.zip',\n                       '--replace-existing'])\n\nprint(\"##octopus[stdout-default]\")\n\nprint(\"Done\")\n"
        "SerializeSpace.Exported.Space.DummySecrets" = "True"
        "SerializeSpace.Exported.Space.IgnoredLibraryVariableSet" = "Octoterra"
        "SerializeSpace.ThisInstance.Server.Url" = "#{Octopus.Source.Server}"
        "Octopus.Action.Template.Id" = var.octopus_serialize_actiontemplateid
        "SerializeSpace.Exported.Space.IncludeStepTemplates" = "True"
        "SerializeSpace.Exported.Space.IgnoredAccounts" = "Octoterra AWS Account"
      }
      environments                       = []
      excluded_environments              = []
      channels                           = []
      tenant_tags                        = []
      features                           = []
    }

    properties   = {}
    target_roles = []
  }
}

resource "octopusdeploy_runbook" "deploy_space" {
  project_id         = octopusdeploy_project.space_management_project.id
  name               = "__ 2. Deploy Space"
  description        = "Deploy the serialized Terraform module to a space"
  multi_tenancy_mode = "Untenanted"
  connectivity_policy {
    allow_deployments_to_no_targets = false
    exclude_unhealthy_targets       = false
    skip_machine_behavior           = "SkipUnavailableMachines"
  }
  retention_policy {
    quantity_to_keep = 10
  }
  environment_scope           = "All"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "deploy_space" {
  runbook_id = octopusdeploy_runbook.deploy_space.id

  step {
    condition           = "Success"
    name                = "Octopus - Populate Octoterra Space (S3 Backend)"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.TerraformApply"
      name                               = "Octopus - Populate Octoterra Space (S3 Backend)"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = true
      is_required                        = false
      # Use the ubuntu worker pool if it is present, or use the default otherwise
      worker_pool_id = len(data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools) == 0 ? "" : data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools[0].id
      properties                         = {
        "OctoterraApply.AWS.S3.BucketName" = var.terraform_state_bucket
        "OctoterraApply.AWS.S3.BucketRegion" =  var.terraform_state_bucket_region
        "OctoterraApply.AWS.Account" = "Terraform.AWS.Account"
        "OctoterraApply.AWS.S3.BucketKey" = "Project_#{Octopus.Project.Name | Replace \"[^A-Za-z0-9]\" \"_\"}"
        "Octopus.Action.Terraform.Workspace" = "#{OctoterraApply.Terraform.Workspace.Name}"
        "Octopus.Action.AwsAccount.UseInstanceRole" = "False"
        "Octopus.Action.AwsAccount.Variable" = "OctoterraApply.AWS.Account"
        "Octopus.Action.Aws.Region" = "#{OctoterraApply.AWS.S3.BucketRegion}"
        "Octopus.Action.Template.Id" = var.octopus_deploys3_actiontemplateid
        "Octopus.Action.Terraform.RunAutomaticFileSubstitution" = "False"
        "Octopus.Action.Terraform.AdditionalInitParams" = "-backend-config=\"bucket=#{OctoterraApply.AWS.S3.BucketName}\" -backend-config=\"region=#{OctoterraApply.AWS.S3.BucketRegion}\" -backend-config=\"key=#{OctoterraApply.AWS.S3.BucketKey}\" #{if OctoterraApply.Terraform.AdditionalInitParams}#{OctoterraApply.Terraform.AdditionalInitParams}#{/if}"
        "Octopus.Action.Terraform.TemplateDirectory" = "space_population"
        "Octopus.Action.Package.DownloadOnTentacle" = "False"
        "Octopus.Action.Terraform.AllowPluginDownloads" = "True"
        "OctoterraApply.Octopus.ServerUrl" = "#{Octopus.Destination.Server}"
        "Octopus.Action.RunOnServer" = "true"
        "Octopus.Action.Terraform.PlanJsonOutput" = "False"
        "Octopus.Action.Terraform.AzureAccount" = "False"
        "OctoterraApply.Octopus.ApiKey" = "#{Octopus.Destination.ApiKey}"
        "Octopus.Action.GoogleCloud.ImpersonateServiceAccount" = "False"
        "Octopus.Action.Terraform.GoogleCloudAccount" = "False"
        "Octopus.Action.Terraform.AdditionalActionParams" = "-var=octopus_server=#{OctoterraApply.Octopus.ServerUrl} -var=octopus_apikey=#{OctoterraApply.Octopus.ApiKey} -var=octopus_space_id=#{OctoterraApply.Octopus.SpaceID} #{if OctoterraApply.Terraform.AdditionalApplyParams}#{OctoterraApply.Terraform.AdditionalApplyParams}#{/if}"
        "Octopus.Action.Terraform.FileSubstitution" = "**/project_variable_sensitive*.tf"
        "Octopus.Action.Script.ScriptSource" = "Package"
        "Octopus.Action.Template.Version" = "3"
        "Octopus.Action.GoogleCloud.UseVMServiceAccount" = "True"
        "Octopus.Action.Terraform.ManagedAccount" = "AWS"
        "Octopus.Action.Aws.AssumeRole" = "False"
        "OctoterraApply.Terraform.Package.Id" = jsonencode({
          "PackageId" = "${replace(var.octopus_space_name, "/[^A-Za-z0-9]/", "_")}"
          "FeedId" = "${data.octopusdeploy_feeds.built_in_feed.feeds[0].id}"
        })
        "OctoterraApply.Terraform.Workspace.Name" = "#{OctoterraApply.Octopus.SpaceID}"
        "OctoterraApply.Octopus.SpaceID" = "#{Octopus.Destination.SpaceID}"
        "OctopusUseBundledTooling" = "False"
        "Octopus.Action.AutoRetry.MaximumCount" = "3"
      }

      container {
        feed_id = octopusdeploy_docker_container_registry.feed_docker.id
        image   = "ghcr.io/octopusdeploylabs/terraform-workertools"
      }

      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []

      primary_package {
        package_id           = "${replace(var.octopus_space_name, "/[^A-Za-z0-9]/", "_")}"
        acquisition_location = "Server"
        feed_id              = "${data.octopusdeploy_feeds.built_in_feed.feeds[0].id}"
        properties           = { PackageParameterName = "OctoterraApply.Terraform.Package.Id", SelectionMode = "deferred" }
      }

      features = []
    }

    properties   = {}
    target_roles = []
  }

}