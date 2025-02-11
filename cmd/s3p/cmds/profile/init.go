package profile

const (
	CreateProfileFlagFilename  = "filename"
	CreateProfileFlagFilenameS = "f"
)

func init() {
	SampleProfile.Flags().StringP(CreateProfileFlagFilename, CreateProfileFlagFilenameS, "./sample-profile.yml", "filename for the new profile")
	_ = SampleProfile.MarkFlagRequired(CreateProfileFlagFilename)
}
