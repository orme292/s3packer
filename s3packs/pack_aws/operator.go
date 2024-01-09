package pack_aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider"
)

func NewAwsOperator(ac *conf.AppConfig) (*AwsOperator, error) {
	svc, client, err := buildUploader(ac)
	if err != nil {
		return nil, err
	}
	return &AwsOperator{
		ac:     ac,
		client: client,
		svc:    svc,
	}, nil
}

func (op *AwsOperator) CreateBucket() (err error) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(op.ac.Bucket.Name),
		ACL:    types.BucketCannedACL(op.ac.Provider.AwsACL),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(op.ac.Bucket.Region),
		},
	}
	_, err = op.client.CreateBucket(context.Background(), input)
	if err != nil {
		op.ac.Log.Error("Unable to create bucket %q: %q", op.ac.Bucket.Name, err.Error())
		return err
	}
	op.ac.Log.Info("Created bucket %q in %q", op.ac.Bucket.Name, op.ac.Bucket.Region)
	return
}

func (op *AwsOperator) Get(key string) (obj *provider.GetObject, err error) {
	return nil, errors.New(ErrorNotImplemented)
}

func (op *AwsOperator) ObjectExists(key string) (exists bool, err error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(op.ac.Bucket.Name),
		Key:    aws.String(key),
	}
	_, err = op.client.HeadObject(context.Background(), input)
	if err != nil {
		if errors.As(err, &s3Error) {
			switch {
			case errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchKey):
				return false, nil
			default:
				return false, errors.New(s("aws error: %q", err))
			}
		}
	}
	return true, nil
}

func (op *AwsOperator) Upload(po provider.PutObject) (err error) {
	obj, ok := po.Object().(*s3.PutObjectInput)
	if !ok {
		return errors.New(ErrorCouldNotAssertObject)
	}
	_, err = op.svc.Upload(context.Background(), obj)
	if err != nil {
		return err
	}
	return nil
}

func (op *AwsOperator) BucketExists() (exists bool, errs provider.Errs) {
	_, err := op.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(op.ac.Bucket.Name),
	})
	if err != nil {
		if errors.As(err, &s3Error) {
			switch {
			case errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchBucket):
				exists = false
				errs.Add(errors.New(s("aws says bucket %q does not exist", op.ac.Bucket.Name)))
			}
		} else {
			exists = false
			errs.Add(errors.New(s("aws error when checking if %q exists: %q", op.ac.Bucket.Name, err)))
		}
		return exists, errs
	} else {
		return true, errs
	}
}
