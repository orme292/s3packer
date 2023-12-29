package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/logbot"
	"github.com/rs/zerolog"
)

type readConfig struct {
	Version int `yaml:"Version"`
	AWS     struct {
		Profile string `yaml:"Profile"`
		Key     string `yaml:"Key"`
		Secret  string `yaml:"Secret"`
		ACL     string `yaml:"ACL"`
		Storage string `yaml:"Storage"`
	} `yaml:"AWS"`
	Bucket struct {
		Create bool   `yaml:"Create"`
		Name   string `yaml:"Name"`
		Region string `yaml:"Region"`
	} `yaml:"Bucket"`
	Objects struct {
		NamePrefix          string `yaml:"NamePrefix"`
		RootPrefix          string `yaml:"RootPrefix"`
		Naming              string `yaml:"Naming"`
		OmitOriginDirectory bool   `yaml:"OmitOriginDirectory"`
	} `yaml:"Objects"`
	Options struct {
		MaxParts   int    `yaml:"MaxParts"`
		MaxUploads int    `yaml:"MaxUploads"`
		Overwrite  string `yaml:"Overwrite"`
	} `yaml:"Options"`
	Tagging struct {
		ChecksumSHA256 bool              `yaml:"Checksum"`
		Origins        bool              `yaml:"Origins"`
		Tags           map[string]string `yaml:"Tags"`
	} `yaml:"Tagging"`
	Uploads struct {
		Files       []string `yaml:"Files"`
		Folders     []string `yaml:"Folders"`
		Directories []string `yaml:"Directories"`
	} `yaml:"Uploads"`
	Logging struct {
		Level    int    `yaml:"Level"`
		Console  bool   `yaml:"Console"`
		File     bool   `yaml:"File"`
		Filepath string `yaml:"Filepath"`
	} `yaml:"Logging"`
	Log *logbot.LogBot
}

type Provider struct {
	Is         ProviderName
	AwsProfile string
	AwsACL     types.ObjectCannedACL
	AwsStorage types.StorageClass
	Key        string
	Secret     string
}

type Bucket struct {
	Create bool
	Name   string
	Region string
}

type Objects struct {
	NamePrefix          string
	RootPrefix          string
	Naming              Naming
	OmitOriginDirectory bool
}

type Opts struct {
	MaxParts   int
	MaxUploads int
	Overwrite  Overwrite
}

type Tag struct {
	ChecksumSHA256       bool
	AwsChecksumAlgorithm types.ChecksumAlgorithm
	AwsChecksumMode      types.ChecksumMode
	Origins              bool
}

type LogOpts struct {
	Level    zerolog.Level
	Console  bool
	File     bool
	Filepath string
}

type AppConfig struct {
	Provider    *Provider
	Bucket      *Bucket
	Objects     *Objects
	Opts        *Opts
	Tags        map[string]string
	Tag         *Tag
	LogOpts     *LogOpts
	Log         *logbot.LogBot
	Files       []string
	Directories []string
}

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

// Overwrite type
type Overwrite string

const (
	OverwriteChecksum Overwrite = "checksum"
	OverwriteNever    Overwrite = "never"
	OverwriteAlways   Overwrite = "always"
)

func (o Overwrite) String() string {
	return string(o)
}

// Naming type
type Naming string

const (
	NamingRelative Naming = "relative"
	NamingAbsolute Naming = "absolute"
)

func (n Naming) String() string {
	return string(n)
}

const (
	Empty = ""
)

func S(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
