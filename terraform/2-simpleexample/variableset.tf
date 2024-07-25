resource "octopusdeploy_library_variable_set" "octopus_library_variable_set" {
  name = "Test"
  description = "Test variable set"
}

resource "octopusdeploy_variable" "secret" {
  name = "Test.SecretVariable"
  type = "Sensitive"
  description = "Test variable"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = "Development"
  scope {
    environments = [octopusdeploy_environment.development_environment.id]
  }
}

resource "octopusdeploy_variable" "secret2" {
  name = "Test.SecretVariable"
  type = "Sensitive"
  description = "Test variable"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = "Test"
  scope {
    environments = [octopusdeploy_environment.test_environment.id]
  }
}

resource "octopusdeploy_variable" "secret3" {
  name = "Test.SecretVariable"
  type = "Sensitive"
  description = "Test variable"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = "Test"
  scope {
    environments = [octopusdeploy_environment.production_environment.id]
  }
}

resource "octopusdeploy_variable" "unscoped" {
  name = "Test.UnscopedVariable"
  type = "Sensitive"
  description = "Test variable"
  is_sensitive = true
  is_editable = true
  owner_id = octopusdeploy_library_variable_set.octopus_library_variable_set.id
  value = "Unscoped"
}
