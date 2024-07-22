package steps

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"time"
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

type CustomCredentials struct {
	State state.State
}

func (c CustomCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.State.AwsAccessKey,
		SecretAccessKey: c.State.AwsSecretKey,
		SessionToken:    "",
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
		AccountID:       "",
	}, nil
}

func (s AwsTerraformStateStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDestinationDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	}, func() {
		s.result.SetText("")

		cfg, err := config.LoadDefaultConfig(
			context.Background(),
			config.WithCredentialsProvider(CustomCredentials{s.getState()}),
			config.WithRegion(s.getState().AwsS3BucketRegion))
		if err != nil {
			s.result.SetText("🔴 Unable to validate the AWS credentials.")
			return
		}

		simpleTokenService := sts.NewFromConfig(cfg)

		_, err = simpleTokenService.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
		if err != nil {
			s.result.SetText("🔴 Unable to validate the AWS credentials.")
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
		Server:            s.State.Server,
		ApiKey:            s.State.ApiKey,
		Space:             s.State.Space,
		DestinationServer: s.State.DestinationServer,
		DestinationApiKey: s.State.DestinationApiKey,
		DestinationSpace:  s.State.DestinationSpace,
		AwsS3Bucket:       s.s3Bucket.Text,
		AwsS3BucketRegion: s.s3Region.Text,
		AwsAccessKey:      s.accessKey.Text,
		AwsSecretKey:      s.secretKey.Text,
	}
}
