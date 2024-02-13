package pack_akamai

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	op.ac.Log.Fatal("CreateBucket() not implemented")
	return
}

func (op *AkamaiOperator) ObjectExists(key string) (exists bool, err error) {
	exists = true
	key = strings.ToUpper(strings.TrimSpace(key))
	return false, nil
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
	exists = true
	return exists, errs
}
