package handler_aws

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type AwsObject struct {
	acl     *types.ObjectCannedACL
	storage *types.StorageClass

	bucket *string
	key    *string

	f *os.File

	tags *string
}

func NewAwsObject(job *provider_v2.QueueJob) provider_v2.Object {

	awsObj := &AwsObject{}
	err := awsObj.Generate(job)
	if err != nil {
		return nil
	}

	return awsObj

}

func (o *AwsObject) Destroy() error {

	err := o.f.Close()
	if err != nil {
		log.Printf("Error closing file (object): %s", err)
	}

	return nil
}

func (o *AwsObject) Generate(job *provider_v2.QueueJob) error {

	err := job.OpenFile()
	if err != nil {
		return err
	}
	o.f = job.F

	o.acl = &job.App.Provider.AWS.ACL
	o.storage = &job.App.Provider.AWS.Storage

	o.bucket = aws.String(job.App.Bucket.Name)
	o.key = aws.String(job.Key)

	o.setTags(job.App.Tags)

	return nil

}

func (o *AwsObject) setTags(input map[string]string) {

	if len(input) == 0 {
		o.tags = aws.String(EmptyString)
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

	o.tags = aws.String(tagString)

}
