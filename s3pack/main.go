package s3pack

import (
	"github.com/orme292/s3packer/config"
)

// TODO: Overwrite options -- only-if checksum changes (overwrite: always, on-change, never)
// TODO: LogBot, support sprintf style formatting
// TODO: Config, add naming section for KeyNamingMethod, pathPrefix, etc
// TODO: Config, rename indexes from camel case to dashed "pathPrefix" to "path-prefix"
// TODO: Add some console styling, maybe a progress bar.
// TODO: Add silent option
// TODO: Add option to create sample profile YAML
// TODO: Update all comments for each function/method.
// TODO: Update CHANGELOG.md
// TODO: Update README.md
// TODO: Update VERSION
// TODO: Add more readable log output, check log levels make sense
// TODO: Consider ErrorAs implementation and hard coding error messages in Const
// TODO: Upgrade to AWS SDK v2
// TODO: OCI support

func IndividualFileUploader(c *config.Configuration, files []string) (err error, uploaded, ignored int) {
	objList, err := NewObjectList(c, files)
	if err != nil {
		return
	}

	return objList.Upload(c)
}

func DirectoryUploader(c *config.Configuration, dirs []string) (err error, uploaded, ignored int) {
	dirLists, err := NewRootList(c, dirs)
	if err != nil {
		return
	}

	return dirLists.Upload()
}
