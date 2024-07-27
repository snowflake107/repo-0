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
	result             *widget.Label
	subscriptionId     *widget.Entry
	tenantId           *widget.Entry
	applicationId      *widget.Entry
	password           *widget.Entry
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
			next.Enable()
		} else {
			next.Disabled()
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
		ServerExternal:            "",
		ApiKey:                    s.State.ApiKey,
		Space:                     s.State.Space,
		DestinationServer:         s.State.DestinationServer,
		DestinationServerExternal: "",
		DestinationApiKey:         s.State.DestinationApiKey,
		DestinationSpace:          s.State.DestinationSpace,
		AwsAccessKey:              s.State.AwsAccessKey,
		AwsSecretKey:              s.State.AwsSecretKey,
		AwsS3Bucket:               s.State.AwsS3Bucket,
		AwsS3BucketRegion:         s.State.AwsS3BucketRegion,
		PromptForDelete:           s.State.PromptForDelete,
		UseContainerImages:        s.State.UseContainerImages,
		AzureResourceGroupName:    s.resourceGroupName.Text,
		AzureStorageAccountName:   s.storageAccountName.Text,
		AzureContainerName:        s.containerName.Text,
		AzureSubscriptionId:       s.State.AzureSubscriptionId,
		AzureTenantId:             s.State.AzureTenantId,
		AzureApplicationId:        s.State.AzureApplicationId,
		AzurePassword:             s.State.AzurePassword,
	}
}
