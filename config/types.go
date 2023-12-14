package config

import (
	"github.com/orme292/s3packer/logbot"
)

/*
Configuration is the main configuration object for the s3packer application.

Authentication is the AWS key and access token.
Bucket is the bucket name and region, etc.
Options are the options for the s3 upload.
Dirs are the directories to upload.
Files are the files to upload.
Logging is the logging configuration.
Logger is the logbot object.
*/
type Configuration struct {
	Version        int               `yaml:"Version"`
	Authentication map[string]any    `yaml:"Authentication"`
	Bucket         map[string]any    `yaml:"Bucket"`
	Options        map[string]any    `yaml:"Options"`
	Naming         map[string]any    `yaml:"Naming"`
	Dirs           []string          `yaml:"Dirs"`
	Files          []string          `yaml:"Files"`
	Tags           map[string]string `yaml:"Tags"`
	Logging        map[string]any    `yaml:"Logging"`
	Logger         logbot.LogBot
}

const (
	EmptyString = ""

	ACLPrivate                = "private"
	ACLPublicRead             = "public-read"
	ACLPublicReadWrite        = "public-read-write"
	ACLAuthenticatedRead      = "authenticated-read"
	ACLAwsExecRead            = "aws-exec-read"
	ACLBucketOwnerRead        = "bucket-owner-read"
	ACLBucketOwnerFullControl = "bucket-owner-full-control"
	ACLLogDeliveryWrite       = "log-delivery-write"

	NameMethodRelative = "relative"
	NameMethodAbsolute = "absolute"

	StorageClassStandard           = "STANDARD"
	StorageClassStandardIA         = "STANDARD_IA"
	StorageClassOneZoneIA          = "ONEZONE_IA"
	StorageClassIntelligentTiering = "INTELLIGENT_TIERING"
	StorageClassGlacier            = "GLACIER"
	StorageClassDeepArchive        = "DEEP_ARCHIVE"
)

const (
	ProfileAuthProfile           = "aws-profile"
	ProfileAuthKey               = "key"
	ProfileAuthSecret            = "secret"
	ProfileBucketName            = "name"
	ProfileBucketRegion          = "region"
	ProfileOptionACL             = "acl"
	ProfileOptionObjectPrefix    = "objectPrefix"
	ProfileOptionPathPrefix      = "pathPrefix"
	ProfileOptionsMaxConcurrent  = "maxConcurrentUploads"
	ProfileOptionOverwrite       = "overwrite"
	ProfileOptionStorage         = "storage"
	ProfileOptionTagOrigins      = "tagOrigins"
	ProfileOptionKeyNamingMethod = "keyNamingMethod"
	ProfileOptionOmitOriginDir   = "keyNamingOmitOriginDir" // Hidden Option
	ProfileLoggingToConsole      = "toConsole"
	ProfileLoggingToFile         = "toFile"
	ProfileLoggingFilename       = "filename"
	ProfileLoggingLevel          = "level"
)
