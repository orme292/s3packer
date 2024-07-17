package oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/oracle/oci-go-sdk/v65/objectstorage/transfer"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type OracleOperator struct {
	App    *conf.AppConfig
	Oracle *OracleClient
}

func NewOracleOperator(app *conf.AppConfig) (oper provider_v2.Operator, err error) {

	client, err := NewOracleClient(app.Provider.OCI.Profile, app.Provider.OCI.Compartment)
	if err != nil {
		return nil, err
	}

	oper = &OracleOperator{
		App:    app,
		Oracle: client,
	}

	return oper, nil

}

func (oper *OracleOperator) BucketCreate() error {

	request := objectstorage.CreateBucketRequest{
		NamespaceName: common.String(oper.Oracle.details.namespace),
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			CompartmentId: common.String(oper.Oracle.details.compartment),
			Name:          &oper.App.Bucket.Name,
			FreeformTags: map[string]string{
				"createdBy": "s3packer",
			},
		},
	}

	_, err := oper.Oracle.obj.CreateBucket(context.Background(), request)
	if err != nil {
		return err
	}

	return nil

}

func (oper *OracleOperator) BucketExists() (bool, error) {

	request := objectstorage.HeadBucketRequest{
		NamespaceName: common.String(oper.Oracle.details.namespace),
		BucketName:    common.String(oper.App.Bucket.Name),
	}

	response, err := oper.Oracle.obj.HeadBucket(context.Background(), request)
	if err != nil {
		return false, err
	}

	if response.RawResponse.StatusCode != 200 {
		return false, fmt.Errorf("returned non 200 status code")
	}

	return true, nil

}

func (oper *OracleOperator) BucketDelete() error {
	return nil
}

func (oper *OracleOperator) ObjectExists(obj provider_v2.Object) (bool, error) {

	oobj, ok := obj.(*OracleObject)
	if !ok {
		return true, fmt.Errorf("trouble building object to check")
	}

	request := objectstorage.HeadObjectRequest{
		NamespaceName: common.String(oper.Oracle.details.namespace),
		BucketName:    common.String(oper.App.Bucket.Name),
		ObjectName:    common.String(oobj.key),
	}

	response, err := oper.Oracle.obj.HeadObject(context.Background(), request)
	if err != nil {
		return true, err
	}

	if response.ETag != common.String(EmptyString) {
		return true, nil
	}

	return false, fmt.Errorf("object etag not found: %s", oobj.key)

}

func (oper *OracleOperator) ObjectDelete(key string) error {
	return nil
}

func (oper *OracleOperator) ObjectUpload(obj provider_v2.Object) error {

	oobj, ok := obj.(*OracleObject)
	if !ok {
		return fmt.Errorf("trouble building object to upload")
	}

	request := transfer.UploadRequest{
		NamespaceName:       common.String(oper.Oracle.details.namespace),
		BucketName:          common.String(oper.App.Bucket.Name),
		ObjectName:          common.String(oobj.key),
		ObjectStorageClient: oper.Oracle.obj,
		Metadata:            oobj.tags,
		StorageTier:         oper.App.Provider.OCI.PutStorage,
	}

	oobj.setTagsWithWorkaround(oobj.job.Metadata.SizeBytes)
	response, err := oper.Oracle.manager.UploadFile(context.Background(), transfer.UploadFileRequest{
		UploadRequest: request,
		FilePath:      oobj.job.Metadata.FullPath(),
	})
	if err != nil {
		return err
	}

	if response.SinglepartUploadResponse != nil {
		if response.SinglepartUploadResponse.RawResponse.StatusCode != 200 {
			return fmt.Errorf("upload returned non 200 status code [%d]", response.SinglepartUploadResponse.RawResponse.StatusCode)
		}
	}

	if response.MultipartUploadResponse != nil {
		if response.MultipartUploadResponse.RawResponse.StatusCode != 200 {
			return fmt.Errorf("upload returned non 200 status code [%d]", response.SinglepartUploadResponse.RawResponse.StatusCode)
		}
	}

	return nil

}

func (oper *OracleOperator) GetObjectTags(key string) (map[string]string, error) {
	return make(map[string]string), nil
}

func (oper *OracleOperator) Support() *provider_v2.Supports {

	return provider_v2.NewSupports(true, false, false, false)

}
