package steps

import (
	"fyne.io/fyne/v2/test"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"testing"
)

func TestStartSpaceExportStep_GetContainer(t *testing.T) {
	return

	testWindow := test.NewWindow(nil)
	defer testWindow.Close()

	step := StartSpaceExportStep{
		Wizard: wizard.Wizard{},
		BaseStep: BaseStep{State: state.State{
			BackendType:             "",
			Server:                  "http://172.17.0.1:8080/",
			ApiKey:                  "API-AAAAA",
			Space:                   "Spaces-1",
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

	if err := step.Execute(); err != nil {
		t.Error(err)
	}
}
