package cmds

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	UseProfileFilenameFlag  = "filename"
	UseProfileFilenameFlagS = "f"
)

var rootCmd = &cobra.Command{
	Use:   "s3p",
	Short: "s3p is a tool for uploading files to object storage services",
	Long: "s3p is a tool used to upload and backup files to object storage services like Amazon S3, Google Cloud" +
		"Storage,  Linode Object Storage, and Oracle Cloud Object Storage.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	addProfileCmd()
	addUseCmd()

	UseCmd.Flags().StringP(UseProfileFilenameFlag, UseProfileFilenameFlagS, "./sample-profile.yml", "filename for the new profile")
	_ = UseCmd.MarkFlagRequired(UseProfileFilenameFlag)
}
