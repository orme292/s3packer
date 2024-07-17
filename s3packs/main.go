package s3packs

import (
	"errors"
	"fmt"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider"
	s3aws "github.com/orme292/s3packer/s3packs/providers/aws"
	s3linode "github.com/orme292/s3packer/s3packs/providers/linode"
	s3oracle "github.com/orme292/s3packer/s3packs/providers/oracle"
	"github.com/orme292/s3packer/tuipack"
)

func Init(app *conf.AppConfig) (*provider.Stats, error) {

	operFn, objFn, err := getProviderFunctions(app.Provider.Is)
	if err != nil {
		return nil, errors.New("unable to find the correct provider")
	}

	handler, err := provider.NewHandler(app, operFn, objFn)
	if err != nil {
		return &provider.Stats{}, err
	}

	err = handler.Init()
	if err != nil {
		return &provider.Stats{}, err
	}

	app.Tui.SendOutput(tuipack.NewLogMsg("Finished.", tuipack.ScrnLfCheck,
		tuipack.INFO, "Finished"))
	app.Tui.ToScreenHeader("s3packer is finished!")

	return handler.Stats, nil

}

func getProviderFunctions(name conf.ProviderName) (provider.OperGenFunc, provider.ObjectGenFunc, error) {

	switch name {
	case conf.ProviderNameAWS:
		return s3aws.NewAwsOperator, s3aws.NewAwsObject, nil

	case conf.ProviderNameLinode:
		return s3linode.NewLinodeOperator, s3linode.NewLinodeObject, nil

	case conf.ProviderNameOCI:
		return s3oracle.NewOracleOperator, s3oracle.NewOracleObject, nil

	default:
		return nil, nil, fmt.Errorf("unable to determine the provider")

	}

}
