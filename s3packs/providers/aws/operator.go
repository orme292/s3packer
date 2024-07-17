package provider_aws

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type AwsOperator struct {
	App *conf.AppConfig
	AWS *AwsClient
}

func NewAwsOperator(app *conf.AppConfig) (oper provider_v2.Operator, err error) {

	client := AwsClient{
		details: &details{
			profile: app.Provider.AWS.Profile,
			key:     app.Provider.AWS.Key,
			secret:  app.Provider.AWS.Secret,
			region:  app.Bucket.Region,
		},
	}

	err = client.getClient()
	if err != nil {
		return nil, err
	}
	if client.s3 == nil || client.manager == nil {
		return nil, errors.New("could not initialize AWS client. check your credentials")
	}

	oper = &AwsOperator{
		App: app,
		AWS: &client,
	}

	return oper, nil

}

func (oper *AwsOperator) BucketCreate() error {

	input := &s3.CreateBucketInput{
		Bucket: aws.String(oper.App.Bucket.Name),
		ACL:    types.BucketCannedACL(oper.App.Provider.AWS.ACL),
	}

	location := types.BucketLocationConstraint(oper.App.Bucket.Region)
	if location != "us-east-1" {
		input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
			LocationConstraint: location,
		}
	}

	_, err := oper.AWS.s3.CreateBucket(
		context.Background(), input)
	if err != nil {
		return fmt.Errorf("error while creating bucket: %s", err.Error())
	}

	return nil

}

func (oper *AwsOperator) BucketExists() (bool, error) {

	input := &s3.HeadBucketInput{
		Bucket: &oper.App.Bucket.Name,
	}

	_, err := oper.AWS.s3.HeadBucket(context.Background(), input)
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

func (oper *AwsOperator) BucketDelete() error {

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(oper.App.Bucket.Name),
	}

	_, err := oper.AWS.s3.DeleteBucket(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error deleting bucket: %s", err.Error())
	}

	return nil

}

func (oper *AwsOperator) ObjectDelete(key string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(oper.App.Bucket.Name),
		Key:    aws.String(key),
	}
	_, err := oper.AWS.s3.DeleteObject(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error deleting object: %s", err.Error())
	}

	return nil

}

func (oper *AwsOperator) ObjectExists(obj provider_v2.Object) (bool, error) {

	awsObj, ok := obj.(*AwsObject)
	if !ok {
		return true, fmt.Errorf("trouble building object to check")
	}

	input := &s3.HeadObjectInput{
		Bucket: &awsObj.bucket,
		Key:    &awsObj.key,
	}

	_, err := oper.AWS.s3.HeadObject(context.Background(), input)
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

func (oper *AwsOperator) ObjectUpload(obj provider_v2.Object) error {

	awsObj, ok := obj.(*AwsObject)
	if !ok {
		return fmt.Errorf("trouble building object to upload")
	}

	if awsObj.job.Metadata.HasChanged() {
		return fmt.Errorf("file changed during upload: %s", awsObj.job.Metadata.FullPath())
	}

	input := &s3.PutObjectInput{
		ACL:               awsObj.acl,
		Body:              awsObj.f,
		Bucket:            &awsObj.bucket,
		Key:               &awsObj.key,
		ChecksumAlgorithm: oper.App.Provider.AWS.AwsChecksumAlgorithm,
		StorageClass:      awsObj.storage,
		Tagging:           aws.String(awsObj.tags),
	}

	_, err := oper.AWS.manager.Upload(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error uploading [%s]: %s", err.Error(), awsObj.key)
	}

	return nil

}

func (oper *AwsOperator) GetObjectTags(key string) (map[string]string, error) {

	input := &s3.GetObjectTaggingInput{
		Bucket: &oper.App.Bucket.Name,
		Key:    &key,
	}

	_, err := oper.AWS.s3.GetObjectTagging(context.Background(), input)
	if err != nil {
		return make(map[string]string), fmt.Errorf("error getting object tags: %s", err.Error())
	}

	return make(map[string]string), nil

}

func (oper *AwsOperator) Support() *provider_v2.Supports {

	return provider_v2.NewSupports(true, true, true, false)

}
