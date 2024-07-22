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
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/query"
	"github.com/mcasperson/OctoterraWizard/internal/strutil"
	"github.com/mcasperson/OctoterraWizard/internal/wizard"
	"github.com/samber/lo"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed modules/project_management/terraform.tf
var runbookModule string

type ProjectExportStep struct {
	BaseStep
	Wizard        wizard.Wizard
	createProject *widget.Button
	infinite      *widget.ProgressBarInfinite
	result        *widget.Label
	logs          *widget.Entry
	next          *widget.Button
	previous      *widget.Button
}

func (s ProjectExportStep) GetContainer(parent fyne.Window) *fyne.Container {

	bottom, thisPrevious, thisNext := s.BuildNavigation(func() {
		s.Wizard.ShowWizardStep(SpaceExportStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	}, func() {
		s.Wizard.ShowWizardStep(FinishStep{
			Wizard:   s.Wizard,
			BaseStep: BaseStep{State: s.State}})
	})
	s.next = thisNext
	s.previous = thisPrevious
	s.next.Disable()

	intro := widget.NewLabel(strutil.TrimMultilineWhitespace(`Each project gets two runbooks, one to serialize it to a Terraform module, and the second to deploy it.`))
	s.infinite = widget.NewProgressBarInfinite()
	s.infinite.Start()
	s.infinite.Hide()
	s.result = widget.NewLabel("")
	s.logs = widget.NewEntry()
	s.logs.Disable()
	s.logs.MultiLine = true
	s.logs.SetMinRowsVisible(20)
	s.createProject = widget.NewButton("Add Runbooks", func() { s.createNewProject(parent) })
	middle := container.New(layout.NewVBoxLayout(), intro, s.createProject, s.infinite, s.result, s.logs)

	content := container.NewBorder(nil, bottom, nil, nil, middle)

	return content
}

func (s ProjectExportStep) createNewProject(parent fyne.Window) {
	s.logs.SetText("")
	s.next.Disable()
	s.previous.Disable()
	s.infinite.Show()
	s.createProject.Disable()
	s.result.SetText("🔵 Creating runbooks. This can take a little while.")

	go func() {
		defer s.previous.Enable()
		defer s.infinite.Hide()
		defer s.createProject.Enable()

		myclient, err := octoclient.CreateClient(s.State)

		if err != nil {
			s.logs.SetText("🔴 Failed to create the client:\n" + err.Error())
			return
		}

		allProjects, err := s.getProjects(myclient)

		if err != nil {
			s.logs.SetText("🔴 Failed to get all the projects:\n" + err.Error())
			return
		}

		allProjects = lo.Filter(allProjects, func(project *projects.Project, index int) bool {
			return project.Name != spaceManagementProject
		})

		lvsExists, lvs, err := query.LibraryVariableSetExists(myclient)

		if err != nil {
			s.logs.SetText("🔴 Failed to get the library variable set Octoterra:\n" + err.Error())
			return
		}

		if !lvsExists {
			s.logs.SetText("🔴 The library variable set Octoterra could not be found")
			return
		}

		// First look deletes any existing projects
		for _, project := range allProjects {
			if project.Name == spaceManagementProject {
				continue
			}

			projExists, runbook, err := s.runbookExists(myclient, project.ID, "__ 1. Serialize Project")

			if err != nil {
				s.logs.SetText("🔴 Failed to find runbook:\n" + err.Error())
				return
			}

			if projExists {
				dialog.NewConfirm("Project Group Exists", "The runbook \"__ 1. Serialize Project\" already exists in project "+project.Name+". Do you want to delete it? It is usually safe to delete this resource.", func(b bool) {
					if b {
						if err := s.deleteRunbook(myclient, runbook); err != nil {
							s.result.SetText("🔴 Failed to delete the resource")
							s.logs.SetText(err.Error())
						} else {
							s.createNewProject(parent)
						}
					}
				}, parent).Show()

				// We can't go further until the resource is deleted
				return
			}
		}

		// Find the step template ID
		serializeProjectTemplate, err, message := query.GetStepTemplateId(myclient, s.State, "Octopus - Serialize Project to Terraform")

		if err != nil {
			s.result.SetText(message)
			s.logs.SetText(err.Error())
			return
		}

		deploySpaceTemplateS3, err, message := query.GetStepTemplateId(myclient, s.State, "Octopus - Populate Octoterra Space (S3 Backend)")

		if err != nil {
			s.result.SetText(message)
			s.logs.SetText(err.Error())
			return
		}

		for _, project := range allProjects {
			// Save and apply the module
			dir, err := ioutil.TempDir("", "octoterra")
			if err != nil {
				s.logs.SetText("An error occurred while creating a temporary directory:\n" + err.Error())
				return
			}

			filePath := filepath.Join(dir, "terraform.tf")

			if err := os.WriteFile(filePath, []byte(runbookModule), 0644); err != nil {
				s.logs.SetText("🔴 An error occurred while writing the Terraform file:\n" + err.Error())
				return
			}

			initCmd := exec.Command("terraform", "init", "-no-color")
			initCmd.Dir = dir

			var initStdout, initStderr bytes.Buffer
			initCmd.Stdout = &initStdout
			initCmd.Stderr = &initStderr

			if err := initCmd.Run(); err != nil {
				s.result.SetText("🔴 Terraform init failed.")
				s.logs.SetText(initStdout.String() + initCmd.String())
				return
			}

			applyCmd := exec.Command("terraform",
				"apply",
				"-auto-approve",
				"-no-color",
				"-var=octopus_serialize_actiontemplateid="+serializeProjectTemplate,
				"-var=octopus_deploys3_actiontemplateid="+deploySpaceTemplateS3,
				"-var=octopus_server="+s.State.Server,
				"-var=octopus_apikey="+s.State.ApiKey,
				"-var=octopus_space_id="+s.State.Space,
				"-var=octopus_project_id="+project.ID,
				"-var=terraform_state_bucket="+s.State.AwsS3Bucket,
				"-var=terraform_state_bucket_region="+s.State.AwsS3BucketRegion,
				"-var=octopus_destination_server="+s.State.DestinationServer,
				"-var=octopus_destination_apikey="+s.State.DestinationApiKey,
				"-var=octopus_destination_space_id="+s.State.DestinationSpace)
			applyCmd.Dir = dir

			var stdout, stderr bytes.Buffer
			applyCmd.Stdout = &stdout
			applyCmd.Stderr = &stderr

			if err := applyCmd.Run(); err != nil {
				s.result.SetText("🔴 Terraform apply failed")
				s.logs.SetText(stdout.String() + stderr.String())
				return
			} else {
				s.result.SetText("Terraform apply succeeded")
				s.logs.SetText(stdout.String() + stderr.String())
			}

			// link the library variable set
			project, err := myclient.Projects.GetByID(project.ID)

			if err != nil {
				s.logs.SetText("🔴 Failed to get the project:\n" + err.Error())
				return
			}

			project.IncludedLibraryVariableSets = append(project.IncludedLibraryVariableSets, lvs.ID)

			_, err = projects.Update(myclient, project)

			if err != nil {
				s.logs.SetText("🔴 Failed to update the project:\n" + err.Error())
				return
			}

		}

		s.result.SetText("🟢 Added runbooks to all projects")
		s.next.Enable()
	}()
}

func (s ProjectExportStep) getProjects(myclient *client.Client) ([]*projects.Project, error) {
	if allprojects, err := myclient.Projects.GetAll(); err != nil {
		return nil, err
	} else {
		return allprojects, nil
	}
}

func (s ProjectExportStep) deleteRunbook(myclient *client.Client, runbook *runbooks.Runbook) error {
	if err := myclient.Runbooks.DeleteByID(runbook.ID); err != nil {
		return err
	}

	return nil
}

func (s ProjectExportStep) runbookExists(myclient *client.Client, projectId string, runbookName string) (bool, *runbooks.Runbook, error) {
	if runbook, err := runbooks.GetByName(myclient, myclient.GetSpaceID(), projectId, runbookName); err == nil {
		if runbook == nil {
			return false, nil, nil
		}
		return true, runbook, nil
	} else {
		return false, nil, err
	}
}
