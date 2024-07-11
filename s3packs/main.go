package s3packs

import (
	"fmt"
	"os"

	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/handler_aws"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

func Do(app *conf.AppConfig) (*provider_v2.Stats, error) {

	h, err := provider_v2.NewHandler(app, handler_aws.NewAwsOperator, handler_aws.NewAwsObject)
	if err != nil {
		app.Tui.ScreenQuit()
		fmt.Printf("Error creating AWS handler: %s\n", err)
		os.Exit(1)
	}

	err = h.Run()
	if err != nil {
		app.Tui.ScreenQuit()
		os.Exit(1)
	}

	app.Tui.ToScreen("Finished.", true)
	app.Tui.ToScreenHeader("s3packer is finished!")

	app.Tui.ScreenQuit()

	return h.Stats, nil

}

// func Do(ac *conf.AppConfig) (stats *objectify.Stats, errs provider.Errs) {
// 	ops, fn, err := getProvider(ac)
// 	if err != nil {
// 		errs.Add(err)
// 		return
// 	}
// 	p, err := provider.NewProcessor(ac, ops, fn)
// 	if err != nil {
// 		errs.Add(err)
// 		return
// 	}
//
// 	if p != nil {
// 		errs = p.Run()
// 	} else {
// 		ac.Log.Fatal("Processor is empty.")
// 	}
// 	return p.Stats, errs
// }
//
// func getProvider(ac *conf.AppConfig) (ops provider.Operator, fn provider.IteratorFunc, err error) {
// 	switch ac.Provider.Is {
// 	case conf.ProviderNameAWS:
// 		ops, err = pack_aws.NewAwsOperator(ac)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		return ops, pack_aws.AwsIteratorFunc, nil
// 	case conf.ProviderNameOCI:
// 		ops, err = pack_oci.NewOracleOperator(ac)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		return ops, pack_oci.OracleIteratorFunc, nil
// 	case conf.ProviderNameLinode:
// 		ops, err = pack_akamai.NewAkamaiOperator(ac)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		return ops, pack_akamai.AkamaiIteratorFunc, nil
// 	default:
// 		return nil, nil, errors.New("unknown provider")
// 	}
// }
