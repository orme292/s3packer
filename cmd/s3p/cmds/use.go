package cmds

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"s3p/internal/conf"
	"s3p/internal/provider"
	s3aws "s3p/internal/providers/aws"
	s3gcloud "s3p/internal/providers/gcloud"
	s3linode "s3p/internal/providers/linode"
	s3oracle "s3p/internal/providers/oracle"
)

var UseCmd = &cobra.Command{
	Use:   "use",
	Short: "use an upload profile",
	Long:  "use an upload profile to configure s3p to upload files to a specific object storage service",
	Run:   useProfile,
}

func addUseCmd() {
	rootCmd.AddCommand(UseCmd)
}

func useProfile(cmd *cobra.Command, args []string) {
	filename, err := cmd.Flags().GetString(UseProfileFilenameFlag)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to retrieve '%s' flag: %v", UseProfileFilenameFlag, err))
	}

	builder := conf.NewBuilder(filename)
	app, err := builder.FromYaml()
	if err != nil {
		log.Fatalf("Failed to load profile: %v", err)
	}

	fmt.Printf("\ns3p\n")
	fmt.Printf("Using %s and bucket %q\n\n", app.Provider.Is.Title(), app.Bucket)
	time.Sleep(1 * time.Second)

	stats, err := appInit(app)
	if err != nil {
		log.Printf("s3p exited with error: %q\n\n", err.Error())
		os.Exit(1)
	}

	hrb := stats.ReadableString()
	msg := fmt.Sprintf("%s uploaded, %s skipped", hrb[stats.ObjectsBytes],
		hrb[stats.SkippedBytes])
	app.Log.Info(msg)
	app.Log.Info(stats.String())

	os.Exit(0)
}

func appInit(app *conf.AppConfig) (*provider.Stats, error) {
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

	app.Log.Info("Finished.")

	return handler.Stats, nil
}

func getProviderFunctions(name conf.ProviderName) (provider.OperGenFunc, provider.ObjectGenFunc, error) {
	switch name {
	case conf.ProviderNameAWS:
		return s3aws.NewAwsOperator, s3aws.NewAwsObject, nil
	case conf.ProviderNameGoogle:
		return s3gcloud.NewGCloudOperator, s3gcloud.NewCloudObject, nil
	case conf.ProviderNameLinode:
		return s3linode.NewLinodeOperator, s3linode.NewLinodeObject, nil
	case conf.ProviderNameOCI:
		return s3oracle.NewOracleOperator, s3oracle.NewOracleObject, nil
	default:
		return nil, nil, fmt.Errorf("unable to determine the provider")
	}
}

func startSigWatcher() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		os.Exit(0)
	}()
}
