package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/validators"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"net/url"
)

type OctopusDestinationDetails struct {
	BaseStep
	Wizard  wizard.Wizard
	server  *widget.Entry
	apiKey  *widget.Entry
	spaceId *widget.Entry
	result  *widget.Label
}

func (s OctopusDestinationDetails) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	}, func() {
		s.result.SetText("")

		if !validators.ValidateDestinationCreds(s.getState()) {
			s.result.SetText("🔴 Unable to connect to the Octopus server. Please check the URL, API key, and Space ID.")
			return
		}

		s.Wizard.ShowWizardStep(AwsTerraformStateStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.getState()}})
	})

	s.result = widget.NewLabel("")

	validation := func(input string) {
		next.Disabled()

		if s.server != nil && s.server.Text != "" && s.apiKey != nil && s.apiKey.Text != "" && s.spaceId != nil && s.spaceId.Text != "" {
			return
		}

		next.Enable()
	}

	validation("")

	introText := widget.NewLabel("Enter the URL, API key, and Space ID of the Octopus instance you want to export to (i.e. the destination server).")
	linkUrl, _ := url.Parse("https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key")
	link := widget.NewHyperlink("Learn how to create an API key.", linkUrl)

	serverLabel := widget.NewLabel("Destination Server URL")
	s.server = widget.NewEntry()
	s.server.SetPlaceHolder("https://octopus.example.com")
	s.server.SetText(s.State.DestinationServer)

	apiKeyLabel := widget.NewLabel("Destination API Key")
	s.apiKey = widget.NewPasswordEntry()
	s.apiKey.SetPlaceHolder("API-xxxxxxxxxxxxxxxxxxxxxxxxxx")
	s.apiKey.SetText(s.State.DestinationApiKey)

	spaceIdLabel := widget.NewLabel("Destination Space ID")
	s.spaceId = widget.NewEntry()
	s.spaceId.SetPlaceHolder("Spaces-#")
	s.spaceId.SetText(s.State.DestinationSpace)

	s.server.OnChanged = validation
	s.apiKey.OnChanged = validation
	s.spaceId.OnChanged = validation

	formLayout := container.New(layout.NewFormLayout(), serverLabel, s.server, apiKeyLabel, s.apiKey, spaceIdLabel, s.spaceId)

	middle := container.New(layout.NewVBoxLayout(), introText, link, formLayout, s.result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s OctopusDestinationDetails) getState() state.State {
	return state.State{
		BackendType:       "",
		Server:            s.State.Server,
		ApiKey:            s.State.ApiKey,
		Space:             s.State.Space,
		DestinationServer: s.server.Text,
		DestinationApiKey: s.apiKey.Text,
		DestinationSpace:  s.spaceId.Text,
		AwsS3Bucket:       s.State.AwsS3Bucket,
		AwsS3BucketRegion: s.State.AwsS3BucketRegion,
		AwsAccessKey:      s.State.AwsAccessKey,
		AwsSecretKey:      s.State.AwsSecretKey,
		PromptForDelete:   s.State.PromptForDelete,
	}
}
