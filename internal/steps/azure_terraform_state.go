package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/validators"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"strings"
)

type AzureTerraformStateStep struct {
	BaseStep
	Wizard             wizard.Wizard
	resourceGroupName  *widget.Entry
	storageAccountName *widget.Entry
	containerName      *widget.Entry
	result             *widget.Label
	subscriptionId     *widget.Entry
	tenantId           *widget.Entry
	applicationId      *widget.Entry
	password           *widget.Entry
	previous           *widget.Button
	next               *widget.Button
}

func (s AzureTerraformStateStep) GetContainer(parent fyne.Window) *fyne.Container {
	bottom, previous, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(BackendSelectionStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	}, func() {
		s.result.SetText("🔵 Validating Azure credentials and storage account container.")
		s.subscriptionId.Disable()
		s.containerName.Disable()
		s.tenantId.Disable()
		s.applicationId.Disable()
		s.storageAccountName.Disable()
		s.resourceGroupName.Disable()
		s.password.Disable()
		s.previous.Disable()
		s.next.Disable()

		defer s.subscriptionId.Enable()
		defer s.containerName.Enable()
		defer s.tenantId.Enable()
		defer s.applicationId.Enable()
		defer s.storageAccountName.Enable()
		defer s.resourceGroupName.Enable()
		defer s.password.Enable()
		defer s.previous.Enable()
		defer s.next.Enable()

		validationFailed := false
		newState := s.getState()
		exists, err := validators.AzureContainerExists(newState.AzureTenantId, newState.AzureApplicationId, newState.AzurePassword, newState.AzureStorageAccountName, newState.AzureContainerName)

		if err != nil {
			s.result.SetText("🔴 Unable to validate the credentials. Please check the credentials and storage account details.")
			validationFailed = true
		} else if !exists {
			s.result.SetText("🔴 Unable to find the Azure storage container.")
			validationFailed = true
		}

		rgExists, err := validators.AzureResourceGroupExists(newState.AzureTenantId, newState.AzureApplicationId, newState.AzureSubscriptionId, newState.AzurePassword, newState.AzureResourceGroupName)

		if err != nil {
			s.result.SetText("🔴 Unable to validate the credentials. Please check the credentials and storage account details.")
			validationFailed = true
		} else if !rgExists {
			s.result.SetText("🔴 Unable to find the Azure resource group.")
			validationFailed = true
		}

		nextCallback := func(result bool) {
			if result {
				s.Wizard.ShowWizardStep(SpreadVariablesStep{
					Wizard:   s.Wizard,
					BaseStep: BaseStep{State: s.getState()}})
			}
		}

		if validationFailed {
			dialog.NewConfirm("Azure Validation failed", "Validation of the Azure details failed. Do you wish to continue anyway?", nextCallback, s.Wizard.Window).Show()
		} else {
			nextCallback(true)
		}
	})
	s.next = next
	s.previous = previous

	s.next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Terraform manages its state in an storage account inAzure. Please provide the details of the storage account that will be used to store the Terraform state.
	`))

	s.result = widget.NewLabel("")

	subscriptionIdLabel := widget.NewLabel("Azure Subscription ID")
	s.subscriptionId = widget.NewEntry()
	s.subscriptionId.SetPlaceHolder("")
	s.subscriptionId.SetText(s.State.AzureSubscriptionId)

	tenantIdLabel := widget.NewLabel("Azure Tenant ID")
	s.tenantId = widget.NewEntry()
	s.tenantId.SetPlaceHolder("")
	s.tenantId.SetText(s.State.AzureTenantId)

	applicationIdLabel := widget.NewLabel("Azure Application ID")
	s.applicationId = widget.NewEntry()
	s.applicationId.SetPlaceHolder("")
	s.applicationId.SetText(s.State.AzureApplicationId)

	passwordLabel := widget.NewLabel("Azure Password")
	s.password = widget.NewPasswordEntry()
	s.password.SetPlaceHolder("")
	s.password.SetText(s.State.AzurePassword)

	azureResourceGroupLabel := widget.NewLabel("Azure Resource Group")
	s.resourceGroupName = widget.NewEntry()
	s.resourceGroupName.SetPlaceHolder("")
	s.resourceGroupName.SetText(s.State.AzureResourceGroupName)

	azureStorageAccountNameLabel := widget.NewLabel("Azure Storage Account Name")
	s.storageAccountName = widget.NewEntry()
	s.storageAccountName.SetPlaceHolder("")
	s.storageAccountName.SetText(s.State.AzureStorageAccountName)

	azureContainerNameLabel := widget.NewLabel("Azure Container Name")
	s.containerName = widget.NewEntry()
	s.containerName.SetPlaceHolder("my-container")
	s.containerName.SetText(s.State.AzureContainerName)

	validation := func(input string) {
		if s.resourceGroupName != nil && s.resourceGroupName.Text != "" && s.storageAccountName != nil && s.storageAccountName.Text != "" && s.containerName != nil && s.containerName.Text != "" {
			s.next.Enable()
		} else {
			s.next.Disabled()
		}
	}

	validation("")

	s.resourceGroupName.OnChanged = validation
	s.storageAccountName.OnChanged = validation
	s.containerName.OnChanged = validation

	formLayout := container.New(
		layout.NewFormLayout(),
		subscriptionIdLabel,
		s.subscriptionId,
		tenantIdLabel,
		s.tenantId,
		applicationIdLabel,
		s.applicationId,
		passwordLabel,
		s.password,
		azureResourceGroupLabel,
		s.resourceGroupName,
		azureStorageAccountNameLabel,
		s.storageAccountName,
		azureContainerNameLabel,
		s.containerName)

	middle := container.New(layout.NewVBoxLayout(), label1, formLayout, s.result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s AzureTerraformStateStep) getState() state.State {
	return state.State{
		BackendType:               s.State.BackendType,
		Server:                    s.State.Server,
		ServerExternal:            s.State.ServerExternal,
		ApiKey:                    s.State.ApiKey,
		Space:                     s.State.Space,
		DestinationServer:         s.State.DestinationServer,
		DestinationServerExternal: s.State.DestinationServerExternal,
		DestinationApiKey:         s.State.DestinationApiKey,
		DestinationSpace:          s.State.DestinationSpace,
		AwsAccessKey:              s.State.AwsAccessKey,
		AwsSecretKey:              s.State.AwsSecretKey,
		AwsS3Bucket:               s.State.AwsS3Bucket,
		AwsS3BucketRegion:         s.State.AwsS3BucketRegion,
		PromptForDelete:           s.State.PromptForDelete,
		UseContainerImages:        s.State.UseContainerImages,
		AzureResourceGroupName:    strings.TrimSpace(s.resourceGroupName.Text),
		AzureStorageAccountName:   strings.TrimSpace(s.storageAccountName.Text),
		AzureContainerName:        strings.TrimSpace(s.containerName.Text),
		AzureSubscriptionId:       strings.TrimSpace(s.State.AzureSubscriptionId),
		AzureTenantId:             strings.TrimSpace(s.tenantId.Text),
		AzureApplicationId:        strings.TrimSpace(s.applicationId.Text),
		AzurePassword:             strings.TrimSpace(s.password.Text),
	}
}
