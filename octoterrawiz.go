package main

import (
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/mcasperson/OctoterraWizard/internal/steps"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
)

func main() {
	wiz := wizard.NewWizard()
	wiz.ShowWizardStep(steps.WelcomeStep{
		Wizard:   *wiz,
		BaseStep: steps.BaseStep{State: state.State{}},
	})
	wiz.Window.ShowAndRun()
}
