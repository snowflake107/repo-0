package octoclient

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"net/url"
	"strings"
)

func CreateClient(state state.State) (*client.Client, error) {
	return createClient(
		strings.TrimSpace(state.Server),
		strings.TrimSpace(state.ApiKey),
		strings.TrimSpace(state.Space))
}

func CreateDestinationClient(state state.State) (*client.Client, error) {
	return createClient(
		strings.TrimSpace(state.DestinationServer),
		strings.TrimSpace(state.DestinationApiKey),
		strings.TrimSpace(state.DestinationSpace))
}

func createClient(server string, apikey string, space string) (*client.Client, error) {
	apiURL, err := url.Parse(server)
	if err != nil {
		_ = fmt.Errorf("error parsing URL for Octopus API: %v", err)
		return nil, err
	}

	// the first parameter for NewClient can accept a http.Client if you wish to
	// override the default; also, the spaceID may be an empty string (i.e. "") if
	// you wish to load the default space
	octopusClient, err := client.NewClient(nil, apiURL, apikey, space)
	if err != nil {
		_ = fmt.Errorf("error creating API client: %v", err)
		return nil, err
	}

	return octopusClient, nil
}
