package query

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/samber/lo"
)

func LibraryVariableSetExists(myclient *client.Client) (bool, *variables.LibraryVariableSet, error) {
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
