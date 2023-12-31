package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"gopkg.in/yaml.v3"
)

// load() reads the profile file and returns a readConfig struct. The only validated fields are the logging fields,
// version, and the files and directories. The rest of the fields are left as they are until each individual
// method is called.
func (r *readConfig) load(file string) (err error) {
	r.Version = 2
	r.Logging.Console = true
	r.Logging.File = false
	r.Logging.Level = int(logbot.WARN)
	file, err = filepath.Abs(file)
	if err != nil {
		return errors.New(S(ErrorProfilePath, err.Error()))
	}

	f, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return errors.New(S(ErrorOpeningProfile, err.Error()))
	}

	err = yaml.Unmarshal(f, &r)
	if err != nil {
		return errors.New(S(ErrorReadingYaml, err.Error()))
	}

	err = r.validateLogging()
	r.Log = &logbot.LogBot{
		Level:       logbot.ParseIntLevel(r.Logging.Level),
		FlagConsole: r.Logging.Console,
		FlagFile:    r.Logging.File,
		Path:        r.Logging.Filepath,
	}
	if err != nil {
		r.Log.Warn(err.Error())
	}

	err = r.validateFiles()
	if err != nil {
		return err
	}
	return nil
}

// getBucket() returns a Bucket struct. If the bucket name or region is not specified, an error is returned.
// Create is not implemented, so it's value doesn't matter.
func (r *readConfig) getBucket() (b *Bucket, err error) {
	if r.Bucket.Name == Empty || r.Bucket.Region == Empty {
		return nil, errors.New(ErrorBucketNotSpecified)
	}
	return &Bucket{
		Create: r.Bucket.Create,
		Name:   r.Bucket.Name,
		Region: r.Bucket.Region,
	}, nil
}

// getTargets() returns a slice of files and directories to be uploaded. If no files or directories are specified,
// an error is returned.
// Directories and Folders slices are consolidated here, since they are just two different names for the same thing.
// TODO: Check files and dirs for duplicate entries.
// TODO: Add support for globs.
func (r *readConfig) getTargets() (files []string, dirs []string, err error) {
	for _, file := range r.Uploads.Files {
		s, err := os.Stat(file)
		if err != nil {
			r.Log.Warn("%s: %q", ErrorGettingFileInfo, file)
		} else {
			if s.IsDir() == true {
				dirs = append(dirs, strings.TrimRight(file, "/"))
				r.Log.Warn("%s: %q", ErrorFileIsDirectory, file)
			} else {
				files = append(files, file)
			}
		}
	}
	for _, folder := range r.Uploads.Folders {
		dirs = append(dirs, strings.TrimRight(folder, "/"))
	}
	for _, dir := range r.Uploads.Directories {
		dirs = append(dirs, strings.TrimRight(dir, "/"))
	}
	if len(files) == 0 && len(dirs) == 0 {
		return nil, nil, errors.New(ErrorNoReadableFiles)
	}
	return
}

// getLogging() returns a LogOpts struct. If the logging file is enabled and the path is specified, then the
// path is converted to an absolute path. Any actual validation is handled elsewhere.
func (r *readConfig) getLogging() (lo *LogOpts, err error) {
	var abs string
	if r.Logging.File == true && r.Logging.Filepath != Empty {
		abs, err = filepath.Abs(filepath.Clean(r.Logging.Filepath))
		if err != nil {
			return nil, errors.New(S("%s: %s", ErrorLoggingFilepath, err.Error()))
		}
	}

	return &LogOpts{
		Level:    logbot.ParseIntLevel(r.Logging.Level),
		Console:  r.Logging.Console,
		File:     r.Logging.File,
		Filepath: abs,
	}, nil
}

// getObjects() returns an Objects struct. If the naming method is not specified, then the default is used, but
// an error is returned.
func (r *readConfig) getObjects() (o *Objects, err error) {
	var method Naming
	switch strings.ToLower(strings.Trim(r.Objects.Naming, " ")) {
	case NamingAbsolute.String():
		method = NamingAbsolute
	case NamingRelative.String():
		method = NamingRelative
	default:
		method = NamingAbsolute
		err = errors.New(InvalidNamingMethod)
	}
	return &Objects{
		NamePrefix:          r.Objects.NamePrefix,
		RootPrefix:          r.Objects.RootPrefix,
		Naming:              method,
		OmitOriginDirectory: r.Objects.OmitOriginDirectory,
	}, err
}

// getOpts() returns an Opts struct. If the overwrite method is not specified, then the default is used, and the
// default should always be to never overwrite an object. OverwriteChecksum support is not implemented. A MaxParts
// value greater or less than 1 is not supported.
func (r *readConfig) getOpts() (opts *Opts, err error) {
	var overwrite Overwrite
	switch strings.ToLower(strings.Trim(r.Options.Overwrite, " ")) {
	case OverwriteAlways.String():
		overwrite = OverwriteAlways
	// case OverwriteChecksum.String():
	//	 overwrite = OverwriteChecksum
	case OverwriteNever.String():
		overwrite = OverwriteNever
	default:
		overwrite = OverwriteNever
		err = errors.New(InvalidOverwriteMethod)
	}
	return &Opts{
		MaxParts:   r.Options.MaxParts,
		MaxUploads: r.Options.MaxUploads,
		Overwrite:  overwrite,
	}, err
}

// getProvider() returns a Provider struct. If the provider is not specified, then an error is returned.
// There is only support for a single provider right now, so there's no real complexity here.
func (r *readConfig) getProvider() (p *Provider, err error) {
	if r.AWS.Profile != Empty || r.AWS.Key != Empty || r.AWS.Secret != Empty {
		err = r.validateProviderAWS()
		if err != nil {
			return nil, err
		}
		acl, err := awsMatchACL(r.AWS.ACL)
		if err != nil {
			r.Log.Warn(err.Error())
		}
		class, err := awsMatchStorage(r.AWS.Storage)
		if err != nil {
			r.Log.Warn(err.Error())
		}
		return &Provider{
			Is:         ProviderNameAWS,
			AwsProfile: r.AWS.Profile,
			AwsACL:     acl,
			AwsStorage: class,
			Key:        r.AWS.Key,
			Secret:     r.AWS.Secret,
		}, nil
	}
	return &Provider{Is: ProviderNameNone}, errors.New(ErrorProviderNotSpecified)
}

// getTag() returns a Tag struct.
func (r *readConfig) getTag() (t *Tag, err error) {
	return &Tag{
		ChecksumSHA256:       r.Tagging.ChecksumSHA256,
		AwsChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		AwsChecksumMode:      types.ChecksumModeEnabled,
		Origins:              r.Tagging.Origins,
	}, nil
}

// getValidTags() returns a map of valid string tags. This is both a get and validate method. It removes
// unsupported symbols from the tag key/values specified in the profile. Changes are logged as Warnings, but they
// don't halt execution.
func (r *readConfig) getValidTags() (tags map[string]string, err error) {
	tags = make(map[string]string)
	for k, v := range r.Tagging.Tags {
		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
		nk := reg.ReplaceAllString(k, "")
		nv := reg.ReplaceAllString(v, "")
		if nk != k {
			r.Log.Warn(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, k, nk))
		}
		if nv != v {
			r.Log.Warn(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, v, nv))
		}
		tags[nk] = nv
	}
	return tags, nil
}

// validateFiles() checks to make sure that at least one file or directory is specified. If not, then an error
// is returned.
func (r *readConfig) validateFiles() (err error) {
	if len(r.Uploads.Files) == 0 && len(r.Uploads.Folders) == 0 && len(r.Uploads.Directories) == 0 {
		err = errors.New(ErrorNoFilesSpecified)
	}
	return
}

// validateLogging() checks to make sure that if logging to a file is enabled, then a path is specified. If not,
// then an error is returned. Whether the actual file is accessible or not is not checked.
func (r *readConfig) validateLogging() (err error) {
	if r.Logging.File == true && r.Logging.Filepath == Empty {
		err = errors.New(ErrorLoggingFilepathNotSpecified)
		r.Logging.File = false
	}
	return
}

// versionOK() checks that the profile is at version 2, otherwise an error is returned. If there are future versions
// of the profile, then this method will be fleshed out. For now, there's only the un-versioned profile and version 2.
func (r *readConfig) versionOK() (v int, err error) {
	if r.Version != 2 {
		return r.Version, errors.New(ErrorUnsupportedProfileVersion)
	}
	return r.Version, nil
}

// validateProviderAWS() checks that the AWS profile and keys are not both specified. If they are, then an error
// is returned. If A key is provided, but not a secret, or vice versa, then an error is returned, also.
func (r *readConfig) validateProviderAWS() (err error) {
	if r.AWS.Profile != Empty && (r.AWS.Key != Empty || r.AWS.Secret != Empty) {
		err = errors.New(ErrorAWSProfileAndKeys)
	}
	if (r.AWS.Key == Empty && r.AWS.Secret != Empty) || (r.AWS.Key != Empty && r.AWS.Secret == Empty) {
		err = errors.New(ErrorAWSKeyOrSecretNotSpecified)
	}
	return
}
