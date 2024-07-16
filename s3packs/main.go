package s3packs

import (
	"errors"
	"fmt"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/handler_aws"
	"github.com/orme292/s3packer/s3packs/handler_linode"
	"github.com/orme292/s3packer/s3packs/provider_v2"
	"github.com/orme292/s3packer/tuipack"
)

func Init(app *conf.AppConfig) (*provider_v2.Stats, error) {

	operFn, objFn, err := getProviderFunctions(app.Provider.Is)
	if err != nil {
		return nil, errors.New("unable to find the correct provider")
	}

	h, err := provider_v2.NewHandler(app, operFn, objFn)
	if err != nil {
		return &provider_v2.Stats{}, err
	}

	err = h.Run()
	if err != nil {
		return &provider_v2.Stats{}, err
	}

	app.Tui.SendOutput(tuipack.NewLogMsg("Finished.", tuipack.ScrnLfCheck,
		tuipack.INFO, "Finished"))
	app.Tui.ToScreenHeader("s3packer is finished!")

	return h.Stats, nil

}

func getProviderFunctions(name conf.ProviderName) (provider_v2.OperGenFunc, provider_v2.ObjectGenFunc, error) {

	switch name {
	case conf.ProviderNameAWS:
		return handler_aws.NewAwsOperator, handler_aws.NewAwsObject, nil

	case conf.ProviderNameLinode:
		return handler_linode.NewLinodeOperator, handler_linode.NewLinodeObject, nil

	default:
		return nil, nil, fmt.Errorf("Unable to determine the provider")

	}

}
