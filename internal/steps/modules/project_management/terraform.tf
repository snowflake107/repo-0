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
variable "octopus_project_name" {
  type        = string
  nullable    = false
  sensitive   = false
  description = "The project name to apply"
}

data "octopusdeploy_feeds" "docker_feed" {
  feed_type    = "Docker"
  ids          = null
  partial_name = "Octoterra Docker Feed"
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
  environment_scope           = "All"
  environments                = []
  default_guided_failure_mode = "EnvironmentDefault"
  force_package_download      = true
}

resource "octopusdeploy_runbook_process" "runbook" {
  runbook_id = octopusdeploy_runbook.runbook.id

  step {
    condition           = "Success"
    name                = "Octopus - Serialize Project to Terraform"
    package_requirement = "LetOctopusDecide"
    start_trigger       = "StartAfterPrevious"

    action {
      action_type                        = "Octopus.Script"
      name                               = "Octopus - Serialize Project to Terraform"
      condition                          = "Success"
      run_on_server                      = true
      is_disabled                        = false
      can_be_used_for_project_versioning = false
      is_required                        = false
      # Use the ubuntu worker pool if it is present, or use the default otherwise
      worker_pool_id = length(data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools) == 0 ? "" : data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools[0].id
      properties                         = {
        "SerializeProject.ThisInstance.Server.Url" = "#{Octopus.Source.Server}"
        "Octopus.Action.Template.Id" = var.octopus_serialize_actiontemplateid
        "SerializeProject.ThisInstance.Terraform.Backend" = "s3"
        "Octopus.Action.Template.Version" = "9"
        "SerializeProject.Exported.Project.Name" = "#{Octopus.Project.Name}"
        "Octopus.Action.Script.Syntax" = "Python"
        "Octopus.Action.RunOnServer" = "true"
        "Octopus.Action.Script.ScriptSource" = "Inline"
        "Octopus.Action.Script.ScriptBody" = "import argparse\nimport os\nimport re\nimport socket\nimport subprocess\nimport sys\nfrom datetime import datetime\nfrom urllib.parse import urlparse\nfrom itertools import chain\nimport platform\nfrom urllib.request import urlretrieve\nimport zipfile\n\n# If this script is not being run as part of an Octopus step, return variables from environment variables.\n# Periods are replaced with underscores, and the variable name is converted to uppercase\nif \"get_octopusvariable\" not in globals():\n    def get_octopusvariable(variable):\n        return os.environ[re.sub('\\\\.', '_', variable.upper())]\n\n# If this script is not being run as part of an Octopus step, print directly to std out.\nif \"printverbose\" not in globals():\n    def printverbose(msg):\n        print(msg)\n\n\ndef printverbose_noansi(output):\n    \"\"\"\n    Strip ANSI color codes and print the output as verbose\n    :param output: The output to print\n    \"\"\"\n    output_no_ansi = re.sub('\\x1b\\[[0-9;]*m', '', output)\n    printverbose(output_no_ansi)\n\n\ndef get_octopusvariable_quiet(variable):\n    \"\"\"\n    Gets an octopus variable, or an empty string if it does not exist.\n    :param variable: The variable name\n    :return: The variable value, or an empty string if the variable does not exist\n    \"\"\"\n    try:\n        return get_octopusvariable(variable)\n    except:\n        return ''\n\n\ndef execute(args, cwd=None, env=None, print_args=None, print_output=printverbose_noansi):\n    \"\"\"\n        The execute method provides the ability to execute external processes while capturing and returning the\n        output to std err and std out and exit code.\n    \"\"\"\n    process = subprocess.Popen(args,\n                               stdout=subprocess.PIPE,\n                               stderr=subprocess.PIPE,\n                               text=True,\n                               cwd=cwd,\n                               env=env)\n    stdout, stderr = process.communicate()\n    retcode = process.returncode\n\n    if print_args is not None:\n        print_output(' '.join(args))\n\n    if print_output is not None:\n        print_output(stdout)\n        print_output(stderr)\n\n    return stdout, stderr, retcode\n\n\ndef is_windows():\n    return platform.system() == 'Windows'\n\n\ndef init_argparse():\n    parser = argparse.ArgumentParser(\n        usage='%(prog)s [OPTION] [FILE]...',\n        description='Serialize an Octopus project to a Terraform module'\n    )\n    parser.add_argument('--ignore-all-changes',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IgnoreAllChanges') or get_octopusvariable_quiet(\n                            'Exported.Project.IgnoreAllChanges') or 'false',\n                        help='Set to true to set the \"lifecycle.ignore_changes\" ' +\n                             'setting on each exported resource to \"all\"')\n    parser.add_argument('--ignore-variable-changes',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IgnoreVariableChanges') or get_octopusvariable_quiet(\n                            'Exported.Project.IgnoreVariableChanges') or 'false',\n                        help='Set to true to set the \"lifecycle.ignore_changes\" ' +\n                             'setting on each exported octopus variable to \"all\"')\n    parser.add_argument('--terraform-backend',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.ThisInstance.Terraform.Backend') or get_octopusvariable_quiet(\n                            'ThisInstance.Terraform.Backend') or 'pg',\n                        help='Set this to the name of the Terraform backend to be included in the generated module.')\n    parser.add_argument('--server-url',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.ThisInstance.Server.Url') or get_octopusvariable_quiet(\n                            'ThisInstance.Server.Url'),\n                        help='Sets the server URL that holds the project to be serialized.')\n    parser.add_argument('--api-key',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.ThisInstance.Api.Key') or get_octopusvariable_quiet(\n                            'ThisInstance.Api.Key'),\n                        help='Sets the Octopus API key.')\n    parser.add_argument('--space-id',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Space.Id') or get_octopusvariable_quiet(\n                            'Exported.Space.Id') or get_octopusvariable_quiet('Octopus.Space.Id'),\n                        help='Set this to the space ID containing the project to be serialized.')\n    parser.add_argument('--project-name',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.Name') or get_octopusvariable_quiet(\n                            'Exported.Project.Name') or get_octopusvariable_quiet(\n                            'Octopus.Project.Name'),\n                        help='Set this to the name of the project to be serialized.')\n    parser.add_argument('--upload-space-id',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Octopus.UploadSpace.Id') or get_octopusvariable_quiet(\n                            'Octopus.UploadSpace.Id') or get_octopusvariable_quiet('Octopus.Space.Id'),\n                        help='Set this to the space ID of the Octopus space where ' +\n                             'the resulting package will be uploaded to.')\n    parser.add_argument('--ignore-cac-managed-values',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IgnoreCacValues') or get_octopusvariable_quiet(\n                            'Exported.Project.IgnoreCacValues') or 'false',\n                        help='Set this to true to exclude cac managed values like non-secret variables, ' +\n                             'deployment processes, and project versioning into the Terraform module. ' +\n                             'Set to false to have these values embedded into the module.')\n    parser.add_argument('--exclude-cac-project-settings',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.ExcludeCacProjectValues') or get_octopusvariable_quiet(\n                            'Exported.Project.ExcludeCacProjectValues') or 'false',\n                        help='Set this to true to exclude CaC settings like git connections from the exported module.')\n    parser.add_argument('--ignored-library-variable-sets',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IgnoredLibraryVariableSet') or get_octopusvariable_quiet(\n                            'Exported.Project.IgnoredLibraryVariableSet'),\n                        help='A comma separated list of library variable sets to ignore.')\n    parser.add_argument('--ignored-accounts',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IgnoredAccounts') or get_octopusvariable_quiet(\n                            'Exported.Project.IgnoredAccounts'),\n                        help='A comma separated list of accounts to ignore.')\n    parser.add_argument('--include-step-templates',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.IncludeStepTemplates') or get_octopusvariable_quiet(\n                            'Exported.Project.IncludeStepTemplates') or 'false',\n                        help='Set this to true to include step templates in the exported module. ' +\n                             'This disables the default behaviour of detaching step templates.')\n    parser.add_argument('--lookup-project-link-tenants',\n                        action='store',\n                        default=get_octopusvariable_quiet(\n                            'SerializeProject.Exported.Project.LookupProjectLinkTenants') or get_octopusvariable_quiet(\n                            'Exported.Project.LookupProjectLinkTenants') or 'false',\n                        help='Set this option to link tenants and create tenant project variables.')\n\n\n    return parser.parse_known_args()\n\n\ndef ensure_octo_cli_exists():\n    if is_windows():\n        print(\"Checking for the Octopus CLI\")\n        try:\n            stdout, _, exit_code = execute(['octo', 'help'])\n            printverbose(stdout)\n            if not exit_code == 0:\n                raise \"Octo CLI not found\"\n        except:\n            print(\"Downloading the Octopus CLI\")\n            urlretrieve('https://download.octopusdeploy.com/octopus-tools/9.0.0/OctopusTools.9.0.0.win-x64.zip',\n                        'OctopusTools.zip')\n            with zipfile.ZipFile('OctopusTools.zip', 'r') as zip_ref:\n                zip_ref.extractall(os.getcwd())\n\n\ndef check_docker_exists():\n    try:\n        stdout, _, exit_code = execute(['docker', 'version'])\n        printverbose(stdout)\n        if not exit_code == 0:\n            raise \"Docker not found\"\n    except:\n        print('Docker must be installed: https://docs.docker.com/get-docker/')\n        sys.exit(1)\n\n\ncheck_docker_exists()\nensure_octo_cli_exists()\nparser, _ = init_argparse()\n\n# Variable precondition checks\nif len(parser.server_url) == 0:\n    print(\"--server-url, ThisInstance.Server.Url, or SerializeProject.ThisInstance.Server.Url must be defined\")\n    sys.exit(1)\n\nif len(parser.api_key) == 0:\n    print(\"--api-key, ThisInstance.Api.Key, or ThisInstance.Api.Key must be defined\")\n    sys.exit(1)\n\noctoterra_image = 'ghcr.io/octopussolutionsengineering/octoterra-windows' if is_windows() else 'ghcr.io/octopussolutionsengineering/octoterra'\noctoterra_mount = 'C:/export' if is_windows() else '/export'\n\nprint(\"Pulling the Docker images\")\nexecute(['docker', 'pull', octoterra_image])\n\nif not is_windows():\n    execute(['docker', 'pull', 'ghcr.io/octopusdeploylabs/octo'])\n\n# Find out the IP address of the Octopus container\nparsed_url = urlparse(parser.server_url)\noctopus = socket.getaddrinfo(parsed_url.hostname, '80')[0][4][0]\n\nprint(\"Octopus hostname: \" + parsed_url.hostname)\nprint(\"Octopus IP: \" + octopus.strip())\n\n# Build the arguments to ignore library variable sets\nignores_library_variable_sets = parser.ignored_library_variable_sets.split(',')\nignores_library_variable_sets_args = [['-excludeLibraryVariableSet', x] for x in ignores_library_variable_sets]\n\n# Build the arguments to ignore accounts\nignored_accounts = parser.ignored_accounts.split(',')\nignored_accounts = [['-excludeAccounts', x] for x in ignored_accounts]\n\nos.mkdir(os.getcwd() + '/export')\n\nexport_args = ['docker', 'run',\n               '--rm',\n               '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n               '-v', os.getcwd() + '/export:' + octoterra_mount,\n               octoterra_image,\n               # the url of the instance\n               '-url', parser.server_url,\n               # the api key used to access the instance\n               '-apiKey', parser.api_key,\n               # add a postgres backend to the generated modules\n               '-terraformBackend', parser.terraform_backend,\n               # dump the generated HCL to the console\n               '-console',\n               # dump the project from the current space\n               '-space', parser.space_id,\n               # the name of the project to serialize\n               '-projectName', parser.project_name,\n               # ignoreProjectChanges can be set to ignore all changes to the project, variables, runbooks etc\n               '-ignoreProjectChanges=' + parser.ignore_all_changes,\n               # use data sources to lookup external dependencies (like environments, accounts etc) rather\n               # than serialize those external resources\n               '-lookupProjectDependencies',\n               # for any secret variables, add a default value set to the octostache value of the variable\n               # e.g. a secret variable called \"database\" has a default value of \"#{database}\"\n               '-defaultSecretVariableValues',\n               # detach any step templates, allowing the exported project to be used in a new space\n               '-detachProjectTemplates=' + str(not parser.include_step_templates),\n               # allow the downstream project to move between project groups\n               '-ignoreProjectGroupChanges',\n               # allow the downstream project to change names\n               '-ignoreProjectNameChanges',\n               # CaC enabled projects will not export the deployment process, non-secret variables, and other\n               # CaC managed project settings if ignoreCacManagedValues is true. It is usually desirable to\n               # set this value to true, but it is false here because CaC projects created by Terraform today\n               # save some variables in the database rather than writing them to the Git repo.\n               '-ignoreCacManagedValues=' + parser.ignore_cac_managed_values,\n               # Excluding CaC values means the resulting module does not include things like git credentials.\n               # Setting excludeCaCProjectSettings to true and ignoreCacManagedValues to false essentially\n               # converts a CaC project back to a database project.\n               '-excludeCaCProjectSettings=' + parser.exclude_cac_project_settings,\n               # This value is always true. Either this is an unmanaged project, in which case we are never\n               # reapplying it; or it is a variable configured project, in which case we need to ignore\n               # variable changes, or it is a shared CaC project, in which case we don't use Terraform to\n               # manage variables.\n               '-ignoreProjectVariableChanges=' + parser.ignore_variable_changes,\n               # To have secret variables available when applying a downstream project, they must be scoped\n               # to the Sync environment. But we do not need this scoping in the downstream project, so the\n               # Sync environment is removed from any variable scopes when serializing it to Terraform.\n               '-excludeVariableEnvironmentScopes', 'Sync',\n               # Exclude any variables starting with \"Private.\"\n               '-excludeProjectVariableRegex', 'Private\\\\..*',\n               # Capture the octopus endpoint, space ID, and space name as output vars. This is useful when\n               # querying th Terraform state file to know which space and instance the resources were\n               # created in. The scripts used to update downstream projects in bulk work by querying the\n               # Terraform state, finding all the downstream projects, and using the space name to only process\n               # resources that match the current tenant (because space names and tenant names are the same).\n               # The output variables added by this option are octopus_server, octopus_space_id, and\n               # octopus_space_name.\n               '-includeOctopusOutputVars',\n               # Where steps do not explicitly define a worker pool and reference the default one, this\n               # option explicitly exports the default worker pool by name. This means if two spaces have\n               # different default pools, the exported project still uses the pool that the original project\n               # used.\n               '-lookUpDefaultWorkerPools',\n               # Link any tenants that were originally link to the project and create project tenant variables\n               '-lookupProjectLinkTenants=' + parser.lookup_project_link_tenants,\n               # Add support for experimental step templates\n               '-experimentalEnableStepTemplates=' + parser.include_step_templates,\n               # The directory where the exported files will be saved\n               '-dest', octoterra_mount,\n               # This is a management runbook that we do not wish to export\n               '-excludeRunbookRegex', '__ .*'] + list(chain(*ignores_library_variable_sets_args)) + list(chain(*ignored_accounts))\n\nprint(\"Exporting Terraform module\")\n_, _, octoterra_exit = execute(export_args)\n\nif not octoterra_exit == 0:\n    print(\"Octoterra failed. Please check the logs for more information.\")\n    sys.exit(1)\n\ndate = datetime.now().strftime('%Y.%m.%d.%H%M%S')\n\nprint(\"Creating Terraform module package\")\nif is_windows():\n    execute(['octo',\n             'pack',\n             '--format', 'zip',\n             '--id', re.sub('[^0-9a-zA-Z]', '_', parser.project_name),\n             '--version', date,\n             '--basePath', os.getcwd() + '\\\\export',\n             '--outFolder', 'C:\\\\export'])\nelse:\n    _, _, _ = execute(['docker', 'run',\n                       '--rm',\n                       '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n                       '-v', os.getcwd() + \"/export:/export\",\n                       'ghcr.io/octopusdeploylabs/octo',\n                       'pack',\n                       '--format', 'zip',\n                       '--id', re.sub('[^0-9a-zA-Z]', '_', parser.project_name),\n                       '--version', date,\n                       '--basePath', '/export',\n                       '--outFolder', '/export'])\n\nprint(\"Uploading Terraform module package\")\nif is_windows():\n    _, _, _ = execute(['octo',\n                       'push',\n                       '--apiKey', parser.api_key,\n                       '--server', parser.server_url,\n                       '--space', parser.upload_space_id,\n                       '--package', 'C:\\\\export\\\\' +\n                       re.sub('[^0-9a-zA-Z]', '_', parser.project_name) + '.' + date + '.zip',\n                       '--replace-existing'])\nelse:\n    _, _, _ = execute(['docker', 'run',\n                       '--rm',\n                       '--add-host=' + parsed_url.hostname + ':' + octopus.strip(),\n                       '-v', os.getcwd() + \"/export:/export\",\n                       'ghcr.io/octopusdeploylabs/octo',\n                       'push',\n                       '--apiKey', parser.api_key,\n                       '--server', parser.server_url,\n                       '--space', parser.upload_space_id,\n                       '--package', '/export/' +\n                       re.sub('[^0-9a-zA-Z]', '_', parser.project_name) + '.' + date + '.zip',\n                       '--replace-existing'])\n\nprint(\"##octopus[stdout-default]\")\n\nprint(\"Done\")\n",
        "SerializeProject.Exported.Space.Id" = "#{Octopus.Space.Id}"
        "SerializeProject.Exported.Project.IgnoreVariableChanges" = "True"
        "SerializeProject.Exported.Project.IgnoreCacValues" = "False"
        "Exported.Project.IgnoredLibraryVariableSet" = "Octoterra"
        "SerializeProject.Exported.Project.IgnoreAllChanges" = "True"
        "SerializeProject.ThisInstance.Api.Key" = "#{Octopus.Source.ApiKey}"
        "SerializeProject.Exported.Project.IncludeStepTemplates" = "True"
        "SerializeProject.Exported.Project.IgnoredAccounts" = "Octoterra AWS Account"
        "SerializeProject.Exported.Project.ExcludeCacProjectValues" = "True"
        "SerializeProject.Exported.Project.LookupProjectLinkTenants" = "True"
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

resource "octopusdeploy_runbook" "deploy_project" {
  project_id         = var.octopus_project_id
  name               = "__ 2. Deploy Project"
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

resource "octopusdeploy_runbook_process" "deploy_project" {
  runbook_id = octopusdeploy_runbook.deploy_project.id

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
      worker_pool_id = length(data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools) == 0 ? "" : data.octopusdeploy_worker_pools.ubuntu_worker_pool.worker_pools[0].id
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
          "PackageId" = "${replace(var.octopus_project_name, "/[^A-Za-z0-9]/", "_")}"
          "FeedId" = "${data.octopusdeploy_feeds.built_in_feed.feeds[0].id}"
        })
        "OctoterraApply.Terraform.Workspace.Name" = "#{OctoterraApply.Octopus.SpaceID}"
        "OctoterraApply.Octopus.SpaceID" = "#{Octopus.Destination.SpaceID}"
        "OctopusUseBundledTooling" = "False"
        "Octopus.Action.AutoRetry.MaximumCount" = "3"
      }

#       container {
#         feed_id = data.octopusdeploy_feeds.docker_feed.feeds[0].id
#         image   = "ghcr.io/octopusdeploylabs/terraform-workertools"
#       }

      environments          = []
      excluded_environments = []
      channels              = []
      tenant_tags           = []

      primary_package {
        package_id           = "${replace(var.octopus_project_name, "/[^A-Za-z0-9]/", "_")}"
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