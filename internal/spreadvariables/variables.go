package spreadvariables

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"github.com/samber/lo"
	"slices"
	"strings"
)

type OwnerVariablePair struct {
	OwnerID     string
	VariableSet *variables.VariableSet
}

func findSecretVariablesWithSharedName(variableSet *variables.VariableSet) ([]string, error) {
	groupedVariables := []string{}
	for _, variable := range variableSet.Variables {

		// The variable has to be sensitive
		if !variable.IsSensitive {
			continue
		}

		if variable.Type != "Sensitive" {
			continue
		}

		// If it has an empty scope but no other variables share the name we can leave this variable unprocessed
		if variable.Scope.IsEmpty() && !lo.ContainsBy(variableSet.Variables, func(item *variables.Variable) bool {
			return item.Name == variable.Name && item.ID != variable.ID
		}) {
			continue
		}

		if slices.Index(groupedVariables, variable.Name) == -1 {
			groupedVariables = append(groupedVariables, variable.Name)
		}
	}

	return groupedVariables, nil
}

func buildUniqueVariableName(variable *variables.Variable, usedNamed []string) string {
	name := variable.Name

	if variable.Scope.IsEmpty() {
		name += "_Unscoped"
	}

	if len(variable.Scope.Environments) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.Environments, "_"))
	}

	if len(variable.Scope.Machines) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.Machines, "_"))
	}

	if len(variable.Scope.Roles) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.Roles, "_"))
	}

	if len(variable.Scope.Actions) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.Actions, "_"))
	}

	if len(variable.Scope.TenantTags) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.TenantTags, "_"))
	}

	if len(variable.Scope.Channels) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.Channels, "_"))
	}

	if len(variable.Scope.ProcessOwners) > 0 {
		name += fmt.Sprintf("_%s", strings.Join(variable.Scope.ProcessOwners, "_"))
	}

	startingName := name
	index := 1
	for slices.Index(usedNamed, name) != -1 {
		name = startingName + "_" + fmt.Sprint(index)
		index++
	}

	return name
}

func spreadVariables(client *client.Client, ownerId string, variableSet *variables.VariableSet) error {
	groupedVariables, err := findSecretVariablesWithSharedName(variableSet)

	if err != nil {
		return err
	}

	// Get a list of all the existing variable names. We can't reuse any of these names.
	usedNames := lo.Uniq(lo.Map(variableSet.Variables, func(item *variables.Variable, index int) string {
		return item.Name
	}))

	for _, groupedVariable := range groupedVariables {
		for _, variable := range variableSet.Variables {
			if groupedVariable != variable.Name {
				continue
			}

			// Copy the original variable
			originalVar := *variable

			// Get a unique name
			uniqueName := buildUniqueVariableName(variable, usedNames)

			// Create a new variable with the original name and scopes referencing the new unscoped variable
			referenceVar := originalVar

			jsonData, err := json.Marshal(referenceVar.Scope)
			if err != nil {
				return err
			}

			// Note the original scope of this variable
			referenceVar.Description += "\n\nReplaced variable ID\n\n" + referenceVar.ID
			referenceVar.Description += "\n\nOriginal Scope\n\n" + string(jsonData)

			referenceVar.IsSensitive = false
			referenceVar.Type = "String"
			referenceVar.ID = ""
			reference := "#{" + uniqueName + "}"
			referenceVar.Value = &reference

			fmt.Println("Recreating " + referenceVar.Name + " referencing " + reference)

			_, err = variables.AddSingle(client, client.GetSpaceID(), ownerId, &referenceVar)

			if err != nil {
				return err
			}

			// Update the original variable with the new name and no scopes
			originalName := variable.Name
			usedNames = append(usedNames, uniqueName)

			if variable.Value != nil {
				panic("The value of the variable must be nil here, otherwise we may be overriding sensitive values")
			}

			fmt.Println("Renaming " + originalName + " to " + uniqueName + " and removing scopes for " + ownerId)

			jsonData, err = json.Marshal(variable.Scope)
			if err != nil {
				return err
			}

			// Note the original scope of this variable
			referenceVar.Description += "\n\nOriginal Name\n\n" + variable.Name
			// Note the original scope of this variable
			referenceVar.Description += "\n\nOriginal Scope\n\n" + string(jsonData)

			variable.Name = uniqueName
			variable.Scope = variables.VariableScope{}

			_, err = variables.UpdateSingle(client, client.GetSpaceID(), ownerId, variable)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SpreadAllVariables(state state.State) error {
	myclient, err := octoclient.CreateClient(state)

	if err != nil {
		return err
	}

	libraryVariableSets, err := myclient.LibraryVariableSets.GetAll()

	if err != nil {
		return err
	}

	projects, err := myclient.Projects.GetAll()

	if err != nil {
		return err
	}

	variableSets := []OwnerVariablePair{}

	for _, libraryVariableSet := range libraryVariableSets {
		variableSet, err := variables.GetVariableSet(myclient, myclient.GetSpaceID(), libraryVariableSet.VariableSetID)

		if err != nil {
			return errors.New("Failed to get variable set for library variable set " + libraryVariableSet.Name + ". Error was " + err.Error())
		}

		variableSets = append(variableSets, OwnerVariablePair{
			OwnerID:     libraryVariableSet.ID,
			VariableSet: variableSet,
		})
	}

	for _, project := range projects {
		variableSet, err := variables.GetVariableSet(myclient, myclient.GetSpaceID(), project.VariableSetID)

		if err != nil {
			return errors.New("Failed to get variable set for project " + project.Name + ". Error was " + err.Error())
		}

		variableSets = append(variableSets, OwnerVariablePair{
			OwnerID:     project.ID,
			VariableSet: variableSet,
		})
	}

	for _, variableSet := range variableSets {

		err = spreadVariables(myclient, variableSet.OwnerID, variableSet.VariableSet)

		if err != nil {
			return err
		}
	}

	return nil
}
