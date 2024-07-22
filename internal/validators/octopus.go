package validators

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/state"
)

func ValidateSourceCreds(state state.State) bool {
	if myclient, err := octoclient.CreateClient(state); err != nil {
		return false
	} else {
		if _, err := spaces.GetByID(myclient, myclient.GetSpaceID()); err != nil {
			return false
		}
	}

	return true
}

func ValidateDestinationCreds(state state.State) bool {
	if myclient, err := octoclient.CreateDestinationClient(state); err != nil {
		return false
	} else {
		if _, err := spaces.GetByID(myclient, myclient.GetSpaceID()); err != nil {
			return false
		}
	}

	return true
}
