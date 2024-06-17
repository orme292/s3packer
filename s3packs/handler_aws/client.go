package handler_aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type AwsClient struct {
	details details
	config  aws.Config
	manager *manager.Uploader
	s3      *s3.Client
}

type details struct {
	profile string
	key     string
	secret  string
	region  string
	parts   int
}

func (client *AwsClient) build() error {
	return nil
}

func (client *AwsClient) getClient() error {

	err := client.getClientConfig()
	if err != nil {
		return fmt.Errorf("error configuring aws client: %s", err.Error())
	}

	if client.details.parts == 0 {
		client.details.parts = int(manager.MaxUploadParts)
	}

	client.s3 = s3.NewFromConfig(client.config)

	return nil

}

func (client *AwsClient) getClientConfig() error {

	var err error

	if client.details.profile != EmptyString {
		err = client.getClientConfigFromProfile()
	} else {
		err = client.getClientConfigFromKeys()
	}
	if err != nil {
		return err
	}

	return nil

}

func (client *AwsClient) getClientConfigFromKeys() error {

	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		client.details.key, client.details.secret, ""))

	conf, err := config.LoadDefaultConfig(
		context.Background(), config.WithCredentialsProvider(creds),
		awsLoadOpts(client.details.region))
	if err != nil {
		return err
	}

	client.config = conf

	return nil

}

func (client *AwsClient) getClientConfigFromProfile() error {

	conf, err := config.LoadDefaultConfig(
		context.Background(), config.WithSharedConfigProfile(client.details.profile),
		awsLoadOpts(client.details.region))
	if err != nil {
		return err
	}

	client.config = conf

	return nil

}

func (client *AwsClient) getUploadManager() error {

	if client.s3 == nil {
		return fmt.Errorf("AWS client not configured")
	}

	client.manager = manager.NewUploader(client.s3, func(u *manager.Uploader) {
		u.MaxUploadParts = int32(client.details.parts)
		u.LeavePartsOnError = false
	})

	return nil

}

func awsLoadOpts(region string) func(*config.LoadOptions) error {

	return func(o *config.LoadOptions) error {
		o.Region = region
		return nil
	}

}
