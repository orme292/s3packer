package pack_oci

import (
	"context"
	"errors"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
	"github.com/oracle/oci-go-sdk/v49/objectstorage/transfer"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider"
)

func NewOracleOperator(ac *conf.AppConfig) (*OracleOperator, error) {
	ociIdentity, ociClient, err := buildClients(ac)
	if err != nil {
		return nil, err
	}
	ociNamespace, err := getNamespace(ac, ociClient)
	if err != nil {
		return nil, err
	}
	return &OracleOperator{
		ac:        ac,
		identity:  ociIdentity,
		client:    ociClient,
		namespace: ociNamespace,
		um:        transfer.NewUploadManager(),
	}, err
}

func (op *OracleOperator) SupportsMultipartUploads() bool { return false }

func (op *OracleOperator) CreateBucket() (err error) {
	request := objectstorage.CreateBucketRequest{
		NamespaceName: common.String(op.namespace),
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			CompartmentId: common.String(op.ac.Provider.OCI.Compartment),
			Name:          &op.ac.Bucket.Name,
			Metadata: map[string]string{
				"createdBy": "s3packer",
			},
		},
	}

	response, err := op.client.CreateBucket(context.Background(), request)
	if err != nil {
		return
	}

	op.ac.Log.Info("Created bucket %q (OCID:%s)", response.Bucket.Name, response.Bucket.Id)
	return
}

func (op *OracleOperator) ObjectExists(key string) (exists bool, err error) {
	request := objectstorage.HeadObjectRequest{
		NamespaceName: common.String(op.namespace),
		BucketName:    common.String(op.ac.Bucket.Name),
		ObjectName:    common.String(key),
	}

	response, err := op.client.HeadObject(context.Background(), request)
	if err != nil {
		return
	}

	op.ac.Log.Debug("Object %q with status code %d and ETag %q", key, response.RawResponse.StatusCode, response.ETag)
	if response.ETag != common.String(EmptyString) {
		return true, nil
	}
	return false, errors.New(s("object %q not found", key))
}

func (op *OracleOperator) Upload(po provider.PutObject) (err error) {
	ur := po.Object().(transfer.UploadRequest)
	response, err := op.um.UploadFile(context.Background(), transfer.UploadFileRequest{
		UploadRequest: ur,
		FilePath:      po.Fo().AbsPath,
	})
	if err != nil {
		return err
	}
	if response.SinglepartUploadResponse != nil {
		if response.SinglepartUploadResponse.RawResponse.StatusCode != 200 {
			return errors.New(s("(single) upload failed with status code %d", response.SinglepartUploadResponse.RawResponse.StatusCode))
		}
	}
	if response.MultipartUploadResponse != nil {
		if response.MultipartUploadResponse.RawResponse.StatusCode != 200 {
			return errors.New(s("(multi) upload failed with status code %d", response.MultipartUploadResponse.RawResponse.StatusCode))
		}
	}
	op.ac.Log.Debug("Object %q upload with UploadID %v", *ur.ObjectName, response.Type)
	return nil
}

func (op *OracleOperator) UploadMultipart(po provider.PutObject) (err error) {
	return errors.New("explicit multipart uploads are not used")
}

func (op *OracleOperator) BucketExists() (exists bool, errs provider.Errs) {
	exists = false
	request := objectstorage.HeadBucketRequest{
		NamespaceName: common.String(op.namespace),
		BucketName:    common.String(op.ac.Bucket.Name),
	}

	response, err := op.client.HeadBucket(context.Background(), request)
	if err != nil {
		errs.Add(err)
		return exists, errs
	}
	op.ac.Log.Debug("HeadBucket %q returned status code %d", op.ac.Bucket.Name, response.RawResponse.StatusCode)

	if response.RawResponse.StatusCode != 200 {
		op.ac.Log.Info("Bucket %q not found", op.ac.Bucket.Name)
		return exists, errs
	}
	op.ac.Log.Info("Found bucket named %q", op.ac.Bucket.Name)
	exists = true
	return exists, errs
}
