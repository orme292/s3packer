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
	InvalidAWSACL          = "invalid aws acl"
	ErrorAWSProfileAndKeys = "you provided AWS profile name and key/secret pair, use profile or keys"
	ErrorAWSMissingAuth    = "must provide a valid AWS key pair"
	ErrorAWSAuthNeeded     = "must provide AWS profile name or key pair"
)

// Google Cloud Constants
const (
	GCBucketACLAuthenticatedRead = "authenticatedread"
	GCBucketACLPrivate           = "private"
	GCBucketACLPublicRead        = "publicread"
	GCBucketACLPublicReadWrite   = "publicreadwrite"
	GCBucketACLProjectPrivate    = "projectprivate"
)

const (
	GCObjectACLAuthenticatedRead = "authenticatedread"
	GCObjectACLPrivate           = "private"
	GCObjectACLPublicRead        = "publicread"
	GCObjectACLProjectPrivate    = "projectprivate"
	GCObjectACLBucketOwnerFull   = "bucketownerfullcontrol"
	GCObjectACLBucketOwnerRead   = "bucketownerread"
)

const (
	GCStorageStandard = "STANDARD"
	GCStorageNearline = "NEARLINE"
	GCStorageColdline = "COLDLINE"
	GCStorageArchive  = "ARCHIVE"
)

const (
	GCLocationTypeRegion = "region"
	GCLocationTypeDual   = "dual-region"
	GCLocationTypeMulti  = "multi-region"
)

const (
	GCRegionNANE1    = "NORTHAMERICA-NORTHEAST1"
	GCRegionNANE2    = "NORTHAMERICA-NORTHEAST2"
	GCRegionNAS1     = "NORTHAMERICA-SOUTH1"
	GCRegionUSC1     = "US-CENTRAL1"
	GCRegionUSE1     = "US-EAST1"
	GCRegionUSE4     = "US-EAST4"
	GCRegionUSE5     = "US-EAST5"
	GCRegionUSS1     = "US-SOUTH1"
	GCRegionUSW1     = "US-WEST1"
	GCRegionUSW2     = "US-WEST2"
	GCRegionUSW3     = "US-WEST3"
	GCRegionUSW4     = "US-WEST4"
	GCRegionSAE1     = "SOUTHAMERICA-EAST1"
	GCRegionSAW1     = "SOUTHAMERICA-WEST1"
	GCRegionEUC2     = "EUROPE-CENTRAL2"
	GCRegionEUN1     = "EUROPE-NORTH1"
	GCRegionEUSW1    = "EUROPE-SOUTHWEST1"
	GCRegionEUW1     = "EUROPE-WEST1"
	GCRegionEUW2     = "EUROPE-WEST2"
	GCRegionEUW3     = "EUROPE-WEST3"
	GCRegionEUW4     = "EUROPE-WEST4"
	GCRegionEUW6     = "EUROPE-WEST6"
	GCRegionEUW8     = "EUROPE-WEST8"
	GCRegionEUW9     = "EUROPE-WEST9"
	GCRegionEUW12    = "EUROPE-WEST12"
	GCRegionASIAE1   = "ASIA-EAST1"
	GCRegionASIAE2   = "ASIA-EAST2"
	GCRegionASIANE1  = "ASIA-NORTHEAST1"
	GCRegionASIANE2  = "ASIA-NORTHEAST2"
	GCRegionASIANE3  = "ASIA-NORTHEAST3"
	GCRegionASIASE1  = "ASIA-SOUTHEAST1"
	GCRegionASIAS1   = "ASIA-SOUTH1"
	GCRegionASIAS2   = "ASIA-SOUTH2"
	GCRegionASIASE2  = "ASIA-SOUTHEAST2"
	GCRegionMEC1     = "MIDDLEEAST-CENTRAL1"
	GCRegionMEC2     = "MIDDLEEAST-CENTRAL2"
	GCRegionMEW1     = "MIDDLEEAST-WEST1"
	GCRegionAUSSE1   = "AUSTRALIA-SOUTHEAST1"
	GCRegionAUSSE2   = "AUSTRALIA-SOUTHEAST2"
	GCRegionAFRICAS1 = "AFRICA-SOUTH1"
)

const (
	GCDualAsia = "ASIA1"
	GCDualEur4 = "EUR4"
	GCDualEur5 = "EUR5"
	GCDualEur7 = "EUR7"
	GCDualEur8 = "EUR8"
	GCDualNAM4 = "NAM4"
)

const (
	GCMultiAsia = "ASIA"
	GCMultiEU   = "EU"
	GCMultiUS   = "US"
)

const (
	InvalidGCBucketACL    = "invalid google cloud bucket acl"
	InvalidGCObjectACL    = "invalid google cloud object acl"
	InvalidGCLocationType = "invalid google cloud location type"
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
	LinodeInvalidRegion = "invalid Linode region provided"
	LinodeAuthNeeded    = "Linode authentication not specified"
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
