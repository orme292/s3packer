package profile

import (
	"github.com/spf13/cobra"
)

var CreateProfileCmd = &cobra.Command{
	Use:   "create",
	Short: "create an upload profile with interactive prompts",
	Long:  "create an upload profile using interactive prompts for the configuration and service authentication values",
	Run:   createProfile,
}

func createProfile(cmd *cobra.Command, args []string) {

}
