package validators

import (
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"os"
	"testing"
)

func TestS3BucketExists(t *testing.T) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_DEFAULT_REGION")
	bucket := os.Getenv("AWS_DEFAULT_BUCKET")

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
		AwsS3Bucket:               bucket,
		AwsS3BucketRegion:         region,
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

	err := TestS3Bucket(testState)
	if err != nil {
		t.Fatalf("Error checking if bucket exists: %v", err)
	}
}
