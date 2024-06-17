package handler_aws

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
	app *conf.AppConfig
	aws *AwsClient
}

func NewAwsOperator(app *conf.AppConfig) (provider_v2.Operator, error) {

	client := AwsClient{
		details: details{
			profile: app.Provider.AWS.Profile,
			key:     app.Provider.AWS.Key,
			secret:  app.Provider.AWS.Secret,
			region:  app.Bucket.Region,
			parts:   app.Opts.MaxParts,
		},
	}

	err := client.build()
	if err != nil {
		return nil, err
	}

	oper := &AwsOperator{
		app: app,
		aws: &client,
	}

	return oper, nil

}

func (oper *AwsOperator) BucketCreate() error {

	location := types.BucketLocationConstraint(oper.app.Bucket.Region)
	bucketConf := &types.CreateBucketConfiguration{
		LocationConstraint: location,
	}

	input := &s3.CreateBucketInput{
		Bucket:                    aws.String(oper.app.Bucket.Name),
		ACL:                       types.BucketCannedACL(oper.app.Provider.AWS.ACL),
		CreateBucketConfiguration: bucketConf,
	}

	output, err := oper.aws.s3.CreateBucket(
		context.Background(), input)
	if err != nil {
		return fmt.Errorf("error while creating bucket: %s", err.Error())
	}

	fmt.Printf("oper.BucketCreate: RESULT METADATA: %+v\n", output)

	return nil

}

func (oper *AwsOperator) BucketExists() (bool, error) {

	input := &s3.HeadBucketInput{
		Bucket: aws.String(oper.app.Bucket.Name),
	}

	output, err := oper.aws.s3.HeadBucket(context.Background(), input)
	if err != nil {
		if errors.As(err, &s3Error) {
			if errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchBucket) {
				return false, fmt.Errorf("bucket not found")
			}
		}
		return false, fmt.Errorf("error trying to find bucket: %s", err.Error())
	}

	fmt.Printf("oper.BucketExists: RESULT METADATA: %+v\n", output)

	return true, nil

}

func (oper *AwsOperator) BucketDelete() error {

	input := &s3.DeleteBucketInput{
		Bucket: aws.String(oper.app.Bucket.Name),
	}

	output, err := oper.aws.s3.DeleteBucket(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error deleting bucket: %s", err.Error())
	}

	fmt.Printf("oper.BucketDelete: RESULT METADATA: %+v\n", output)

	return nil

}

func (oper *AwsOperator) ObjectExists(key string) (bool, error) {

	input := &s3.HeadObjectInput{
		Bucket: aws.String(oper.app.Bucket.Name),
		Key:    aws.String(key),
	}

	output, err := oper.aws.s3.HeadObject(context.Background(), input)
	if err != nil {
		if errors.As(err, &s3Error) {
			if errors.As(s3Error, &s3NoSuchKey) || errors.As(s3Error, &s3NotFound) {
				return false, fmt.Errorf("object not found")
			}
		}
		return false, fmt.Errorf("error trying to find object: %s", err.Error())
	}

	fmt.Printf("oper.ObjectExists(%s): RESULT METADATA: %+v\n", key, output)

	return true, nil

}

func (oper *AwsOperator) ObjectDelete(key string) error {

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(oper.app.Bucket.Name),
		Key:    aws.String(key),
	}
	output, err := oper.aws.s3.DeleteObject(context.Background(), input)
	if err != nil {
		return fmt.Errorf("error deleting object: %s", err.Error())
	}

	fmt.Printf("oper.ObjectDelete(%s): RESULT METADATA: %+v\n", key, output)

	return nil

}

func (oper *AwsOperator) ObjectUpload() error {

	return nil

}

func (oper *AwsOperator) ObjectUploadMultipart() error {
	return nil
}

func (oper *AwsOperator) GetObjectTags(key string) (map[string]string, error) {

	input := &s3.GetObjectTaggingInput{
		Bucket: &oper.app.Bucket.Name,
		Key:    &key,
	}

	output, err := oper.aws.s3.GetObjectTagging(context.Background(), input)
	if err != nil {
		return make(map[string]string), fmt.Errorf("error getting object tags: %s", err.Error())
	}

	fmt.Printf("oper.ObjectDelete(%s): RESULT METADATA: %+v\n", key, output)

	return make(map[string]string), nil

}

func (oper *AwsOperator) Support() *provider_v2.Supports {

	return &provider_v2.Supports{
		BucketCreate:          true,
		BucketDelete:          false,
		ObjectDelete:          false,
		ObjectUploadMultipart: false,
	}

}
