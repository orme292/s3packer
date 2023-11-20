package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/orme292/s3packer/logbot"
	"gopkg.in/yaml.v2"
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
	Authentication map[string]string `yaml:"Authentication"`
	Bucket         map[string]string `yaml:"Bucket"`
	Options        map[string]any    `yaml:"Options"`
	Dirs           []string          `yaml:"Dirs"`
	Files          []string          `yaml:"Files"`
	Logging        map[string]any    `yaml:"Logging"`
	Logger         logbot.LogBot
}

/*
New returns a new Configuration object with default values for logging.
*/
func New() Configuration {
	return Configuration{
		Logger: logbot.LogBot{
			Level:       logbot.DEBUG,
			FlagConsole: true,
			FlagFile:    false,
			Path:        "log.log",
		},
	}
}

/*
Load loads and sanitizes the configuration from the specified file.

The file is expected to be in YAML format.
1. Load the file.
2. Unmarshal the YAML into the Configuration object.
3. Configure logging.
4. Process Authentication.
5. Process Bucket details.
6. Process Options.
7. Return the Configuration object.
*/
func (c *Configuration) Load(file string) error {
	f, err := os.ReadFile(file)
	if err != nil {
		return errors.New(err.Error())
	}

	fmt.Println("Using profile:", file)

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return errors.New(err.Error())
	}

	// Configure logging first
	if c.Logging["console"] == "" {
		c.Logger.FlagConsole = true
	} else {
		c.Logger.FlagConsole = c.Logging["console"].(bool)
	}
	if c.Logging["file"] == "" {
		c.Logger.FlagFile = false
	} else {
		c.Logger.FlagFile = c.Logging["file"].(bool)
	}
	if c.Logging["file"] == true && c.Logging["path"] == "" {
		return errors.New("log to file is true but path not provided")
	}
	if c.Logging["path"] != "" {
		path, err := filepath.Abs(c.Logging["path"].(string))
		if err != nil {
			return errors.New("unable to get absolute path of log file: " + err.Error())
		}
		c.Logger.Path = path
	}

	// Process Authentication
	if c.Authentication["key"] == "" {
		return errors.New("authentication key is empty")
	}
	if c.Authentication["access"] == "" {
		return errors.New("authentication access token is empty")
	}

	// Process Bucket details
	if c.Bucket["name"] == "" {
		return errors.New("bucket name is empty")
	}
	if c.Bucket["region"] == "" {
		return errors.New("bucket region is empty")
	}

	// Process Options
	/*
		AWS Canned ACLs
		c.Options["acl"] takes an aws canned ACL value:
		- private: Owner gets FULL_CONTROL. No one else has access rights (default).
		- public-read: Owner gets FULL_CONTROL. The AllUsers group gets READ access.
		- public-read-write: Owner gets FULL_CONTROL. The AllUsers group gets READ and WRITE access.
		- authenticated-read: Owner gets FULL_CONTROL. Amazon EC2 gets READ access to GET an AMI bundle.
		- aws-exec-read: Owner gets FULL_CONTROL. The AuthenticatedUsers group gets READ access.
		- bucket-owner-read: Object owner gets FULL_CONTROL. Bucket owner gets READ access.
		- bucket-owner-full-control: Both the object owner and the bucket owner get FULL_CONTROL over the object.
		- log-delivery-write: The LogDelivery group gets WRITE and READ_ACP permissions on the bucket.

		See: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl for more information.
	*/
	c.Options["acl"] = strings.ToLower(strings.TrimSpace(c.Options["acl"].(string)))
	if c.Options["acl"] == "" {
		c.Options["acl"] = "private"
	} else {
		switch strings.ToLower(c.Options["acl"].(string)) {
		case "private":
		case "public-read":
		case "public-read-write":
		case "authenticated-read":
		case "aws-exec-read":
		case "bucket-owner-read":
		case "bucket-owner-full-control":
		case "log-delivery-write":
		default:
			c.Logger.Warn("Invalid ACL specified, using default.")
			c.Options["acl"] = "private"
		}
		c.Logger.Info("ACL set to " + c.Options["acl"].(string))
	}
	if c.Options["overwrite"] == "" {
		c.Options["overwrite"] = false
	}
	c.Options["prefix"] = strings.TrimSpace(c.Options["prefix"].(string))

	/*
		AWS S3 Storage Classes
		c.Options["storage"] takes an aws storage class:
		- STANDARD: Standard storage class for frequently accessed data.
		- STANDARD_IA: Standard-Infrequent Access storage class for infrequently accessed data.
		- ONEZONE_IA: One Zone-Infrequent Access storage class for infrequently accessed data that can be recreated.
		- INTELLIGENT_TIERING: Intelligent Tiering storage class for data with unknown or changing access patterns.
		- GLACIER: Glacier storage class for long-term data archival.
		- DEEP_ARCHIVE: Deep Archive storage class for long-term data archival with the lowest cost.
	*/
	c.Options["storage"] = strings.ToUpper(strings.TrimSpace(c.Options["storage"].(string)))
	if c.Options["storage"] == "" {
		c.Options["storage"] = "STANDARD"
	} else {
		switch c.Options["storage"].(string) {
		case "STANDARD":
		case "STANDARD_IA":
		case "ONEZONE_IA":
		case "INTELLIGENT_TIERING":
		case "GLACIER":
		case "DEEP_ARCHIVE":
		default:
			c.Logger.Warn("Invalid storage class specified, using default.")
			c.Options["storage"] = strings.ToUpper("STANDARD")
		}
		c.Logger.Info("Storage Class set to " + c.Options["storage"].(string))
	}
	if len(c.Dirs) == 0 && len(c.Files) == 0 {
		return errors.New("no files or directories specified")
	}

	return nil
}
