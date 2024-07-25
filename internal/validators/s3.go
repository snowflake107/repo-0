package validators

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/mcasperson/OctoterraWizard/internal/state"
	"log"
)

func TestS3Bucket(state state.State) error {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithCredentialsProvider(CustomCredentials{state}),
		config.WithRegion(state.AwsS3BucketRegion))
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg)

	if _, err := bucketExists(s3Client, state.AwsS3Bucket); err != nil {
		return err
	}

	return nil
}

func bucketExists(s3Client *s3.Client, bucketName string) (bool, error) {
	_, err := s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	} else {
		log.Printf("Bucket %v exists and you already own it.", bucketName)
	}

	return exists, err
}
