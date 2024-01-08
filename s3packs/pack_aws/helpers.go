package pack_aws

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

/*
BuildUploader builds and returned a manager.Uploader object. It takes a config.Configuration object. The func creates
a session by calling NewConfig and passes the aws config to manager.NewUploader.
*/
func buildUploader(ac *conf.AppConfig) (uploader *manager.Uploader, client *s3.Client, err error) {
	client, err = buildClient(ac)
	uploader = manager.NewUploader(client, func(u *manager.Uploader) {
		u.MaxUploadParts = int32(ac.Opts.MaxParts)
	})
	return
}

func buildClient(ac *conf.AppConfig) (client *s3.Client, err error) {
	cfg, err := newConfig(ac)
	client = s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = ac.Bucket.Region
	})
	return
}

/*
NewConfig creates a new aws.Config object using the authentication information in the profile. It takes a
config.Configuration object. It determines whether to use a profile or a keypair based on the presence of a profile
name in the profile configuration.

NewConfig calls NewConfigWithProfile or NewConfigWithKeypair
*/
func newConfig(ac *conf.AppConfig) (cfg aws.Config, err error) {
	if ac.Provider.AwsProfile != EmptyString {
		cfg, err = newConfigWithProfile(ac)
	} else {
		cfg, err = newConfigWithKeypair(ac)
	}
	if err != nil {
		ac.Log.Fatal("Unable to create build config: %q", err.Error())
	}

	return
}

func newConfigWithKeypair(ac *conf.AppConfig) (cfg aws.Config, err error) {
	creds := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		ac.Provider.Key, ac.Provider.Secret, ""))
	opts := func(o *config.LoadOptions) error {
		o.Region = ac.Bucket.Region
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.Background(), config.WithCredentialsProvider(creds), opts)
	return
}

func newConfigWithProfile(ac *conf.AppConfig) (cfg aws.Config, err error) {
	opts := func(o *config.LoadOptions) error {
		o.Region = ac.Bucket.Region
		return nil
	}
	cfg, err = config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(ac.Provider.AwsProfile),
		opts)
	return
}

func s(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func awsTag(tags map[string]string) string {
	if len(tags) == 0 {
		return EmptyString
	}
	var tag string
	for k, v := range tags {
		if tag == EmptyString {
			tag = s("%s=%s", k, v)
		} else {
			tag = s("%s&%s=%s", tag, k, v)
		}
	}
	return tag
}
