package profile

const (
	CreateProfileFlagFilename  = "filename"
	CreateProfileFlagFilenameS = "f"
)

func init() {
	CreateProfileCmd.Flags().StringP(CreateProfileFlagFilename, CreateProfileFlagFilenameS, "./sample-profile.yml", "filename for the new profile")
	_ = CreateProfileCmd.MarkFlagRequired(CreateProfileFlagFilename)

	CreateSampleProfile.Flags().StringP(CreateProfileFlagFilename, CreateProfileFlagFilenameS, "./sample-profile.yml", "filename for the new profile")
	_ = CreateSampleProfile.MarkFlagRequired(CreateProfileFlagFilename)
}
