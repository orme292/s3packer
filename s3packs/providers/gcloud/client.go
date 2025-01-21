package gcloud

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GoogleClient struct {
	Ctx     context.Context
	Storage *storage.Client
	Bucket  *storage.BucketHandle
	cfg     *googleCfg
}

type googleCfg struct {
	adc string
}

func (client *GoogleClient) getClient() error {

	client.Ctx = context.Background()

	c, err := storage.NewClient(client.Ctx, option.WithCredentialsFile(client.cfg.adc))
	if err != nil {
		return fmt.Errorf("error configuring storage client: %s", err.Error())
	}

	client.Storage = c

	return nil

}

func (client *GoogleClient) getBucket(name string) {

	b := client.Storage.Bucket(name)
	client.Bucket = b

}

func (client *GoogleClient) refreshBucket() {
	client.getBucket(client.Bucket.BucketName())
}
