package main

import (
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"os"
)

func main() {
	wiz := wizard.NewWizard()
	wiz.ShowWizardStep(steps.WelcomeStep{
		Wizard: *wiz,
		BaseStep: steps.BaseStep{State: state.State{
			BackendType: "",
			Server:      os.Getenv("OCTOPUS_CLI_SERVER"),
			ApiKey:      os.Getenv("OCTOPUS_CLI_API_KEY"),
			Space:       "Spaces-1",
		}},
	})
	wiz.Window.ShowAndRun()
}
