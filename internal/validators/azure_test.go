package validators

import (
	"os"
	"testing"
)

func TestAzureContainerExists(t *testing.T) {
	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	accountName := os.Getenv("OCTOTERRAWIZ_AZURE_STORAGE_ACCOUNT")
	containerName := os.Getenv("OCTOTERRAWIZ_AZURE_CONTAINER")

	exists, err := AzureContainerExists(tenantID, clientID, clientSecret, accountName, containerName)
	if err != nil {
		t.Fatalf("Error checking if container exists: %v", err)
	}

	if !exists {
		t.Errorf("Expected container to exist, but it does not.")
	}
}
