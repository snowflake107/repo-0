package main

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/samber/lo"
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
				BackendType:             "",
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

			if len(lvsVariable.Variables) != 7 {
				t.Fatalf("Expected 7 variables, got %v", len(lvsVariable.Variables))
			}

			// There must be 4 sensitive variables, and they must all be unscoped
			sensitiveVariables := lo.Filter(lvsVariable.Variables, func(item *variables.Variable, index int) bool {
				return item.IsSensitive &&
					len(item.Scope.Environments) == 0 &&
					len(item.Scope.Roles) == 0 &&
					len(item.Scope.Machines) == 0 &&
					len(item.Scope.Actions) == 0 &&
					len(item.Scope.TenantTags) == 0 &&
					len(item.Scope.ProcessOwners) == 0 &&
					len(item.Scope.Channels) == 0
			})

			if len(sensitiveVariables) != 4 {
				t.Fatalf("Expected 4 variables, got %v", len(sensitiveVariables))
			}

			// The three sensitive variables that shared a name must now have 3 regular variables each scoped
			// to an environment
			originalVariables := lo.Filter(lvsVariable.Variables, func(item *variables.Variable, index int) bool {
				return item.Name == "Test.SecretVariable" && !item.IsSensitive && len(item.Scope.Environments) == 1
			})

			if len(originalVariables) != 3 {
				t.Fatalf("Expected 3 variables, got %v", len(originalVariables))
			}
		}

		return nil
	})
}
