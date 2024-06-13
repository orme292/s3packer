package conf

const (
	Empty = ""
)

// AWS Constants
const (
	AwsACLPrivate                = "private"
	AwsACLPublicRead             = "public-read"
	AwsACLPublicReadWrite        = "public-read-write"
	AwsACLAuthenticatedRead      = "authenticated-read"
	AwsACLAwsExecRead            = "aws-exec-read"
	AwsACLBucketOwnerRead        = "bucket-owner-read"
	AwsACLBucketOwnerFullControl = "bucket-owner-full-control"
)

const (
	AwsClassStandard           = "STANDARD"
	AwsClassReducedRedundancy  = "REDUCED_REDUNDANCY"
	AwsClassGlacierIR          = "GLACIER_IR"
	AwsClassSnow               = "SNOW"
	AwsClassStandardIA         = "STANDARD_IA"
	AwsClassOneZoneIA          = "ONEZONE_IA"
	AwsClassIntelligentTiering = "INTELLIGENT_TIERING"
	AwsClassGlacier            = "GLACIER"
	AwsClassDeepArchive        = "DEEP_ARCHIVE"
)

const (
	InvalidAWSACL                   = "invalid aws acl"
	ErrorAWSProfileAndKeys          = "you provided AWS profile name and key/secret pair, use profile or keys"
	ErrorAWSKeyOrSecretNotSpecified = "must specify both AWS key and secret"
	ErrorAWSAuthNeeded              = "must provider either AWS profile name or key/secret pair"
)

// OCI Constants
const (
	OciDefaultProfile = "DEFAULT"
)

const (
	ErrorOCICompartmentNotSpecified = "OCI compartment will be tenancy root"
	ErrorOCIAuthNotSpecified        = "OCI auth not specified"
	ErrorOCIStorageNotSpecified     = "OCI storage tier is not valid"
)

const (
	OracleStorageTierStandard         = "standard"
	OracleStorageTierInfrequentAccess = "infrequentaccess" // the case is strange because of the
	OracleStorageTierArchive          = "archive"
)

// Linode Constants
const (
	LinodeClusterAmsterdam  = "nl-ams-1.linodeobjects.com"
	LinodeClusterAtlanta    = "us-southeast-1.linodeobjects.com"
	LinodeClusterChennai    = "in-maa-1.linodeobjects.com"
	LinodeClusterChicago    = "us-ord-1.linodeobjects.com"
	LinodeClusterFrankfurt  = "eu-central-1.linodeobjects.com"
	LinodeClusterJakarta    = "id-cgk-1.linodeobjects.com"
	LinodeClusterLosAngeles = "us-lax-1.linodeobjects.com"
	LinodeClusterMiami      = "us-mia-1.linodeobjects.com"
	LinodeClusterMilan      = "it-mil-1.linodeobjects.com"
	LinodeClusterNewark     = "us-east-1.linodeobjects.com"
	LinodeClusterOsaka      = "jp-osa-1.linodeobjects.com"
	LinodeClusterParis      = "fr-par-1.linodeobjects.com"
	LinodeClusterSaoPaulo   = "br-gru-1.linodeobjects.com"
	LinodeClusterSeattle    = "us-sea-1.linodeobjects.com"
	LinodeClusterSingapore  = "ap-south-1.linodeobjects.com"
	LinodeClusterStockholm  = "se-sto-1.linodeobjects.com"
	LinodeClusterAshburn    = "us-iad-1.linodeobjects.com"
)

const (
	LinodeRegionAmsterdam  = "nl-ams-1"
	LinodeRegionAtlanta    = "us-southeast-1"
	LinodeRegionChennai    = "in-maa-1"
	LinodeRegionChicago    = "us-ord-1"
	LinodeRegionFrankfurt  = "eu-central-1"
	LinodeRegionJakarta    = "id-cgk-1"
	LinodeRegionLosAngeles = "us-lax-1"
	LinodeRegionMiami      = "us-mia-1"
	LinodeRegionMilan      = "it-mil-1"
	LinodeRegionNewark     = "us-east-1"
	LinodeRegionOsaka      = "jp-osa-1"
	LinodeRegionParis      = "fr-par-1"
	LinodeRegionSaoPaulo   = "br-gru-1"
	LinodeRegionSeattle    = "us-sea-1"
	LinodeRegionSingapore  = "ap-south-1"
	LinodeRegionStockholm  = "se-sto-1"
	LinodeRegionAshburn    = "us-iad-1"
)

const (
	LinodeInvalidRegion           = "invalid Linode region provided"
	LinodeKeyOrSecretNotSpecified = "Linode access keys not specified"
)

// Conf Errors
const (
	InvalidNamingType      = "NamingType should be \"relative\" or \"absolute\""
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
	ErrorNoFilesSpecified            = "no files or directories specified"
	ErrorNoReadableFiles             = "no readable files or directories specified"
	ErrorUnsupportedProfileVersion   = "profile version not supported"
	ErrorProviderNotSpecified        = "no valid provider specified"
	ErrorBucketInfoNotSpecified      = "bucket name or bucket region not specified"
)
