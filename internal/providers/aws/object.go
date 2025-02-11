package aws

import (
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"s3p/internal/provider"
)

type AwsObject struct {
	job *provider.Job

	f *os.File

	key string

	acl     types.ObjectCannedACL
	storage types.StorageClass

	bucket string

	tags string
}

func NewAwsObject(job *provider.Job) provider.Object {
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
		reg := regexp.MustCompile(`[^a-zA-Z0-9_./=\+\-:@\s]+`)
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
