package pack_akamai

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider"
)

func NewAkamaiOperator(ac *conf.AppConfig) (*AkamaiOperator, error) {
	svc, client, err := buildUploader(ac)
	if err != nil {
		return nil, err
	}
	return &AkamaiOperator{
		ac:     ac,
		client: client,
		svc:    svc,
	}, nil
}

func (op *AkamaiOperator) SupportsMultipartUploads() bool { return false }

func (op *AkamaiOperator) CreateBucket() (err error) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(op.ac.Bucket.Name),
		ACL:    types.BucketCannedACLPrivate,
	}
	_, err = op.client.CreateBucket(context.Background(), input)
	if err != nil {
		op.ac.Log.Fatal("Unable to create bucket %q in %q: %q", op.ac.Bucket.Name, op.ac.Bucket.Region, err.Error())
		return err
	}
	op.ac.Log.Info("Created bucket %q in %q", op.ac.Bucket.Name, op.ac.Bucket.Region)
	return
}

func (op *AkamaiOperator) ObjectExists(key string) (exists bool, err error) {
	exists = false
	input := &s3.HeadObjectInput{
		Bucket: aws.String(op.ac.Bucket.Name),
		Key:    aws.String(key),
	}
	_, err = op.client.HeadObject(context.Background(), input)
	if err != nil {
		return exists, errors.New(s("object %q not found", key))
	}
	exists = true
	return exists, err
}

func (op *AkamaiOperator) Upload(po provider.PutObject) (err error) {
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

func (op *AkamaiOperator) UploadMultipart(po provider.PutObject) (err error) {
	return errors.New("multipart uploads are not supported")
}

func (op *AkamaiOperator) BucketExists() (exists bool, errs provider.Errs) {
	exists = false
	input := s3.HeadBucketInput{
		Bucket: aws.String(op.ac.Bucket.Name),
	}
	_, err := op.client.HeadBucket(context.Background(), &input)
	if err != nil {
		errs.Add(errors.New(s("akamai says bucket %q does not exist", op.ac.Bucket.Name)))
		return exists, errs
	}
	exists = true
	return exists, errs
}
