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
New returns a new Configuration object with default values for logging.
*/
func New() Configuration {
	return Configuration{
		Logger: logbot.LogBot{
			Level:       logbot.INFO,
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

	fmt.Println("Using profile", file)

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return errors.New(err.Error())
	}

	err = c.Validate()
	if err != nil {
		return err
	}

	c.Logger.FlagConsole = c.Logging["toConsole"].(bool)
	c.Logger.FlagFile = c.Logging["toFile"].(bool)
	if c.Logging["toFile"].(bool) == true && c.Logging["path"].(string) != "" {
		path, err := filepath.Abs(c.Logging["path"].(string))
		if err != nil {
			return errors.New("unable to get absolute path of log file: " + err.Error())
		}
		c.Logger.Path = path
	}
	c.Logger.SetLogLevel(c.Logger.ParseIntLevel(c.Logging["level"]))

	c.sanitizeACL()
	c.sanitizeStorageType()

	if c.Options["overwrite"] == "" {
		c.Options["overwrite"] = false
	}
	c.Options["prefix"] = strings.TrimSpace(c.Options["prefix"].(string))

	if c.Files == nil && c.Dirs == nil {
		return errors.New("no files or directories specified")
	}
	if len(c.Files) == 0 && len(c.Dirs) == 0 {
		return errors.New("no files or directories specified")
	}

	return nil
}
