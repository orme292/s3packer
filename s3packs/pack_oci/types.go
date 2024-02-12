package pack_oci

import (
	"os"

	"github.com/oracle/oci-go-sdk/v49/identity"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
	"github.com/oracle/oci-go-sdk/v49/objectstorage/transfer"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
)

const (
	EmptyString = ""

	OracleDefaultProfile = "DEFAULT"
)

type OracleIterator struct {
	provider *conf.Provider
	fol      objectify.FileObjList
	stage    struct {
		index int
		fo    *objectify.FileObj
		f     *os.File
	}
	group     int
	err       error
	ac        *conf.AppConfig
	client    objectstorage.ObjectStorageClient
	um        *transfer.UploadManager
	namespace string
}

type OracleOperator struct {
	ac        *conf.AppConfig
	identity  identity.IdentityClient
	client    objectstorage.ObjectStorageClient
	um        *transfer.UploadManager
	namespace string
}
