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

func (r *readConfig) getTargets() (files []string, dirs []string, err error) {
	for _, file := range r.Uploads.Files {
		s, err := os.Stat(file)
		if err != nil {
			r.Log.Warn("%s: %q", ErrorGettingFileInfo, file)
		} else if s.IsDir() == true {
			dirs = append(dirs, strings.TrimRight(file, "/"))
			r.Log.Warn("%s: %q", ErrorFileIsDirectory, file)
		} else {
			files = append(files, file)
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

func (r *readConfig) getLogging() (lo *LogOpts, err error) {
	abs, err := filepath.Abs(filepath.Clean(r.Logging.Filepath))
	if err != nil {
		return nil, errors.New(S("%s: %s", ErrorLoggingFilepath, err.Error()))
	}
	return &LogOpts{
		Level:    logbot.ParseIntLevel(r.Logging.Level),
		Console:  r.Logging.Console,
		File:     r.Logging.File,
		Filepath: abs,
	}, nil
}

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

func (r *readConfig) getOpts() (opts *Opts, err error) {
	var overwrite Overwrite
	switch strings.ToLower(strings.Trim(r.Options.Overwrite, " ")) {
	case OverwriteAlways.String():
		overwrite = OverwriteAlways
	//case OverwriteChecksum.String():
	//	overwrite = OverwriteChecksum
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

func (r *readConfig) getTag() (t *Tag, err error) {
	return &Tag{
		ChecksumSHA256:       r.Tagging.ChecksumSHA256,
		AwsChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		AwsChecksumMode:      types.ChecksumModeEnabled,
		Origins:              r.Tagging.Origins,
	}, nil
}

func (r *readConfig) getValidTags() (tags map[string]string, err error) {
	tags = make(map[string]string)
	for k, v := range r.Tagging.Tags {
		reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
		nk := reg.ReplaceAllString(k, "")
		nv := reg.ReplaceAllString(v, "")
		if nk != k {
			r.Log.Info(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, k, nk))
		}
		if nv != v {
			r.Log.Info(fmt.Sprintf("%s: %q is now %q", InvalidTagChars, v, nv))
		}
		tags[nk] = nv
	}
	return tags, nil
}

func (r *readConfig) validateFiles() (err error) {
	if len(r.Uploads.Files) == 0 && len(r.Uploads.Folders) == 0 && len(r.Uploads.Directories) == 0 {
		err = errors.New(ErrorNoFilesSpecified)
	}
	return
}

func (r *readConfig) validateLogging() (err error) {
	if r.Logging.File == true && r.Logging.Filepath == Empty {
		err = errors.New(ErrorLoggingFilepathNotSpecified)
		r.Logging.File = false
	}
	return
}

func (r *readConfig) versionOK() (v int, err error) {
	if r.Version != 2 {
		return r.Version, errors.New(ErrorUnsupportedProfileVersion)
	}
	return r.Version, nil
}

func (r *readConfig) validateProviderAWS() (err error) {
	if r.AWS.Profile != Empty && (r.AWS.Key != Empty || r.AWS.Secret != Empty) {
		err = errors.New(ErrorAWSProfileAndKeys)
	}
	return
}
