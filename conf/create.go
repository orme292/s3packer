package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// outputProfile is used to write out a sample configuration profile. It is based on readConfig, but only includes
// necessary fields. This prevents any hidden or unsupported fields from being revealed.
type outputProfile struct {
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

// Create takes a filename as a string, and writes out a sample configuration profile. The file must not exist.
// The structure is built using a new struct, outputProfile, that is based on readConfig
func Create(filename string) (err error) {
	r := outputProfile{}
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
		"hostname": "this host",
		"author":   "me",
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

	ok, err := canCreate(filename)
	if !ok {
		return fmt.Errorf("cannot create file %s: %s", filename, err.Error())
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Printf("error closing file: %q\n", err.Error())
			os.Exit(1)
		}
	}(f)

	_, err = f.Write(o)
	if err != nil {
		return err
	}

	fmt.Printf("--- Writing:\n%s\n\n", string(o))
	fmt.Printf("Wrote new profile to %q\n", filename)
	return nil
}

// canCreate checks whether a file can be created. It returns true if the file does not exist, and false if it does
// or if another error occurs. To figure it out if the program has permissions ot create the file, it attempts to
// create the file. If creation succeeds, then the file is immediately removed.
func canCreate(f string) (bool, error) {
	filename, err := filepath.Abs(filepath.Clean(f))
	if err != nil {
		return false, err
	}

	// Resolve G304: Potential file inclusion via variable
	if strings.Contains(filename, "..") {
		return false, fmt.Errorf("invalid filename: %s", filename)
	}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o640)
			if err != nil {
				return false, err
			}
			err = file.Close()
			if err != nil {
				return false, err
			}
			err = os.Remove(filename)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	return false, fmt.Errorf("file %s already exists", filename)
}
