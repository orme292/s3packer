package provider_linode

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type LinodeOperator struct {
	App    *conf.AppConfig
	Linode *LinodeClient
}

func (oper *LinodeOperator) BucketCreate() error {

	input := &s3.CreateBucketInput{
		Bucket: aws.String(oper.App.Bucket.Name),
		ACL:    oper.App.Provider.Linode.BucketACL,
		// ACL:    types.BucketCannedACL(oper.App.Provider.AWS.ACL),
	}

	_, err := oper.Linode.s3.CreateBucket(
		context.Background(), input)
	if err != nil {
		return fmt.Errorf("error while creating bucket: %s", err.Error())
	}

	return nil

}

func (oper *LinodeOperator) BucketExists() (bool, error) {

	input := &s3.HeadBucketInput{
		Bucket: &oper.App.Bucket.Name,
	}

	_, err := oper.Linode.s3.HeadBucket(context.Background(), input)
	if err != nil {
		if errors.As(err, &s3Error) {
			if errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchBucket) {
				return false, fmt.Errorf("bucket not found")
			}
		}
		return false, fmt.Errorf("error trying to find bucket: %s", err.Error())
	}

	return true, nil

}

func (oper *LinodeOperator) BucketDelete() error {
	return nil
}

func (oper *LinodeOperator) ObjectDelete(key string) error {
	return nil
}

func (oper *LinodeOperator) ObjectExists(obj provider_v2.Object) (bool, error) {

	linObj, ok := obj.(*LinodeObject)
	if !ok {
		return true, fmt.Errorf("trouble building object to check")
	}

	input := &s3.HeadObjectInput{
		Bucket: &linObj.bucket,
		Key:    &linObj.key,
	}

	_, err := oper.Linode.s3.HeadObject(context.Background(), input)
	if err != nil {
		if errors.As(err, &s3Error) {
			if errors.As(s3Error, &s3NoSuchKey) || errors.As(s3Error, &s3NotFound) {
				return false, fmt.Errorf("object not found")
			}
		}
		return true, fmt.Errorf("error trying to find object: %s", err.Error())
	}

	return true, nil

}

func (oper *LinodeOperator) ObjectUpload(obj provider_v2.Object) error {

	linObj, ok := obj.(*LinodeObject)
	if !ok {
		return fmt.Errorf("trouble building object to upload")
	}

	input := &s3.PutObjectInput{
		Body:   linObj.f,
		Bucket: &linObj.bucket,
		Key:    &linObj.key,
	}

	if linObj.job.Metadata.HasChanged() {
		return fmt.Errorf("file changed during upload: %s", linObj.job.Metadata.FullPath())
	}
	_, err := oper.Linode.manager.Upload(context.Background(), input)
	if err != nil {

		return fmt.Errorf("error uploading [%s]: %s", err.Error(), linObj.key)
	}

	return nil

}

func (oper *LinodeOperator) GetObjectTags(key string) (map[string]string, error) {
	return nil, nil
}

func (oper *LinodeOperator) Support() *provider_v2.Supports {
	return provider_v2.NewSupports(true, false, false, false)
}

func NewLinodeOperator(app *conf.AppConfig) (oper provider_v2.Operator, err error) {

	client := LinodeClient{
		details: &details{
			key:      app.Provider.Linode.Key,
			secret:   app.Provider.Linode.Secret,
			region:   app.Bucket.Region,
			endpoint: app.Provider.Linode.Endpoint,
		},
	}

	err = client.init()
	if err != nil {
		return nil, err
	}
	if client.s3 == nil || client.manager == nil {
		return nil, errors.New("could not initialize client. check credentials")
	}

	oper = &LinodeOperator{
		App:    app,
		Linode: &client,
	}

	return oper, nil
}
