package s3packs

import (
	"errors"
	"fmt"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider_aws"
	"github.com/orme292/s3packer/s3packs/provider_linode"
	"github.com/orme292/s3packer/s3packs/provider_v2"
	"github.com/orme292/s3packer/tuipack"
)

func Init(app *conf.AppConfig) (*provider_v2.Stats, error) {

	operFn, objFn, err := getProviderFunctions(app.Provider.Is)
	if err != nil {
		return nil, errors.New("unable to find the correct provider")
	}

	handler, err := provider_v2.NewHandler(app, operFn, objFn)
	if err != nil {
		return &provider_v2.Stats{}, err
	}

	err = handler.Init()
	if err != nil {
		return &provider_v2.Stats{}, err
	}

	app.Tui.SendOutput(tuipack.NewLogMsg("Finished.", tuipack.ScrnLfCheck,
		tuipack.INFO, "Finished"))
	app.Tui.ToScreenHeader("s3packer is finished!")

	return handler.Stats, nil

}

func getProviderFunctions(name conf.ProviderName) (provider_v2.OperGenFunc, provider_v2.ObjectGenFunc, error) {

	switch name {
	case conf.ProviderNameAWS:
		return provider_aws.NewAwsOperator, provider_aws.NewAwsObject, nil

	case conf.ProviderNameLinode:
		return provider_linode.NewLinodeOperator, provider_linode.NewLinodeObject, nil

	default:
		return nil, nil, fmt.Errorf("unable to determine the provider")

	}

}
