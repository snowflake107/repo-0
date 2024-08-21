package validators

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/spaces"
	"github.com/mcasperson/OctoterraWizard/internal/octoclient"
	"github.com/mcasperson/OctoterraWizard/internal/state"
)

func ValidateSourceCreds(state state.State) error {
	if myclient, err := octoclient.CreateClient(state); err != nil {
		return err
	} else {
		if _, err := spaces.GetByID(myclient, myclient.GetSpaceID()); err != nil {
			return err
		}
	}

	return nil
}

func ValidateDestinationCreds(state state.State) error {
	if myclient, err := octoclient.CreateDestinationClient(state); err != nil {
		return err
	} else {
		if _, err := spaces.GetByID(myclient, myclient.GetSpaceID()); err != nil {
			return err
		}
	}

	return nil
}
