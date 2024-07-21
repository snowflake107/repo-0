package steps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"net/url"
)

type OctopusDetails struct {
	BaseStep
	Wizard wizard.Wizard
}

func (s OctopusDetails) GetContainer() *fyne.Container {

	bottom, _, next := s.BuildNavigation(func() {}, func() {
		s.Wizard.ShowWizardStep(TestTerraformStep{Wizard: s.Wizard})
	})
	next.Disable()

	label1 := widget.NewLabel("Enter the URL, API key, and Space ID of the Octopus instance you want to export.")
	linkUrl, _ := url.Parse("https://octopus.com/docs/octopus-rest-api/how-to-create-an-api-key")
	link := widget.NewHyperlink("Learn how to create an API key.", linkUrl)

	serverLabel := widget.NewLabel("Server URL")
	server := widget.NewEntry()
	server.SetPlaceHolder("https://octopus.example.com")

	apiKeyLabel := widget.NewLabel("API Key")
	apiKey := widget.NewEntry()
	apiKey.SetPlaceHolder("API-xxxxxxxxxxxxxxxxxxxxxxxxxx")

	spaceIdLabel := widget.NewLabel("Space ID")
	spaceId := widget.NewEntry()
	spaceId.SetPlaceHolder("Spaces-#")

	validation := func(input string) {
		if server.Text != "" && apiKey.Text != "" && spaceId.Text != "" {
			next.Enable()
		} else {
			next.Disabled()
		}
	}

	server.OnChanged = validation
	apiKey.OnChanged = validation
	spaceId.OnChanged = validation

	formLayout := container.New(layout.NewFormLayout(), serverLabel, server, apiKeyLabel, apiKey, spaceIdLabel, spaceId)

	middle := container.New(layout.NewVBoxLayout(), label1, link, formLayout)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
