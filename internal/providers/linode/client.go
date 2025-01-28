package linode

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type LinodeClient struct {
	details *details
	s3      *s3.Client

	manager *manager.Uploader
}

type details struct {
	key      string
	secret   string
	region   string
	endpoint string
}

func (client *LinodeClient) init() error {

	creds := credentials.NewStaticCredentialsProvider(
		client.details.key, client.details.secret, "")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(client.details.region))
	if err != nil {
		return err
	}

	cl := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s", client.details.endpoint))
	})

	client.s3 = cl

	mgr := manager.NewUploader(client.s3, func(u *manager.Uploader) {
		u.LeavePartsOnError = false
	})

	client.manager = mgr
	return nil

}
