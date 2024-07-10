package handler_aws

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type AwsObject struct {
	job *provider_v2.Job

	f *os.File

	key string

	acl     types.ObjectCannedACL
	storage types.StorageClass

	bucket string

	tags string
}

func NewAwsObject(job *provider_v2.Job) provider_v2.Object {
	return &AwsObject{
		job: job,
	}
}

func (o *AwsObject) Destroy() error { return o.Post() }

func (o *AwsObject) Generate() error {

	o.acl = o.job.App.Provider.AWS.ACL
	o.storage = o.job.App.Provider.AWS.Storage

	o.bucket = o.job.App.Bucket.Name

	o.key = o.job.Key

	o.setTags(o.job.App.Tags)

	return nil

}

func (o *AwsObject) Post() error { return o.f.Close() }

func (o *AwsObject) Pre() error {

	o.job.Metadata.Update()

	if !o.job.Metadata.IsExists || !o.job.Metadata.IsReadable {
		return fmt.Errorf("file no longer accessible")
	}

	f, err := os.Open(o.job.Metadata.FullPath())
	if err != nil {
		fmt.Printf("Error opening file %s: %s\n", o.job.Metadata.FullPath(), err)
		return err
	}

	o.f = f

	return nil

}

func (o *AwsObject) setTags(input map[string]string) {

	if len(input) == 0 {
		o.tags = EmptyString
		return
	}

	var tagString string
	for k, v := range input {
		if tagString == EmptyString {
			tagString = fmt.Sprintf("%s=%s", k, v)
		} else {
			tagString = fmt.Sprintf("%s&%s=%s", tagString, k, v)
		}
	}

	o.tags = tagString

}
