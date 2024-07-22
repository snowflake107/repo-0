package steps

import (
	"bytes"
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"github.com/samber/lo"
	"image/color"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed modules/space_management/terraform.tf
var module string

type SpaceExportStep struct {
	BaseStep
	Wizard        wizard.Wizard
	createProject *widget.Button
}

func (s SpaceExportStep) GetContainer() *fyne.Container {

	bottom, previous, next := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(WelcomeStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(OctopusDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	next.Disable()

	intro := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		We now must create a project with runbooks to serialize the space to a Terraform module and reapply it to a new space.
		This project is called "Octoterra Space Management" in the project group "Octoterra".
		Click the "Create Project" button to create the project and its associated runbooks.
	`))
	infinite := widget.NewProgressBarInfinite()
	infinite.Start()
	infinite.Hide()
	result := canvas.NewText("", color.White)
	s.createProject = widget.NewButton("Create Project", func() {
		next.Disable()
		previous.Disable()
		infinite.Show()
		s.createProject.Disable()
		result.Text = "Creating project. This can take a little while."

		go func() {
			defer previous.Enable()
			defer infinite.Hide()
			myclient, err := octoclient.CreateClient(s.State)

			if err != nil {
				result.Text = "Failed to create the client:\n" + err.Error()
				return
			}

			// Best effort at deleting existing project and project group
			if project, err := projects.GetByName(myclient, myclient.GetSpaceID(), "Octoterra Space Management"); err == nil {
				if err := myclient.Projects.DeleteByID(project.ID); err != nil {
					fmt.Print(err)
				}
			}

			if projectGroups, err := projectgroups.GetAll(myclient, myclient.GetSpaceID()); err == nil {
				projectGroup := lo.Filter(projectGroups, func(pg *projectgroups.ProjectGroup, index int) bool {
					return pg.Name == "Octoterra"
				})
				if len(projectGroup) != 0 {
					if err := myclient.ProjectGroups.DeleteByID(projectGroup[0].ID); err != nil {
						fmt.Print(err)
					}
				}
			}

			// Save and apply the module
			dir, err := ioutil.TempDir("", "octoterra")
			if err != nil {
				result.Text = "An error occurred while creating a temporary directory:\n" + err.Error()
				return
			}

			filePath := filepath.Join(dir, "terraform.tf")

			if err := os.WriteFile(filePath, []byte(module), 0644); err != nil {
				result.Text = "An error occurred while writing the Terraform file:\n" + err.Error()
				return
			}

			initCmd := exec.Command("terraform", "init")
			initCmd.Dir = dir

			if err := initCmd.Run(); err != nil {
				result.Text = "terraform init failed."
				return
			}

			applyCmd := exec.Command("terraform",
				"apply",
				"-auto-approve",
				"-no-color",
				"-var=octopus_server="+s.State.Server,
				"-var=octopus_apikey="+s.State.ApiKey,
				"-var=octopus_space_id="+s.State.Space)
			applyCmd.Dir = dir

			var stdout, stderr bytes.Buffer
			applyCmd.Stdout = &stdout
			applyCmd.Stderr = &stderr

			if err := applyCmd.Run(); err != nil {
				result.Text = "terraform apply failed:\n" + stderr.String()
				return
			} else {
				result.Text = "terraform apply succeeded"
			}

			next.Enable()
		}()
	})
	middle := container.New(layout.NewVBoxLayout(), intro, s.createProject, infinite, result)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}
