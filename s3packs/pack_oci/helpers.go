package pack_oci

import (
	"github.com/oracle/oci-go-sdk/common"
	"github.com/orme292/s3packer/conf"
)

func buildConfigProvider(ac *conf.AppConfig) (*common.ConfigurationProvider, error) {
	p := common.NewRawConfigurationProvider(
		ac.Provider.OCI.Builder.Tenancy,
		ac.Provider.OCI.Builder.User,
		ac.Provider.OCI.Builder.Region,
		ac.Provider.OCI.Builder.Fingerprint,
		ac.Provider.OCI.Builder.PrivateKey,
		&ac.Provider.OCI.Builder.Passphrase,
	)

	return &p, nil
}
