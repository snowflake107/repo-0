package validators

import (
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"os"
	"testing"
)

func TestAWSCredsExists(t *testing.T) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	testState := state.State{
		BackendType:               "",
		Server:                    "",
		ServerExternal:            "",
		ApiKey:                    "",
		Space:                     "",
		DestinationServer:         "",
		DestinationServerExternal: "",
		DestinationApiKey:         "",
		DestinationSpace:          "",
		AwsAccessKey:              accessKey,
		AwsSecretKey:              secretKey,
		AwsS3Bucket:               "",
		AwsS3BucketRegion:         "",
		PromptForDelete:           false,
		UseContainerImages:        false,
		AzureResourceGroupName:    "",
		AzureStorageAccountName:   "",
		AzureContainerName:        "",
		AzureSubscriptionId:       "",
		AzureTenantId:             "",
		AzureApplicationId:        "",
		AzurePassword:             "",
	}

	err := ValidateAWS(testState)
	if err != nil {
		t.Fatalf("Error checking aws creds: %v", err)
	}
}
