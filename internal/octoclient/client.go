package octoclient

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"net/url"
)

func CreateClient(state state.State) (*client.Client, error) {
	apiURL, err := url.Parse(state.Server)
	if err != nil {
		_ = fmt.Errorf("error parsing URL for Octopus API: %v", err)
		return nil, err
	}

	// the first parameter for NewClient can accept a http.Client if you wish to
	// override the default; also, the spaceID may be an empty string (i.e. "") if
	// you wish to load the default space
	octopusClient, err := client.NewClient(nil, apiURL, state.ApiKey, state.Space)
	if err != nil {
		_ = fmt.Errorf("error creating API client: %v", err)
		return nil, err
	}

	return octopusClient, nil
}
