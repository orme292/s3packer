package cmds

import (
	"github.com/spf13/cobra"
	"s3p/cmd/s3p/cmds/profile"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "create an empty upload profile",
	Long:  "create an empty profile that can be used as baseline for multiple configurations",
}

func addProfileCmd() {
	profileCmd.AddCommand(profile.SampleProfile)

	rootCmd.AddCommand(profileCmd)
}
