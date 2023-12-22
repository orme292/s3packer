// Package s3pack provides functions for uploading files to s3.
// This file implements functions for creating a session.Session object and an s3manager.Uploader object.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	app "github.com/orme292/s3packer/config"
)

/*
BuildUploader builds and returned a manager.Uploader object. It takes a config.Configuration object. The func creates
a session by calling NewConfig and passes the aws config to manager.NewUploader.
*/
func BuildUploader(c *app.Configuration) (uploader *manager.Uploader, err error) {
	client, err := BuildClient(c)
	uploader = manager.NewUploader(client, func(u *manager.Uploader) {
		u.MaxUploadParts = 1
	})
	return
}

func BuildClient(c *app.Configuration) (client *s3.Client, err error) {
	cfg, err := NewConfig(c)
	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = c.Bucket[app.ProfileBucketRegion].(string)
	})
	return
}

/*
NewConfig creates a new aws.Config object using the authentication information in the profile. It takes a
config.Configuration object. It determines whether to use a profile or a keypair based on the presence of a profile
name in the profile configuration.

NewConfig calls NewConfigWithProfile or NewConfigWithKeypair
*/
func NewConfig(c *app.Configuration) (cfg aws.Config, err error) {
	if c.Authentication[app.ProfileAuthProfile].(string) != EmptyString {
		cfg, err = NewConfigWithProfile(c)
	} else {
		cfg, err = NewConfigWithKeypair(c)
	}
	if err != nil {
		c.Logger.Fatal(fmt.Sprintf("Unable to create build config: %q", err.Error()))
	}
	return
}

func NewConfigWithKeypair(c *app.Configuration) (cfg aws.Config, err error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		c.Authentication[app.ProfileAuthKey].(string), c.Authentication[app.ProfileAuthSecret].(string), ""))
	opts := func(o *config.LoadOptions) error {
		o.Region = c.Bucket[app.ProfileBucketRegion].(string)
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds), opts)
	return
}

func NewConfigWithProfile(c *app.Configuration) (cfg aws.Config, err error) {
	opts := func(o *config.LoadOptions) error {
		o.Region = c.Bucket[app.ProfileBucketRegion].(string)
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(c.Authentication[app.ProfileAuthProfile].(string)),
		opts)
	return
}
