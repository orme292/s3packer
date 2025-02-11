package conf

// TagOpts contain the object tagging configuration, but only the ones handled internally by the application.
// Custom tags are put in a separate map named "Tags" inside the AppConfig struct.
type TagOpts struct {
	ChecksumSHA256 bool
	OriginPath     bool
}

func (to *TagOpts) build(inc *ProfileIncoming) error {

	to.OriginPath = inc.TagOptions.OriginPath
	to.ChecksumSHA256 = inc.TagOptions.ChecksumSHA256

	return to.validate()

}

func (to *TagOpts) validate() error {

	// nothing to validate yet
	return nil

}
