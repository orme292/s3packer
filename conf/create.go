package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Create(filename string) (err error) {
	type profile struct {
		Version int `yaml:"Version"`
		AWS     struct {
			Profile string `yaml:"Profile"`
			Key     string `yaml:"Key"`
			Secret  string `yaml:"Secret"`
			ACL     string `yaml:"ACL"`
			Storage string `yaml:"Storage"`
		} `yaml:"AWS"`
		Bucket struct {
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
			OmitOriginDirectory bool   `yaml:"OmitOriginDirectory"`
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

	r := profile{}
	r.Version = 2
	r.AWS.Profile = "default"
	r.AWS.Key = ""
	r.AWS.Secret = ""
	r.AWS.ACL = "private"
	r.AWS.Storage = "standard"
	r.Bucket.Name = "my-bucket"
	r.Bucket.Region = "us-east-1"
	r.Options.MaxUploads = 10
	r.Options.Overwrite = "never"
	r.Tagging.ChecksumSHA256 = true
	r.Tagging.Origins = true
	r.Tagging.Tags = map[string]string{
		"tag1": "value1",
	}
	r.Objects.NamePrefix = ""
	r.Objects.RootPrefix = ""
	r.Objects.Naming = "absolute"
	r.Logging.Level = 2
	r.Logging.Console = true
	r.Logging.File = false
	r.Logging.Filepath = ""
	r.Uploads.Files = []string{
		"file1.txt",
		"file2.txt",
	}
	r.Uploads.Directories = []string{
		"/home/me/dir1",
		"/home/me/dir2",
	}

	o, err := yaml.Marshal(&r)
	if err != nil {
		return err
	}

	filename, err = filepath.Abs(filepath.Clean(filename))
	if err != nil {
		return err
	}

	_, err = os.Stat(filename)
	if !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("Error closing file: %q\n", err.Error())
			os.Exit(1)
		}
	}(f)

	_, err = f.Write(o)
	if err != nil {
		return err
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(o))
	return nil
}
