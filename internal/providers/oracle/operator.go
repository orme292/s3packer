package oci

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/oracle/oci-go-sdk/v65/objectstorage/transfer"

	"s3p/internal/conf"
	"s3p/internal/provider"
)

type OracleOperator struct {
	App    *conf.AppConfig
	Oracle *OracleClient
}

func NewOracleOperator(app *conf.AppConfig) (oper provider.Operator, err error) {

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
	code, msg, t := getResponseCode(response)

	if code >= 200 && code <= 299 {
		return true, nil
	} else if code == 404 {
		return false, nil
	}

	logmsg := fmt.Sprintf("OCI returned [%s] code %d, msg: %s", t.String(), code, msg)
	oper.App.Tui.Info(logmsg)

	return false, err

}

func (oper *OracleOperator) BucketDelete() error {
	return nil
}

func (oper *OracleOperator) ObjectExists(obj provider.Object) (bool, error) {

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
	code, msg, t := getResponseCode(response)

	if code >= 200 && code <= 299 {
		if response.ETag != nil && response.ETag != common.String(EmptyString) {
			return true, nil
		}
		return true, fmt.Errorf("OK, but with unexpected ETag")
	} else if code == 404 {
		return false, fmt.Errorf("object not found")
	}

	logmsg := fmt.Sprintf("OCI returned [%s] code %d, msg: %s", t.String(), code, msg)
	oper.App.Tui.Info(logmsg)

	return true, err
}

func (oper *OracleOperator) ObjectDelete(key string) error {
	return nil
}

func (oper *OracleOperator) ObjectUpload(obj provider.Object) error {

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

func (oper *OracleOperator) Support() *provider.Supports {

	return provider.NewSupports(true, false, false, false)

}

func getResponseCode(response any) (code int, msg string, t reflect.Type) {

	var raw http.Response
	switch v := response.(type) {
	case objectstorage.HeadBucketResponse:
		raw = *v.HTTPResponse()
	case objectstorage.HeadObjectResponse:
		raw = *v.HTTPResponse()
	case transfer.UploadResponse:
		if response.(transfer.UploadResponse).SinglepartUploadResponse != nil {
			raw = *v.SinglepartUploadResponse.HTTPResponse()
		}
		if response.(transfer.UploadResponse).MultipartUploadResponse != nil {
			raw = *v.MultipartUploadResponse.HTTPResponse()
		}
	}

	switch code = raw.StatusCode; {
	case code == 401:
		msg = "not authenticated"
	case code == 403:
		msg = "invalid region"
	case code == 404:
		msg = "not found"
	case code == 409:
		msg = "resource already exists"
	case code == 429:
		msg = "too many requests"
	case code >= 500 && code <= 503:
		msg = "internal server error"
	case code >= 200 && code <= 299:
		msg = "ok"
	default:
		msg = "unknown response"
	}

	return raw.StatusCode, msg, reflect.TypeOf(response)

}
