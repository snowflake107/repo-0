package validators

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"time"
)

type CustomCredentials struct {
	State state.State
}

func (c CustomCredentials) Retrieve(ctx context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID:     c.State.AwsAccessKey,
		SecretAccessKey: c.State.AwsSecretKey,
		SessionToken:    "",
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
		AccountID:       "",
	}, nil
}

func ValidateAWS(state state.State) bool {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(CustomCredentials{state}),
		config.WithRegion(state.AwsS3BucketRegion))
	if err != nil {
		return false
	}

	simpleTokenService := sts.NewFromConfig(cfg)

	_, err = simpleTokenService.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return false
	}

	return true
}
