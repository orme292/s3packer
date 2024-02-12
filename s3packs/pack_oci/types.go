package pack_oci

import (
	"os"

	"github.com/oracle/oci-go-sdk/common"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
)

const (
	EmptyString = ""
)

type OracleIterator struct {
	provider *conf.Provider
	fol      objectify.FileObjList
	stage    struct {
		index int
		fo    *objectify.FileObj
		f     *os.File
	}
	group int
	err   error
	ac    *conf.AppConfig
}

type OracleOperator struct {
	ac *conf.AppConfig
	cp *common.ConfigurationProvider
}
