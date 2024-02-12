package s3packs

import (
	"errors"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
	"github.com/orme292/s3packer/s3packs/pack_aws"
	"github.com/orme292/s3packer/s3packs/pack_oci"
	"github.com/orme292/s3packer/s3packs/provider"
)

func Do(ac *conf.AppConfig) (stats *objectify.Stats, errs provider.Errs) {
	ops, fn, err := getProvider(ac)
	if err != nil {
		errs.Add(err)
		return
	}
	p, err := provider.NewProcessor(ac, ops, fn)
	if err != nil {
		errs.Add(err)
	}
	errs = p.Run()
	return p.Stats, errs
}

func getProvider(ac *conf.AppConfig) (ops provider.Operator, fn provider.IteratorFunc, err error) {
	switch ac.Provider.Is {
	case conf.ProviderNameAWS:
		ops, err = pack_aws.NewAwsOperator(ac)
		if err != nil {
			return nil, nil, err
		}
		return ops, pack_aws.AwsIteratorFunc, nil
	case conf.ProviderNameOCI:
		ops, err = pack_oci.NewOracleOperator(ac)
		if err != nil {
			return nil, nil, err
		}
		return ops, pack_oci.OracleIteratorFunc, nil
	default:
		return nil, nil, errors.New("unknown provider")
	}
}
