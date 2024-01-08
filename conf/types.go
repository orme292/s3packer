package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"github.com/rs/zerolog"
)

/*
readConfig is only used to unmarshal a YAML profile, it is not used in the application.
*/
type readConfig struct {
	// Version will be used for feature support
	Version int `yaml:"Version"`

	// AWS will contain only AWS specific configuration details. Other providers will have their own
	// struct and fields.
	AWS struct {
		Profile string `yaml:"Profile"`
		Key     string `yaml:"Key"`
		Secret  string `yaml:"Secret"`
		ACL     string `yaml:"ACL"`
		Storage string `yaml:"Storage"`
	} `yaml:"AWS"`

	// Bucket should be universal across all providers, though there may be different fields depending on the
	// provider.
	Bucket struct {
		Create bool   `yaml:"Create"`
		Name   string `yaml:"Name"`
		Region string `yaml:"Region"`
	} `yaml:"Bucket"`

	// Objects contains object level configuration details, mostly related to object naming. Tagging is separate.
	Objects struct {
		NamePrefix  string `yaml:"NamePrefix"`
		RootPrefix  string `yaml:"RootPrefix"`
		Naming      string `yaml:"Naming"`
		OmitRootDir bool   `yaml:"OmitRootDir"`
	} `yaml:"Objects"`

	// Options is used to configure the application and how it operates.
	Options struct {
		MaxParts   int    `yaml:"MaxParts"`
		MaxUploads int    `yaml:"MaxUploads"`
		Overwrite  string `yaml:"Overwrite"`
	} `yaml:"Options"`

	// Tagging is used only for object tagging.
	Tagging struct {
		ChecksumSHA256 bool              `yaml:"Checksum"`
		Origins        bool              `yaml:"Origins"`
		Tags           map[string]string `yaml:"Tags"`
	} `yaml:"Tagging"`

	// Uploads contains the list of files, folders, and directories to upload.
	Uploads struct {
		Files       []string `yaml:"Files"`
		Folders     []string `yaml:"Folders"`
		Directories []string `yaml:"Directories"`
	} `yaml:"Uploads"`

	// Logging is used to configure the logging output, which is handled by logbot/zerolog.
	Logging struct {
		Level    int    `yaml:"Level"`
		Console  bool   `yaml:"Console"`
		File     bool   `yaml:"File"`
		Filepath string `yaml:"Filepath"`
	} `yaml:"Logging"`

	// Log is an instance of the logger.
	Log *logbot.LogBot
}

// Provider contains all details related to the provider, for any provider.
// Provider.Is can be used to determine which provider is specified in the loaded profile. The ProviderName type
// is a string enum of the supported providers.
// Provider level configuration fields are prefaced with the provider name, e.g. AwsProfile, AwsACL, etc, whereas
// Key and Secret may be used by multiple providers.
type Provider struct {
	Is         ProviderName
	AwsProfile string
	AwsACL     types.ObjectCannedACL
	AwsStorage types.StorageClass
	Key        string
	Secret     string
}

// Bucket contains all details related to the bucket, for any provider. Create is not implemented.
type Bucket struct {
	Create bool
	Name   string
	Region string
}

// Objects contain the object naming configuration.
type Objects struct {
	NamePrefix string
	RootPrefix string
	Naming     Naming

	// OmitRootDir is used to remove the root directory name from the object's final FormattedKey.
	OmitRootDir bool
}

// Opts contains application level configuration options.
type Opts struct {
	MaxParts   int
	MaxUploads int
	Overwrite  Overwrite
}

// TagOpts contain the object tagging configuration, but only the ones handled internally by the application.
// Custom tags are put in a separate map named "Tags" inside the AppConfig struct.
type TagOpts struct {
	ChecksumSHA256       bool
	AwsChecksumAlgorithm types.ChecksumAlgorithm
	AwsChecksumMode      types.ChecksumMode
	Origins              bool
}

// LogOpts contains the logging configuration, but not an instance of logbot/zerolog.
type LogOpts struct {
	Level    zerolog.Level
	Console  bool
	File     bool
	Filepath string
}

// AppConfig is an application level struct that all profile configuration details are loaded into.
// Log is an instance zerolog.Logger, Files and Directories are the list of files and directories to upload.
type AppConfig struct {
	Provider    *Provider
	Bucket      *Bucket
	Objects     *Objects
	Opts        *Opts
	Tags        map[string]string
	Tag         *TagOpts
	LogOpts     *LogOpts
	Log         *logbot.LogBot
	Files       []string
	Directories []string
}

// ProviderName type is a string enum of the supported providers, meant to make it easier to check
// AppConfig.Provider.Is when figuring out which provider config fields should be used.
// ProviderName.String() will return the string representation of the enum for convenience, either in output or logging.
type ProviderName string

const (
	ProviderNameNone ProviderName = "none"
	ProviderNameAWS  ProviderName = "aws"
	ProviderNameOCI  ProviderName = "oci"
	ProviderNameGCP  ProviderName = "gcp"
)

func (p ProviderName) String() string {
	return string(p)
}

// Overwrite type is a string enum of the supported overwrite methods. OverwriteChecksum is not implemented.
// Overwrite.String() will return the string representation of the enum for convenience, either in output or logging.
type Overwrite string

const (
	OverwriteChecksum Overwrite = "checksum"
	OverwriteNever    Overwrite = "never"
	OverwriteAlways   Overwrite = "always"
)

func (o Overwrite) String() string {
	return string(o)
}

// Naming type is a string enum of the supported object naming methods.
// Naming.String() will return the string representation of the enum of convenience, either in output or logging.
type Naming string

const (
	NamingRelative Naming = "relative"
	NamingAbsolute Naming = "absolute"
)

func (n Naming) String() string {
	return string(n)
}

// S is a shortcut for fmt.Sprintf. Really the only purpose is to reduce the number of times that is called.
func S(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

const (
	Empty = ""
)

// Errors

const (
	InvalidACL             = "invalid acl"
	InvalidNamingMethod    = "invalid object naming method"
	InvalidStorageClass    = "invalid storage class"
	InvalidOverwriteMethod = "invalid overwrite method"
	InvalidTagChars        = "invalid characters removed from tag"

	ErrorProfilePath                 = "error determining profile path"
	ErrorOpeningProfile              = "error opening profile"
	ErrorReadingYaml                 = "error reading yaml"
	ErrorAWSProfileAndKeys           = "both aws profile and keys are specified, use profile or keys"
	ErrorAWSKeyOrSecretNotSpecified  = "profile should specified both key and secret"
	ErrorLoggingFilepathNotSpecified = "path to log file not specified"
	ErrorLoggingFilepath             = "error determining log file path"
	ErrorGettingFileInfo             = "error getting file info"
	ErrorFileIsDirectory             = "listed file is directory"
	ErrorNoFilesSpecified            = "no files, folders, directories specified"
	ErrorNoReadableFiles             = "no readable files or directories specified"
	ErrorUnsupportedProfileVersion   = "profile version 2 required"
	ErrorProviderNotSpecified        = "provider not specified"
	ErrorBucketNotSpecified          = "bucket or region not specified"
)
