package main

import (
	"bytes"
	"encoding/json"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/mcasperson/OctoterraWizard/internal/infrastructure"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/samber/lo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestSpreadVariables(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(
			t,
			container,
			filepath.Join("terraform"),
			"2-simpleexample",
			[]string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		step := steps.SpreadVariablesStep{
			BaseStep: steps.BaseStep{State: state.State{
				BackendType:             "AWS S3",
				Server:                  container.URI,
				ApiKey:                  test.ApiKey,
				Space:                   newSpaceId,
				DestinationServer:       "",
				DestinationApiKey:       "",
				DestinationSpace:        "",
				AwsAccessKey:            "",
				AwsSecretKey:            "",
				AwsS3Bucket:             "",
				AwsS3BucketRegion:       "",
				PromptForDelete:         false,
				AzureResourceGroupName:  "",
				AzureStorageAccountName: "",
				AzureContainerName:      "",
				AzureSubscriptionId:     "",
				AzureTenantId:           "",
				AzureApplicationId:      "",
				AzurePassword:           "",
			}},
		}

		// we must be able to repeat this step with no changes
		for i := 0; i < 3; i++ {
			if err := step.Execute(); err != nil {
				t.Fatalf("Error executing step: %v", err)
			}

			lvs, err := newSpaceClient.LibraryVariableSets.GetAll()

			if err != nil {
				t.Fatalf("Error getting library variable sets: %v", err)
			}

			lvsVariable, err := variables.GetVariableSet(newSpaceClient, newSpaceClient.GetSpaceID(), lvs[0].VariableSetID)

			if err != nil {
				t.Fatalf("Error getting library variable sets: %v", err)
			}

			if len(lvsVariable.Variables) != 10 {
				t.Fatalf("Expected 10 variables, got %v", len(lvsVariable.Variables))
			}

			// There must be one regular variable that was unaltered
			if len(lo.Filter(lvsVariable.Variables, func(item *variables.Variable, index int) bool {
				return !item.IsSensitive && item.Name == "RegularVariable"
			})) != 1 {
				t.Fatalf("Expected 1 regular variable")
			}

			// There must be 5 sensitive variables, and they must all be unscoped
			sensitiveVariables := lo.Filter(lvsVariable.Variables, func(item *variables.Variable, index int) bool {
				return item.IsSensitive && item.Scope.IsEmpty()
			})

			if len(sensitiveVariables) != 5 {
				t.Fatalf("Expected 5 variables, got %v", len(sensitiveVariables))
			}

			// The four sensitive variables that shared a name must now have 4 regular variables each scoped
			// to an environment or are unscoped
			originalVariables := lo.Filter(lvsVariable.Variables, func(item *variables.Variable, index int) bool {
				return item.Name == "Test.SecretVariable" && !item.IsSensitive && (len(item.Scope.Environments) == 1 || len(item.Scope.Environments) == 0)
			})

			if len(originalVariables) != 4 {
				t.Fatalf("Expected 4 variables, got %v", len(originalVariables))
			}

			// Each regular variable must reference a sensitive variable
			for _, variable := range originalVariables {
				matchingSensitiveVar := lo.Filter(sensitiveVariables, func(item *variables.Variable, index int) bool {
					return *variable.Value == "#{"+item.Name+"}"
				})

				if len(matchingSensitiveVar) == 0 {
					t.Fatalf("Should have found a matching sensitive variable for %v", variable.Name)
				}
			}

		}

		return nil
	})
}

func TestProjectSpreadVariables(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(
			t,
			container,
			filepath.Join("terraform"),
			"3-simpleproject",
			[]string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		if err != nil {
			return err
		}

		step := steps.SpreadVariablesStep{
			BaseStep: steps.BaseStep{State: state.State{
				BackendType:             "AWS S3",
				Server:                  container.URI,
				ApiKey:                  test.ApiKey,
				Space:                   newSpaceId,
				DestinationServer:       "",
				DestinationApiKey:       "",
				DestinationSpace:        "",
				AwsAccessKey:            "",
				AwsSecretKey:            "",
				AwsS3Bucket:             "",
				AwsS3BucketRegion:       "",
				PromptForDelete:         false,
				AzureResourceGroupName:  "",
				AzureStorageAccountName: "",
				AzureContainerName:      "",
				AzureSubscriptionId:     "",
				AzureTenantId:           "",
				AzureApplicationId:      "",
				AzurePassword:           "",
			}},
		}

		// we must be able to repeat this step with no changes
		for i := 0; i < 3; i++ {
			if err := step.Execute(); err != nil {
				t.Fatalf("Error executing step: %v", err)
			}

			project, err := projects.GetByName(newSpaceClient, newSpaceClient.GetSpaceID(), "Test")

			if err != nil {
				t.Fatalf("Error getting library project: %v", err)
			}

			variableSet, err := variables.GetVariableSet(newSpaceClient, newSpaceClient.GetSpaceID(), project.VariableSetID)

			if err != nil {
				t.Fatalf("Error getting project variable set: %v", err)
			}

			if len(variableSet.Variables) != 9 {
				t.Fatalf("Expected 9 variables, got %v", len(variableSet.Variables))
			}

			// There must be 4 sensitive variables, and they must all be unscoped
			sensitiveVariables := lo.Filter(variableSet.Variables, func(item *variables.Variable, index int) bool {
				return item.IsSensitive && item.Scope.IsEmpty()
			})

			if len(sensitiveVariables) != 4 {
				t.Fatalf("Expected 4 variables, got %v", len(sensitiveVariables))
			}

			// The three sensitive variables that shared a name must now have 4 regular variables each scoped
			// to an environment or are unscoped
			originalVariables := lo.Filter(variableSet.Variables, func(item *variables.Variable, index int) bool {
				return item.Name == "SensitiveVariable" && !item.IsSensitive && (len(item.Scope.Environments) == 1 || len(item.Scope.Environments) == 0)
			})

			if len(originalVariables) != 4 {
				t.Fatalf("Expected 4 variables, got %v", len(originalVariables))
			}

			// Each regular variable must reference a sensitive variable
			for _, variable := range originalVariables {
				matchingSensitiveVar := lo.Filter(sensitiveVariables, func(item *variables.Variable, index int) bool {
					return *variable.Value == "#{"+item.Name+"}"
				})

				if len(matchingSensitiveVar) == 0 {
					t.Fatalf("Should have found a matching sensitive variable for %v", variable.Name)
				}
			}

			// There must be a regular variable that was not altered
			regularVariable := lo.Filter(variableSet.Variables, func(item *variables.Variable, index int) bool {
				return item.Name == "RegularVariable" && !item.IsSensitive
			})

			if len(regularVariable) != 1 {
				t.Fatalf("Expected 1 variable, got %v", len(regularVariable))
			}

		}

		return nil
	})
}

// TestProjectMigration tests a full space migration running the wizard steps in the same sequence
// a user would from the UI
func TestProjectMigration(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, client *client.Client) error {
		// Act
		newSpaceId, err := testFramework.Act(
			t,
			container,
			filepath.Join("terraform"),
			"3-simpleproject",
			[]string{})

		if err != nil {
			return err
		}

		newSpaceClient, err := octoclient.CreateClient(container.URI, newSpaceId, test.ApiKey)

		newSpace := spaces.NewSpace("Migration")
		newSpace.SpaceManagersTeams = []string{"teams-administrators"}
		space, err := newSpaceClient.Spaces.Add(newSpace)

		if err != nil {
			return err
		}

		state := state.State{
			BackendType:               "AWS S3",
			Server:                    "http://172.17.0.1:8080", // The address used by Octopus when running tasks, which could be in nested containers
			ServerExternal:            container.URI,            // The address used by the wizard
			ApiKey:                    test.ApiKey,
			Space:                     newSpaceId,
			DestinationServer:         "http://172.17.0.1:8080",
			DestinationServerExternal: container.URI,
			DestinationApiKey:         test.ApiKey,
			DestinationSpace:          space.ID,
			AwsAccessKey:              os.Getenv("AWS_ACCESS_KEY_ID"),
			AwsSecretKey:              os.Getenv("AWS_SECRET_ACCESS_KEY"),
			AwsS3Bucket:               os.Getenv("AWS_DEFAULT_BUCKET"),
			AwsS3BucketRegion:         os.Getenv("AWS_DEFAULT_REGION"),
			PromptForDelete:           false,
			UseContainerImages:        false,
			AzureResourceGroupName:    "",
			AzureStorageAccountName:   "",
			AzureContainerName:        "",
			AzureSubscriptionId:       "",
			AzureTenantId:             "",
			AzureApplicationId:        "",
			AzurePassword:             "",
		}

		// need to install pip and terraform onto the Octopus container
		installPip := "{\"Name\":\"AdHocScript\",\"Description\":\"Script run from management console\",\"Arguments\":{\"MachineIds\":[],\"TenantIds\":[],\"TargetRoles\":[],\"EnvironmentIds\":[],\"WorkerIds\":[],\"WorkerPoolIds\":[],\"TargetType\":\"OctopusServer\",\"Syntax\":\"Bash\",\"ScriptBody\":\"apt-get update\\napt-get upgrade\\napt install python3-pip -y\\napt-get update && apt-get install -y gnupg software-properties-common\\nwget -O- https://apt.releases.hashicorp.com/gpg | \\\\\\ngpg --dearmor | \\\\\\ntee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null\\necho \\\"deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \\\\\\nhttps://apt.releases.hashicorp.com $(lsb_release -cs) main\\\" | \\\\\\ntee /etc/apt/sources.list.d/hashicorp.list\\napt update\\napt-get install -y terraform\"},\"SpaceId\":\"Spaces-2\"}"
		req, err := http.NewRequest("POST", container.URI+"/api/tasks", bytes.NewReader([]byte(installPip)))

		if err != nil {
			return err
		}

		response, err := newSpaceClient.HttpSession().DoRawRequest(req)

		if err != nil {
			return err
		}

		responseBody, err := io.ReadAll(response.Body)

		if err != nil {
			return err
		}

		stepTemplates := map[string]any{}
		if err := json.Unmarshal(responseBody, &stepTemplates); err != nil {
			return err
		}

		err = infrastructure.WaitForTask(state, stepTemplates["Id"].(string), func(message string) {})

		if err != nil {
			return err
		}

		if err := (steps.SpreadVariablesStep{BaseStep: steps.BaseStep{State: state}}).Execute(); err != nil {
			t.Fatalf("Error executing SpreadVariablesStep: %v", err)
		}

		if _, err := (steps.StepTemplateStep{BaseStep: steps.BaseStep{State: state}}).Execute(); err != nil {
			t.Fatalf("Error executing StepTemplateStep: %v", err)
		}

		steps.SpaceExportStep{BaseStep: steps.BaseStep{State: state}}.Execute(func(message string, body string, callback func(bool)) {
			callback(true)
		}, func(s string, err error) {
			t.Fatalf("Error executing SpaceExportStep: %v", err)
		}, func(s string) {
			// success
		}, false, false)

		steps.ProjectExportStep{BaseStep: steps.BaseStep{State: state}}.Execute(func(message string, body string, callback func(bool)) {
			callback(true)
		}, func(s string, err error) {
			t.Fatalf("Error executing ProjectExportStep: %v", err)
		}, func(s string) {
			// success
		}, func(s string) {
			// status
		})

		if err := (steps.StartSpaceExportStep{BaseStep: steps.BaseStep{State: state}}).Execute(func(message string) {}); err != nil {
			t.Fatalf("Error executing StartSpaceExportStep: %v", err)
		}

		if err := (steps.StartProjectExportStep{BaseStep: steps.BaseStep{State: state}}).Execute(func(message string) {}); err != nil {
			t.Fatalf("Error executing StartProjectExportStep: %v", err)
		}

		migratedSpaceClient, err := octoclient.CreateClient(container.URI, space.ID, test.ApiKey)

		if err != nil {
			return err
		}

		project, err := projects.GetByName(migratedSpaceClient, migratedSpaceClient.GetSpaceID(), "Test")

		if err != nil {
			return err
		}

		if project.Description != "Test project" {
			t.Fatalf("Expected description to be 'Test project', got %v", project.Description)
		}

		return nil
	})
}
