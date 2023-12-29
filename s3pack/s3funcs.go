// Package s3pack provides functions for uploading files to s3.
// This file implements functions for creating a session.Session object and an s3manager.Uploader object.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/orme292/s3packer/conf"
)

/*
BuildUploader builds and returned a manager.Uploader object. It takes a config.Configuration object. The func creates
a session by calling NewConfig and passes the aws config to manager.NewUploader.
*/
func BuildUploader(a *conf.AppConfig) (uploader *manager.Uploader, err error) {
	client, err := BuildClient(a)
	uploader = manager.NewUploader(client, func(u *manager.Uploader) {
		u.MaxUploadParts = int32(a.Opts.MaxParts)
	})
	return
}

func BuildClient(a *conf.AppConfig) (client *s3.Client, err error) {
	cfg, err := NewConfig(a)
	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = a.Bucket.Region
	})
	return
}

/*
NewConfig creates a new aws.Config object using the authentication information in the profile. It takes a
config.Configuration object. It determines whether to use a profile or a keypair based on the presence of a profile
name in the profile configuration.

NewConfig calls NewConfigWithProfile or NewConfigWithKeypair
*/
func NewConfig(a *conf.AppConfig) (cfg aws.Config, err error) {
	if a.Provider.AwsProfile != EmptyString {
		cfg, err = NewConfigWithProfile(a)
	} else {
		cfg, err = NewConfigWithKeypair(a)
	}
	if err != nil {
		a.Log.Fatal("Unable to create build config: %q", err.Error())
	}

	return
}

func NewConfigWithKeypair(a *conf.AppConfig) (cfg aws.Config, err error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		a.Provider.Key, a.Provider.Secret, ""))
	opts := func(o *config.LoadOptions) error {
		o.Region = a.Bucket.Region
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(creds), opts)
	return
}

func NewConfigWithProfile(a *conf.AppConfig) (cfg aws.Config, err error) {
	opts := func(o *config.LoadOptions) error {
		o.Region = a.Bucket.Region
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(a.Provider.AwsProfile),
		opts)
	return
}
