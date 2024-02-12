package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"gopkg.in/yaml.v3"
)

// loadProfile() reads the profile file and returns a readConfig struct. The only validated fields are the logging fields,
// version, and the files and directories. The rest of the fields are left as they are until each individual
// method is called.
func (rc *readConfig) loadProfile(file string) (err error) {
	rc.Logging.Console = true
	rc.Logging.File = false
	rc.Logging.Level = int(logbot.WARN)

	file, err = filepath.Abs(file)
	if err != nil {
		return errors.New(S(ErrorProfilePath, err.Error()))
	}

	f, err := os.ReadFile(filepath.Clean(file))
	if err != nil {
		return errors.New(S(ErrorOpeningProfile, err.Error()))
	}

	err = yaml.Unmarshal(f, &rc)
	if err != nil {
		return errors.New(S(ErrorReadingYaml, err.Error()))
	}

	err = rc.validateLogging()
	rc.Log = &logbot.LogBot{
		Level:       logbot.ParseIntLevel(rc.Logging.Level),
		FlagConsole: rc.Logging.Console,
		FlagFile:    rc.Logging.File,
		Path:        rc.Logging.Filepath,
	}
	if err != nil {
		rc.Log.Warn(err.Error())
	}

	err = rc.validateFiles()
	if err != nil {
		return err
	}
	return nil
}

// transposeStructBucket() returns a Bucket struct. If the bucket name or region is not specified, an error is returned.
// Create is not implemented, so its value doesn't matter.
func (rc *readConfig) transposeStructBucket() (b *Bucket, err error) {
	if rc.Bucket.Name == Empty || rc.Bucket.Region == Empty {
		return nil, errors.New(ErrorBucketNotSpecified)
	}
	return &Bucket{
		Create: rc.Bucket.Create,
		Name:   rc.Bucket.Name,
		Region: rc.Bucket.Region,
	}, nil
}

// transposeStructFileTargets() returns a slice of files and directories to be uploaded. If no files or directories are specified,
// an error is returned.
// Directories and Folders slices are consolidated here, since they are just two different names for the same thing.
// TODO: Check files and dirs for duplicate entries.
// TODO: Add support for globs.
func (rc *readConfig) transposeStructFileTargets() (files, dirs []string, err error) {
	for _, file := range rc.Uploads.Files {
		s, err := os.Stat(file)
		if err != nil {
			rc.Log.Warn("%s: %q", ErrorGettingFileInfo, file)
		} else {
			if s.IsDir() == true {
				dirs = append(dirs, strings.TrimRight(file, "/"))
				rc.Log.Warn("%s: %q", ErrorFileIsDirectory, file)
			} else {
				files = append(files, file)
			}
		}
	}
	for _, folder := range rc.Uploads.Folders {
		dirs = append(dirs, strings.TrimRight(folder, "/"))
	}
	for _, dir := range rc.Uploads.Directories {
		dirs = append(dirs, strings.TrimRight(dir, "/"))
	}
	if len(files) == 0 && len(dirs) == 0 {
		return nil, nil, errors.New(ErrorNoReadableFiles)
	}
	return
}

// transposeStructLogging() returns a LogOpts struct. If the logging file is enabled and the path is specified, then the
// path is converted to an absolute path. Any actual validation is handled elsewhere.
func (rc *readConfig) transposeStructLogging() (lo *LogOpts, err error) {
	var abs string
	if rc.Logging.File == true && rc.Logging.Filepath != Empty {
		abs, err = filepath.Abs(filepath.Clean(rc.Logging.Filepath))
		if err != nil {
			return nil, errors.New(S("%s: %s", ErrorLoggingFilepath, err.Error()))
		}
	}

	return &LogOpts{
		Level:    logbot.ParseIntLevel(rc.Logging.Level),
		Console:  rc.Logging.Console,
		File:     rc.Logging.File,
		Filepath: abs,
	}, nil
}

// transposeStructObjects() returns an Objects struct. If the naming method is not specified, then the default is used, but
// an error is returned.
func (rc *readConfig) transposeStructObjects() (o *Objects, err error) {
	var method Naming
	switch strings.ToLower(strings.Trim(rc.Objects.Naming, " ")) {
	case NamingAbsolute.String():
		method = NamingAbsolute
	case NamingRelative.String():
		method = NamingRelative
	default:
		method = NamingAbsolute
		err = errors.New(InvalidNamingMethod)
	}
	return &Objects{
		NamePrefix:  strings.TrimPrefix(rc.Objects.NamePrefix, "/"),
		RootPrefix:  formatPath(rc.Objects.RootPrefix),
		Naming:      method,
		OmitRootDir: rc.Objects.OmitRootDir,
	}, err
}

// transposeStructOpts() returns an Opts struct. If the overwrite method is not specified, then the default is used, and the
// default should always be to never overwrite an object. OverwriteChecksum support is not implemented. A MaxParts
// value greater or less than 1 is not supported.
func (rc *readConfig) transposeStructOpts() (opts *Opts, err error) {
	overwrite := OverwriteNever

	switch tidyString(rc.Options.Overwrite) {
	case OverwriteAlways.String():
		overwrite = OverwriteAlways
	case OverwriteNever.String():
		overwrite = OverwriteNever
	default:
		err = errors.New(InvalidOverwriteMethod)
	}

	return &Opts{
		MaxParts:   rc.Options.MaxParts,
		MaxUploads: rc.Options.MaxUploads,
		Overwrite:  overwrite,
	}, err
}

// transposeStructProvider() returns a Provider struct. If the provider is not specified, then an error is returned.
// There is only support for a single provider right now, so there's no real complexity here.
func (rc *readConfig) transposeStructProvider() (p *Provider, err error) {
	provider := whichProvider(rc.Provider)
	switch provider {
	case ProviderNameAWS:
		err = rc.validateProviderAWS()
		if err != nil {
			return nil, err
		}
		return rc.buildProviderAWS(), err
	case ProviderNameOCI:
		err = rc.validateProviderOCI()
		if err != nil {
			return nil, err
		}
		return rc.buildProviderOCI(), err
	default:
		return &Provider{Is: ProviderNameNone}, errors.New(ErrorProviderNotSpecified)
	}
}

func (rc *readConfig) buildProviderAWS() (p *Provider) {
	acl, err := awsMatchACL(rc.AWS.ACL)
	if err != nil {
		rc.Log.Warn(err.Error())
	}

	class, err := awsMatchStorage(rc.AWS.Storage)
	if err != nil {
		rc.Log.Warn(err.Error())
	}

	return &Provider{
		Is: ProviderNameAWS,
		AWS: &ProviderAWS{
			Profile: rc.AWS.Profile,
			Key:     rc.AWS.Key,
			Secret:  rc.AWS.Secret,
			ACL:     acl,
			Storage: class,
		},
		Key:    rc.AWS.Key,
		Secret: rc.AWS.Secret,
	}
}

func (rc *readConfig) buildProviderOCI() (p *Provider) {
	if strings.TrimSpace(strings.ToUpper(rc.OCI.Profile)) == OciDefaultProfile {
		rc.OCI.Profile = OciDefaultProfile
	}
	return &Provider{
		Is: ProviderNameOCI,
		OCI: &ProviderOCI{
			Profile:     strings.TrimSpace(rc.OCI.Profile),
			Compartment: rc.OCI.Compartment,
		},
	}
}

// transposeStructTagOpts() returns a TagOpts struct.
func (rc *readConfig) transposeStructTagOpts() (t *TagOpts, err error) {
	return &TagOpts{
		ChecksumSHA256:       rc.Tagging.ChecksumSHA256,
		AwsChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		AwsChecksumMode:      types.ChecksumModeEnabled,
		Origins:              rc.Tagging.Origins,
	}, nil
}

// transposeStructTags() returns a map of valid string tags. This is both a get and validate method. It removes
// unsupported symbols from the tag key/values specified in the profile. Changes are logged as Warnings, but they
// don't halt execution.
func (rc *readConfig) transposeStructTags() (tags map[string]string, err error) {
	tags = make(map[string]string)
	for k, v := range rc.Tagging.Tags {
		newKey := alphaNumericString(k)
		newValue := alphaNumericString(v)
		if newKey != k {
			rc.Log.Warn(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, k, newKey))
		}
		if newValue != v {
			rc.Log.Warn(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, v, newValue))
		}
		tags[newKey] = newValue
	}
	return tags, nil
}
