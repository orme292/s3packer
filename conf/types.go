package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"github.com/rs/zerolog"
)

/*
createProfile is used to write out a sample configuration profile.
It is based on readConfig and will only include required fields rather than hidden, optional, or unsupported fields.
*/
type createProfile struct {
	Version  int    `yaml:"Version"`
	Provider string `yaml:"Provider"`
	AWS      struct {
		Profile string `yaml:"Profile"`
		Key     string `yaml:"Key"`
		Secret  string `yaml:"Secret"`
		ACL     string `yaml:"ACL"`
		Storage string `yaml:"Storage"`
	} `yaml:"AWS"`
	OCI struct {
		Profile         string `yaml:"Profile"`
		Compartment     string `yaml:"Compartment"`
		AuthTenancy     string `yaml:"AuthTenancy"`
		AuthUser        string `yaml:"AuthUser"`
		AuthRegion      string `yaml:"AuthRegion"`
		AuthFingerprint string `yaml:"AuthFingerprint"`
		AuthPrivateKey  string `yaml:"AuthPrivateKey"`
		AuthPassphrase  string `yaml:"AuthPassphrase"`
	} `yaml:"OCI"`
	Bucket struct {
		Create bool   `yaml:"Create"`
		Name   string `yaml:"Name"`
		Region string `yaml:"Region"`
	} `yaml:"Bucket"`
	Options struct {
		MaxUploads int    `yaml:"MaxUploads"`
		Overwrite  string `yaml:"Overwrite"`
	} `yaml:"Options"`
	Tagging struct {
		ChecksumSHA256 bool              `yaml:"Checksum"`
		Origins        bool              `yaml:"Origins"`
		Tags           map[string]string `yaml:"Tags"`
	} `yaml:"Tagging"`
	Objects struct {
		NamePrefix          string `yaml:"NamePrefix"`
		RootPrefix          string `yaml:"RootPrefix"`
		Naming              string `yaml:"Naming"`
		OmitOriginDirectory bool   `yaml:"OmitRootDir"`
	} `yaml:"Objects"`
	Logging struct {
		Level    int    `yaml:"Level"`
		Console  bool   `yaml:"Console"`
		File     bool   `yaml:"File"`
		Filepath string `yaml:"Filepath"`
	} `yaml:"Logging"`
	Uploads struct {
		Files       []string `yaml:"Files"`
		Directories []string `yaml:"Directories"`
	} `yaml:"Uploads"`
}

/*
readConfig is only used to unmarshal a YAML profile, it is not used in the application.
*/
type readConfig struct {
	// Version will be used for feature support
	Version  int    `yaml:"Version"`
	Provider string `yaml:"Provider"`

	// AWS will contain only AWS specific configuration details. Other providers will have their own
	// struct and fields.
	AWS struct {
		Profile string `yaml:"Profile"`
		Key     string `yaml:"Key"`
		Secret  string `yaml:"Secret"`
		ACL     string `yaml:"ACL"`
		Storage string `yaml:"Storage"`
	} `yaml:"AWS"`

	// OCI will contain only OCI specific configuration details. Other providers will have their own
	OCI struct {
		Profile         string `yaml:"Profile"`
		Compartment     string `yaml:"Compartment"`
		AuthTenancy     string `yaml:"AuthTenancy"`
		AuthUser        string `yaml:"AuthUser"`
		AuthRegion      string `yaml:"AuthRegion"`
		AuthFingerprint string `yaml:"AuthFingerprint"`
		AuthPrivateKey  string `yaml:"AuthPrivateKey"`
		AuthPassphrase  string `yaml:"AuthPassphrase"`
	} `yaml:"OCI"`

	// Bucket should be universal across all providers, though there may be different fields depending on the
	// provider.
	Bucket struct {
		Create bool   `yaml:"Create"`
		Name   string `yaml:"Name"`
		Region string `yaml:"Region"`
	} `yaml:"Bucket"`

	// The Objects struct contains object level configuration details, mostly related to object naming
	// Note: Object tags will be handled in the Tagging struct
	Objects struct {
		NamePrefix  string `yaml:"NamePrefix"`
		RootPrefix  string `yaml:"RootPrefix"`
		Naming      string `yaml:"Naming"`
		OmitRootDir bool   `yaml:"OmitRootDir"`
	} `yaml:"Objects"`

	// The Options struct is used to configure the application and how it operates.
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

	// The Uploads struct contains the list of files, folders, and directories to upload.
	// Folders and Directories will be merged.
	Uploads struct {
		Files       []string `yaml:"Files"`
		Folders     []string `yaml:"Folders"`
		Directories []string `yaml:"Directories"`
	} `yaml:"Uploads"`

	// Logging is used to configure the logging output, which is handled by the 'logbot' package.
	Logging struct {
		Level    int    `yaml:"Level"`
		Console  bool   `yaml:"Console"`
		File     bool   `yaml:"File"`
		Filepath string `yaml:"Filepath"`
	} `yaml:"Logging"`

	// Log is an instance of the logger.
	Log *logbot.LogBot
}

// Provider represents the configuration for a provider.
//
// Fields:
// - Is (ProviderName): The name of the provider. (e.g., "AWS", "OCI")
// - AWS (*ProviderAWS): The configuration for AWS.
// - OCI (*ProviderOCI): The configuration for OCI.
// - Key (string): The provider key.
// - Secret (string): The provider secret.
//
// Usage examples can be found in the surrounding code.
type Provider struct {
	Is     ProviderName
	AWS    *ProviderAWS
	OCI    *ProviderOCI
	Key    string
	Secret string
}

// ProviderAWS represents the AWS provider configuration.
//
// Fields:
// - Profile: The profile name used for authentication.
// - ACL: The access control list for the storage objects.
// - Storage: The storage class for the objects.
// - Key: The AWS access key ID.
// - Secret: The AWS secret access key.
type ProviderAWS struct {
	Profile string
	ACL     types.ObjectCannedACL
	Storage types.StorageClass
	Key     string
	Secret  string
}

// ProviderOCI represents the OCI provider configuration.
type ProviderOCI struct {
	Profile     string
	Compartment string
	Builder     *ProviderOCIBuilder
}

// ProviderOCIBuilder is a struct that holds the required attributes
// for building an OCI provider. These attributes include the tenancy,
// user, region, fingerprint, private key, and passphrase.
//
// Example usage:
//
//	ociBuilder := &ProviderOCIBuilder{
//	  Tenancy:     "my-tenancy",
//	  User:        "my-user",
//	  PrivateKey:  "my-private-key",
//	  Passphrase:  "my-passphrase",
//	  Fingerprint: "my-fingerprint",
//	  Region:      "my-region",
//	}
type ProviderOCIBuilder struct {
	Tenancy     string
	User        string
	Region      string
	Fingerprint string
	PrivateKey  string
	Passphrase  string
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

// LogOpts contains the logging configuration, but not an instance of logbot.
type LogOpts struct {
	Level    zerolog.Level
	Console  bool
	File     bool
	Filepath string
}

// AppConfig is an application level struct that all profile configuration details are loaded into.
// Log is an instance zerolog.Logger, which is built in logbot.
// The Files and Directories structs are the list of files and directories to upload.
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
type Naming string

const (
	NamingRelative Naming = "relative"
	NamingAbsolute Naming = "absolute"
)

// String returns the string representation of the Naming object.
// It converts the Naming object to a string by using the underlying string value.
func (n Naming) String() string {
	return string(n)
}

// S is a shortcut for fmt.Sprintf. The only real purpose is to reduce clutter and line lengths.
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

	ErrorProfilePath    = "error determining profile path"
	ErrorOpeningProfile = "error opening profile"
	ErrorReadingYaml    = "error reading yaml"

	ErrorLoggingFilepathNotSpecified = "path to log file not specified"
	ErrorLoggingFilepath             = "error determining log file path"
	ErrorLoggingLevelTooHigh         = "logging level too high, setting to 5"
	ErrorLoggingLevelTooLow          = "logging level too low, setting to 0"
	ErrorGettingFileInfo             = "error getting file info"
	ErrorFileIsDirectory             = "listed file is a directory"
	ErrorNoFilesSpecified            = "no files, folders, directories specified"
	ErrorNoReadableFiles             = "no readable files or directories specified"
	ErrorUnsupportedProfileVersion   = "profile version not supported"
	ErrorProviderNotSpecified        = "provider not specified"
	ErrorBucketNotSpecified          = "bucket or region not specified"
)
