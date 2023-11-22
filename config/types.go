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
Count is used to count the number of objects uploaded
*/
type Configuration struct {
	Authentication map[string]any `yaml:"Authentication"`
	Bucket         map[string]any `yaml:"Bucket"`
	Options        map[string]any `yaml:"Options"`
	Dirs           []string       `yaml:"Dirs"`
	Files          []string       `yaml:"Files"`
	Logging        map[string]any `yaml:"Logging"`
	Logger         logbot.LogBot
}
