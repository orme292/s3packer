package conf

import (
	"errors"
)

// validateFiles() checks to make sure that at least one file or directory is specified. If not, then an error
// is returned.
func (rc *readConfig) validateFiles() (err error) {
	if len(rc.Uploads.Files) == 0 && len(rc.Uploads.Folders) == 0 && len(rc.Uploads.Directories) == 0 {
		err = errors.New(ErrorNoFilesSpecified)
	}
	return
}

// validateLogging() checks to make sure that if logging to a file is enabled, then a path is specified. If not,
// then an error is returned. Whether the actual file is accessible or not is not checked.
func (rc *readConfig) validateLogging() (err error) {
	if rc.Logging.File == true && rc.Logging.Filepath == Empty {
		err = errors.New(ErrorLoggingFilepathNotSpecified)
		rc.Logging.File = false
	}
	if rc.Logging.Level > 5 {
		rc.Logging.Level = 5
		err = errors.New(ErrorLoggingLevelTooHigh)
	}
	if rc.Logging.Level < -1 {
		rc.Logging.Level = -1
		err = errors.New(ErrorLoggingLevelTooLow)
	}
	return
}

// validateProviderAWS() checks that the AWS profile and keys are not both specified. If they are, then an error
// is returned. If A key is provided, but not a secret, or vice versa, then an error is returned, also.
func (rc *readConfig) validateProviderAWS() (err error) {
	if rc.AWS.Profile != Empty && (rc.AWS.Key != Empty || rc.AWS.Secret != Empty) {
		err = errors.New(ErrorAWSProfileAndKeys)
	}
	if (rc.AWS.Key == Empty && rc.AWS.Secret != Empty) || (rc.AWS.Key != Empty && rc.AWS.Secret == Empty) {
		err = errors.New(ErrorAWSKeyOrSecretNotSpecified)
	}
	return
}

func (rc *readConfig) validateProviderOCI() (err error) {
	if rc.OCI.Profile == Empty {
		return errors.New(ErrorOCIAuthNotSpecified)
	}
	// This isn't fatal. The provider will just retrieve the tenancy root and use that.
	if rc.OCI.Compartment == Empty {
		rc.Log.Warn(ErrorOCICompartmentNotSpecified)
	}
	return nil
}

// validateVersion() checks that the profile is at version 4; otherwise an error is returned.
// If there are future versions of the profile, then this method will be fleshed out.
// For now, there's only support for version 4
func (rc *readConfig) validateVersion() (v int, err error) {
	if rc.Version < 4 || rc.Version > 4 {
		return rc.Version, errors.New(ErrorUnsupportedProfileVersion)
	}
	return rc.Version, nil
}
