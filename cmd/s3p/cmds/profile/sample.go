package profile

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"s3p/internal/conf"
)

var SampleProfile = &cobra.Command{
	Use:   "sample",
	Short: "create an upload profile with sample values",
	Long:  "create an upload profile with the sample configuration and service authentication values",
	Run:   createSampleProfile,
}

func createSampleProfile(cmd *cobra.Command, args []string) {

	filename, err := cmd.Flags().GetString(CreateProfileFlagFilename)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to retrieve '%s' flag: %v", CreateProfileFlagFilename, err))
	}

	builder := conf.NewBuilder(filename)

	fmt.Printf("Creating sample profile in '%s'\n", builder.Filename)

	err = builder.YamlOut()
	if err != nil {
		log.Fatalf("Failed to create sample profile: %v", err)
	} else {
		fmt.Println("Done")
	}

}
