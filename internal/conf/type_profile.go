package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ProfileIncoming struct {
	Version int `yaml:"Version"`

	Provider struct {
		Use     string `yaml:"Use"`
		Profile string `yaml:"Profile"`
		Key     string `yaml:"Key"`
		Secret  string `yaml:"Secret"`
	} `yaml:"Provider"`

	AWS struct {
		ACL     string `yaml:"ACL"`
		Storage string `yaml:"Storage"`
	} `yaml:"AWS"`

	Google struct {
		Project      string `yaml:"Project"`
		LocationType string `yaml:"LocationType"`
		Storage      string `yaml:"Storage"`
		BucketACL    string `yaml:"BucketACL"`
		ObjectACL    string `yaml:"ObjectACL"`
		ADC          string
	} `yaml:"Google"`

	Linode struct {
		Region string `yaml:"Region"`
	} `yaml:"Linode"`

	OCI struct {
		Compartment string `yaml:"Compartment"`
		Storage     string `yaml:"Storage"`
	} `yaml:"OCI"`

	Bucket struct {
		Create bool   `yaml:"Create"`
		Name   string `yaml:"Name"`
		Region string `yaml:"Region"`
	} `yaml:"Bucket"`

	Options struct {
		MaxUploads       int    `yaml:"MaxUploads"`
		FollowSymlinks   bool   `yaml:"FollowSymlinks"`
		WalkDirs         bool   `yaml:"WalkDirs"`
		OverwriteObjects string `yaml:"OverwriteObjects"`
	} `yaml:"Options"`

	TagOptions struct {
		OriginPath     bool `yaml:"OriginPath"`
		ChecksumSHA256 bool `yaml:"ChecksumSHA256"`
	} `yaml:"Tagging"`

	Tags map[string]string `yaml:"Tags"`

	Objects struct {
		NamingType  string `yaml:"NamingType"`
		NamePrefix  string `yaml:"NamePrefix"`
		PathPrefix  string `yaml:"PathPrefix"`
		OmitRootDir bool   `yaml:"OmitRootDir"`
	} `yaml:"Objects"`

	Logging struct {
		Level   int    `yaml:"Level"`
		Screen  bool   `yaml:"Screen"`
		Console bool   `yaml:"Console"`
		File    bool   `yaml:"File"`
		Logfile string `yaml:"Logfile"`
	} `yaml:"Logging"`

	Files []string `yaml:"Files"`
	Dirs  []string `yaml:"Dirs"`
	Skip  []string `yaml:"Skip"` // TODO: Add Support
}

func NewProfile() *ProfileIncoming {
	return &ProfileIncoming{}
}

func (p *ProfileIncoming) LoadFromYaml(filename string) error {

	filename, err := filepath.Abs(filename)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrorProfilePath, err)
	}

	f, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrorOpeningProfile, err)
	}

	err = yaml.Unmarshal(f, p)
	if err != nil {
		return fmt.Errorf("%s: %v", ErrorReadingYaml, err)
	}

	return nil

}

func (p *ProfileIncoming) loadSampleData() {

	p.Provider.Use = "aws"
	p.Provider.Profile = "myAwsProfile"
	p.Provider.Key = "key_value"
	p.Provider.Secret = "secret_value"

	p.AWS.ACL = "private"
	p.AWS.Storage = "intelligent_tiering"

	p.Linode.Region = "us-lax-1"

	p.Google.Project = "s3packer"
	p.Google.LocationType = "region"
	p.Google.Storage = "standard"
	p.Google.BucketACL = "private"
	p.Google.ObjectACL = "private"

	p.OCI.Compartment = "ocid1.compartment.oc1..aaaaaaaaa2qfwzyec6js1ua2ybtyyh3m39ze"
	p.OCI.Storage = "standard"

	p.Bucket.Create = true
	p.Bucket.Region = "us-lax-1"
	p.Bucket.Name = "MyBackupBucket"

	p.Options.MaxUploads = 50
	p.Options.FollowSymlinks = false
	p.Options.WalkDirs = true
	p.Options.OverwriteObjects = "never"

	p.TagOptions.OriginPath = true
	p.TagOptions.ChecksumSHA256 = true

	p.Tags = map[string]string{
		"Author": "Forrest Gump",
		"Title":  "Letters to Jenny",
	}

	p.Objects.NamingType = "relative"
	p.Objects.NamePrefix = "backup-"
	p.Objects.PathPrefix = "/backups/april/2023"
	p.Objects.OmitRootDir = true

	p.Logging.Level = 4
	p.Logging.Screen = false
	p.Logging.Console = true
	p.Logging.File = true
	p.Logging.Logfile = "/var/log/s3p.log"

	p.Files = []string{
		"/documents/to_jenny/letter_1.doc",
		"/documents/to_jenny/letter_2.doc",
		"/documents/to_jenny/letter_3.doc",
	}
	p.Dirs = []string{
		"/documents/from_jenny",
		"/documents/stock_certificates",
	}
	p.Skip = []string{
		"/documents/to_jenny/unsent",
	}

}
