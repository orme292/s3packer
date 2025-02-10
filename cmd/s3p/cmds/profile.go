package cmds

import (
	"github.com/spf13/cobra"
	"s3p/cmd/s3p/cmds/profile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "create or view an upload profile",
	Long:  "create or view an upload profile to configure s3p to upload files to a specific object storage service",
}

func addProfileCmd() {
	profileCmd.AddCommand(profile.CreateProfileCmd)
	profileCmd.AddCommand(profile.CreateSampleProfile)

	rootCmd.AddCommand(profileCmd)
}
