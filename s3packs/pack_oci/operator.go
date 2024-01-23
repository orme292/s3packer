package pack_oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/orme292/s3packer/conf"
)

func NewOracleOperator(ac *conf.AppConfig) (*OracleOperator, error) {
	return nil, nil
}

func (op *OracleOperator) CreateBucket() (err error) {
	op.p = common.DefaultConfigProvider()
	return nil
}
