package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"image/color"
	"os"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func main() {

	wiz := wizard.NewWizard()
	wiz.App.Settings().SetTheme(&myTheme{})

	wiz.ShowWizardStep(steps.WelcomeStep{
		Wizard: *wiz,
		BaseStep: steps.BaseStep{State: state.State{
			BackendType:       "",
			Server:            os.Getenv("OCTOPUS_CLI_SERVER"),
			ApiKey:            os.Getenv("OCTOPUS_CLI_API_KEY"),
			Space:             "Spaces-2048",
			DestinationServer: os.Getenv("OCTOPUS_CLI_SERVER"),
			DestinationApiKey: os.Getenv("OCTOPUS_CLI_API_KEY"),
			DestinationSpace:  "Spaces-2808",
		}},
	})
	wiz.Window.ShowAndRun()
}
