package provider_aws

import (
	"fmt"
	"os"
	"regexp"

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
	cleanString := func(s string) string {
		reg, err := regexp.Compile("[^a-zA-Z0-9_\\.\\/\\=\\+\\-\\:\\@\\s]+")
		if err != nil {
			return ""
		}
		return reg.ReplaceAllString(s, "_")
	}

	buildTags := func(slc map[string]string, tags string) string {

		if len(slc) == 0 {
			return tags
		}

		for k, v := range slc {
			k = cleanString(k)
			v = cleanString(v)
			if tags == EmptyString {
				tags = fmt.Sprintf("%s=%s", k, v)
			} else {
				tags = fmt.Sprintf("%s&%s=%s", tags, k, v)
			}
		}

		return tags

	}

	o.tags = buildTags(o.job.AppTags, "")
	o.tags = buildTags(input, o.tags)

}
