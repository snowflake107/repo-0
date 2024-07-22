package wizard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Wizard struct {
	App    fyne.App
	Window fyne.Window
}

func NewWizard(title string) *Wizard {
	newApp := app.New()
	window := newApp.NewWindow(title)
	window.Resize(fyne.NewSize(800, 600))

	return &Wizard{
		App:    newApp,
		Window: window,
	}
}

type WizardStep interface {
	GetContainer(parent fyne.Window) *fyne.Container
}

func (w *Wizard) ShowWizardStep(step WizardStep) {
	w.Window.SetContent(step.GetContainer(w.Window))
}
