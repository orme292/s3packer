package pack_akamai

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/orme292/s3packer/conf"
)

func buildUploader(ac *conf.AppConfig) (uploader *manager.Uploader, client *s3.Client, err error) {
	client, err = buildClient(ac)
	uploader = manager.NewUploader(client, func(u *manager.Uploader) {
		u.MaxUploadParts = int32(ac.Opts.MaxParts)
		u.LeavePartsOnError = false
	})
	return
}

func buildClient(ac *conf.AppConfig) (client *s3.Client, err error) {
	creds := credentials.NewStaticCredentialsProvider(ac.Provider.Linode.Key,
		ac.Provider.Linode.Secret, "")
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(ac.Bucket.Region),
	)
	if err != nil {
		return nil, err
	}

	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(s("https://%s", ac.Provider.Linode.Endpoint))
	})

	return client, nil
}

func s(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
