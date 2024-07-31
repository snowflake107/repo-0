resource "octopusdeploy_channel" "backend_mainline" {
  name        = "Test"
  project_id  = octopusdeploy_project.deploy_frontend_project.id
  description = "Test channel"
  depends_on  = [octopusdeploy_project.deploy_frontend_project]
  is_default  = true
}