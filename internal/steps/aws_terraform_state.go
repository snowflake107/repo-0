package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/validators"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

type AwsTerraformStateStep struct {
	BaseStep
	Wizard    wizard.Wizard
	accessKey *widget.Entry
	secretKey *widget.Entry
	s3Bucket  *widget.Entry
	s3Region  *widget.Entry
	result    *widget.Label
}

func (s AwsTerraformStateStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDestinationDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	}, func() {
		s.result.SetText("")

		if !validators.ValidateAWS(s.getState()) {
			s.result.SetText("🔴 Unable to validate the credentials. Please check the Access Key, Secret Key, S3 Bucket Name, and S3 Bucket Region.")
			return
		}

		if !validators.TestS3Bucket(s.getState()) {
			s.result.SetText("🔴 Unable to connect to the S3 bucket. Please check that the bucket exists and that the supplied credentials can access it.")
			return
		}

		s.Wizard.ShowWizardStep(SpreadVariablesStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	})
	next.Disable()

	label1 := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		Terraform manages its state in an S3 bucket in AWS. Please provide the details of the S3 bucket that will be used to store the Terraform state.
	`))

	s.result = widget.NewLabel("")

	accessKeyLabel := widget.NewLabel("AWS Access Key")
	s.accessKey = widget.NewEntry()
	s.accessKey.SetPlaceHolder("")
	s.accessKey.SetText(s.State.AwsAccessKey)

	secretKeyLabel := widget.NewLabel("AWS Secret Key")
	s.secretKey = widget.NewPasswordEntry()
	s.secretKey.SetPlaceHolder("")
	s.secretKey.SetText(s.State.AwsSecretKey)

	s3BucketLabel := widget.NewLabel("AWS S3 Bucket Name")
	s.s3Bucket = widget.NewEntry()
	s.s3Bucket.SetPlaceHolder("my-bucket")
	s.s3Bucket.SetText(s.State.AwsS3Bucket)

	apiKeyLabel := widget.NewLabel("AWS S3 Bucket Region")
	s.s3Region = widget.NewEntry()
	s.s3Region.SetPlaceHolder("us-east-1")
	s.s3Region.SetText(s.State.AwsS3BucketRegion)

	validation := func(input string) {
		if s.accessKey != nil && s.accessKey.Text != "" && s.secretKey != nil && s.secretKey.Text != "" && s.s3Bucket != nil && s.s3Bucket.Text != "" && s.s3Region != nil && s.s3Region.Text != "" {
			next.Enable()
		} else {
			next.Disabled()
		}
	}

	validation("")

	s.accessKey.OnChanged = validation
	s.secretKey.OnChanged = validation
	s.s3Bucket.OnChanged = validation
	s.s3Region.OnChanged = validation

	formLayout := container.New(layout.NewFormLayout(), accessKeyLabel, s.accessKey, secretKeyLabel, s.secretKey, s3BucketLabel, s.s3Bucket, apiKeyLabel, s.s3Region)

	middle := container.New(layout.NewVBoxLayout(), label1, formLayout, s.result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s AwsTerraformStateStep) getState() state.State {
	return state.State{
		BackendType:             s.State.BackendType,
		Server:                  s.State.Server,
		ApiKey:                  s.State.ApiKey,
		Space:                   s.State.Space,
		DestinationServer:       s.State.DestinationServer,
		DestinationApiKey:       s.State.DestinationApiKey,
		DestinationSpace:        s.State.DestinationSpace,
		AwsAccessKey:            s.accessKey.Text,
		AwsSecretKey:            s.secretKey.Text,
		AwsS3Bucket:             s.s3Bucket.Text,
		AwsS3BucketRegion:       s.s3Region.Text,
		PromptForDelete:         s.State.PromptForDelete,
		AzureResourceGroupName:  s.State.AzureResourceGroupName,
		AzureStorageAccountName: s.State.AzureStorageAccountName,
		AzureContainerName:      s.State.AzureContainerName,
	}
}
