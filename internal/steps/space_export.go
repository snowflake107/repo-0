package steps

import (
	"bytes"
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projectgroups"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"github.com/samber/lo"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed modules/space_management/terraform.tf
var module string

var spaceManagementProject = "Octoterra Space Management"

type SpaceExportStep struct {
	BaseStep
	Wizard        wizard.Wizard
	createProject *widget.Button
	infinite      *widget.ProgressBarInfinite
	result        *widget.Label
	logs          *widget.Entry
	next          *widget.Button
	previous      *widget.Button
}

func (s SpaceExportStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, thisPrevious, thisNext := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(OctopusDestinationDetails{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(ProjectExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	s.next = thisNext
	s.previous = thisPrevious
	s.next.Disable()

	intro := widget.NewLabel(strutil.TrimMultilineWhitespace(`
		We now must create a project with runbooks to serialize the space to a Terraform module and reapply it to a new space.
		This project is called "Octoterra Space Management" in the project group "Octoterra".
		Click the "Create Project" button to create the project and its associated runbooks.
	`))
	s.infinite = widget.NewProgressBarInfinite()
	s.infinite.Start()
	s.infinite.Hide()
	s.result = widget.NewLabel("")
	s.logs = widget.NewEntry()
	s.logs.Disable()
	s.logs.MultiLine = true
	s.logs.SetMinRowsVisible(20)
	s.createProject = widget.NewButton("Create Project", func() { s.createNewProject(parent) })
	middle := container.New(layout.NewVBoxLayout(), intro, s.createProject, s.infinite, s.result, s.logs)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s SpaceExportStep) createNewProject(parent fyne.Window) {
	s.logs.SetText("")
	s.next.Disable()
	s.previous.Disable()
	s.infinite.Show()
	s.createProject.Disable()
	s.result.SetText("Creating project. This can take a little while.")

	go func() {
		defer s.previous.Enable()
		defer s.infinite.Hide()
		defer s.createProject.Enable()

		myclient, err := octoclient.CreateClient(s.State)

		if err != nil {
			s.logs.SetText("Failed to create the client:\n" + err.Error())
			return
		}

		// Best effort at deleting existing project and project group
		projExists, project, err := s.projectExists(myclient)

		if projExists {
			dialog.NewConfirm("Project Group Exists", "The project "+spaceManagementProject+" already exists. Do you want to delete it? It is usually safe to delete this resource.", func(b bool) {
				if b {
					if err := s.deleteProject(myclient, project); err != nil {
						s.result.SetText("Failed to delete the resource")
						s.logs.SetText(err.Error())
					} else {
						s.createNewProject(parent)
					}
				}
			}, parent).Show()

			// We can't go further until the group is deleted
			return
		}

		pgExists, pggroup, err := s.projectGroupExists(myclient)

		if pgExists {
			dialog.NewConfirm("Project Group Exists", "The project group Octoterra already exists. Do you want to delete it? It is usually safe to delete this resource.", func(b bool) {
				if b {
					if err := s.deleteProjectGroup(myclient, pggroup); err != nil {
						s.result.SetText("Failed to delete the resource")
						s.logs.SetText(err.Error())
					} else {
						s.createNewProject(parent)
					}
				}
			}, parent).Show()

			// We can't go further until the group is deleted
			return
		}

		lvsExists, lvs, err := s.libraryVariableSetExists(myclient)

		if lvsExists {
			dialog.NewConfirm("Library Variable Set Exists", "The library variable set Octoterra already exists. Do you want to delete it? It is usually safe to delete this resource.", func(b bool) {
				if b {
					if err := s.deleteLibraryVariableSet(myclient, lvs); err != nil {
						s.result.SetText("Failed to delete the resource")
						s.logs.SetText(err.Error())
					} else {
						s.createNewProject(parent)
					}
				}
			}, parent).Show()

			// We can't go further until the group is deleted
			return
		}

		// Save and apply the module
		dir, err := ioutil.TempDir("", "octoterra")
		if err != nil {
			s.logs.SetText("An error occurred while creating a temporary directory:\n" + err.Error())
			return
		}

		filePath := filepath.Join(dir, "terraform.tf")

		if err := os.WriteFile(filePath, []byte(module), 0644); err != nil {
			s.logs.SetText("An error occurred while writing the Terraform file:\n" + err.Error())
			return
		}

		initCmd := exec.Command("terraform", "init", "-no-color")
		initCmd.Dir = dir

		var initStdout, initStderr bytes.Buffer
		initCmd.Stdout = &initStdout
		initCmd.Stderr = &initStderr

		if err := initCmd.Run(); err != nil {
			s.result.SetText("Terraform init failed.")
			s.logs.SetText(initStdout.String() + initCmd.String())
			return
		}

		applyCmd := exec.Command("terraform",
			"apply",
			"-auto-approve",
			"-no-color",
			"-var=octopus_server="+s.State.Server,
			"-var=octopus_apikey="+s.State.ApiKey,
			"-var=octopus_space_id="+s.State.Space,
			"-var=octopus_destination_server="+s.State.DestinationServer,
			"-var=octopus_destination_apikey="+s.State.DestinationApiKey,
			"-var=octopus_destination_space_id="+s.State.DestinationSpace)
		applyCmd.Dir = dir

		var stdout, stderr bytes.Buffer
		applyCmd.Stdout = &stdout
		applyCmd.Stderr = &stderr

		if err := applyCmd.Run(); err != nil {
			s.result.SetText("Terraform apply failed")
			s.logs.SetText(stdout.String() + stderr.String())
			return
		} else {
			s.result.SetText("Terraform apply succeeded")
			s.logs.SetText(stdout.String() + stderr.String())
		}

		s.next.Enable()
	}()

}

func (s SpaceExportStep) deleteProjectGroup(myclient *client.Client, projectGroup *projectgroups.ProjectGroup) error {
	if err := myclient.ProjectGroups.DeleteByID(projectGroup.ID); err != nil {
		return err
	}

	return nil
}

func (s SpaceExportStep) deleteProject(myclient *client.Client, project *projects.Project) error {
	if err := myclient.Projects.DeleteByID(project.ID); err != nil {
		return err
	}

	return nil
}

func (s SpaceExportStep) projectExists(myclient *client.Client) (bool, *projects.Project, error) {
	if project, err := projects.GetByName(myclient, myclient.GetSpaceID(), spaceManagementProject); err == nil {
		return true, project, nil
	} else {
		return false, nil, err
	}
}

func (s SpaceExportStep) projectGroupExists(myclient *client.Client) (bool, *projectgroups.ProjectGroup, error) {
	if projectGroups, err := projectgroups.GetAll(myclient, myclient.GetSpaceID()); err == nil {
		groups := lo.Filter(projectGroups, func(pg *projectgroups.ProjectGroup, index int) bool {
			return pg.Name == "Octoterra"
		})

		if len(groups) == 0 {
			return false, nil, nil
		}

		return true, groups[0], nil
	} else {
		return false, nil, err
	}
}

func (s SpaceExportStep) libraryVariableSetExists(myclient *client.Client) (bool, *variables.LibraryVariableSet, error) {
	if resource, err := myclient.LibraryVariableSets.GetByPartialName("Octoterra"); err == nil {
		exatchMatch := lo.Filter(resource, func(item *variables.LibraryVariableSet, index int) bool {
			return item.Name == "Octoterra"
		})

		if len(exatchMatch) == 0 {
			return false, nil, nil
		}

		return true, exatchMatch[0], nil
	} else {
		return false, nil, err
	}
}

func (s SpaceExportStep) deleteLibraryVariableSet(myclient *client.Client, lvs *variables.LibraryVariableSet) error {
	if err := myclient.LibraryVariableSets.DeleteByID(lvs.ID); err != nil {
		return err
	}

	return nil
}
