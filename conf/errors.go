package conf

const (
	InvalidACL             = "invalid acl"
	InvalidNamingMethod    = "invalid object naming method"
	InvalidStorageClass    = "invalid storage class"
	InvalidOverwriteMethod = "invalid overwrite method"
	InvalidTagChars        = "invalid characters removed from tag"
)

const (
	ErrorProfilePath                 = "error determining profile path"
	ErrorOpeningProfile              = "error opening profile"
	ErrorReadingYaml                 = "error reading yaml"
	ErrorAWSProfileAndKeys           = "both aws profile and keys are specified, use profile or keys"
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
