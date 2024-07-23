package state

type State struct {
	BackendType             string
	Server                  string
	ApiKey                  string
	Space                   string
	DestinationServer       string
	DestinationApiKey       string
	DestinationSpace        string
	AwsAccessKey            string
	AwsSecretKey            string
	AwsS3Bucket             string
	AwsS3BucketRegion       string
	PromptForDelete         bool
	AzureResourceGroupName  string
	AzureStorageAccountName string
	AzureContainerName      string
	AzureKeyName            string
}
