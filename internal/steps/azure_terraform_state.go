package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type AzureTerraformStateStep struct {
	BaseStep
	Wizard             wizard.Wizard
	resourceGroupName  *widget.Entry
	storageAccountName *widget.Entry
	containerName      *widget.Entry
	keyName            *widget.Entry
	result             *widget.Label
}

func (s AzureTerraformStateStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDestinationDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	}, func() {
		s.Wizard.ShowWizardStep(SpreadVariablesStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Terraform manages its state in an storage account inAzure. Please provide the details of the storage account that will be used to store the Terraform state.
	`))

	s.result = widget.NewLabel("")

	accessKeyLabel := widget.NewLabel("Azure Resource Group")
	s.resourceGroupName = widget.NewEntry()
	s.resourceGroupName.SetPlaceHolder("")
	s.resourceGroupName.SetText(s.State.AzureResourceGroupName)

	secretKeyLabel := widget.NewLabel("Azure Storage Account Name")
	s.storageAccountName = widget.NewPasswordEntry()
	s.storageAccountName.SetPlaceHolder("")
	s.storageAccountName.SetText(s.State.AzureStorageAccountName)

	s3BucketLabel := widget.NewLabel("Azure Container Name")
	s.containerName = widget.NewEntry()
	s.containerName.SetPlaceHolder("my-container")
	s.containerName.SetText(s.State.AzureContainerName)

	apiKeyLabel := widget.NewLabel("Azure Key Name")
	s.keyName = widget.NewEntry()
	s.keyName.SetPlaceHolder("mykey")
	s.keyName.SetText(s.State.AzureKeyName)

	validation := func(input string) {
		if s.resourceGroupName != nil && s.resourceGroupName.Text != "" && s.storageAccountName != nil && s.storageAccountName.Text != "" && s.containerName != nil && s.containerName.Text != "" && s.keyName != nil && s.keyName.Text != "" {
			next.Enable()
		} else {
			next.Disabled()
		}
	}

	validation("")

	s.resourceGroupName.OnChanged = validation
	s.storageAccountName.OnChanged = validation
	s.containerName.OnChanged = validation
	s.keyName.OnChanged = validation

	formLayout := container.New(layout.NewFormLayout(), accessKeyLabel, s.resourceGroupName, secretKeyLabel, s.storageAccountName, s3BucketLabel, s.containerName, apiKeyLabel, s.keyName)

	middle := container.New(layout.NewVBoxLayout(), label1, formLayout, s.result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s AzureTerraformStateStep) getState() state.State {
	return state.State{
		BackendType:             s.State.BackendType,
		Server:                  s.State.Server,
		ApiKey:                  s.State.ApiKey,
		Space:                   s.State.Space,
		DestinationServer:       s.State.DestinationServer,
		DestinationApiKey:       s.State.DestinationApiKey,
		DestinationSpace:        s.State.DestinationSpace,
		AwsAccessKey:            s.State.AwsAccessKey,
		AwsSecretKey:            s.State.AwsSecretKey,
		AwsS3Bucket:             s.State.AwsS3Bucket,
		AwsS3BucketRegion:       s.State.AwsS3BucketRegion,
		PromptForDelete:         s.State.PromptForDelete,
		AzureResourceGroupName:  s.resourceGroupName.Text,
		AzureStorageAccountName: s.storageAccountName.Text,
		AzureContainerName:      s.containerName.Text,
		AzureKeyName:            s.keyName.Text,
	}
}
