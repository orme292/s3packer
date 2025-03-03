package provider

import (
	"os"
	"os/user"
	"path/filepath"
	"sort"

	"s3p/internal/conf"
)

func getHomeDir() string {
	var path string
	if path = os.Getenv("HOME"); path != "" {
		return path
	}

	usr, err := user.Current()
	if err == nil {
		return usr.HomeDir
	}
	return filepath.Join("home", os.Getenv("USER"))
}

func getFiveFiles() []string {
	homeDir := getHomeDir()
	if homeDir == "" {
		return nil
	}

	entries, err := os.ReadDir(homeDir)
	if err != nil {
		return nil
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(homeDir, entry.Name()))
		}
	}

	sort.Strings(files)

	if len(files) > 5 {
		return files[:5]
	}
	return files
}

func newIncomingProfile() *conf.ProfileIncoming {
	profile := conf.ProfileIncoming{
		Version: 6,
	}
	profile.Provider.Use = "aws"
	profile.Provider.Profile = "default"
	profile.Provider.Key = ""
	profile.Provider.Secret = ""
	profile.AWS.ACL = "private"
	profile.AWS.Storage = "standard"
	profile.Bucket.Name = "s3p_builder_test_bucket"
	profile.Bucket.Create = true
	profile.Bucket.Region = "us-east-1"
	profile.Options.MaxUploads = 50
	profile.Options.FollowSymlinks = true
	profile.Options.WalkDirs = false
	profile.Options.OverwriteObjects = "always"
	profile.TagOptions.OriginPath = true
	profile.TagOptions.ChecksumSHA256 = true
	profile.Tags = map[string]string{}
	profile.Objects.NamingType = "absolute"
	profile.Objects.NamePrefix = "s3p_test_"
	profile.Objects.PathPrefix = "gotest"
	profile.Objects.OmitRootDir = true
	profile.Logging.Level = 5
	profile.Logging.Screen = false
	profile.Logging.Console = true
	profile.Logging.File = false
	profile.Logging.Logfile = ""
	profile.Files = getFiveFiles()
	profile.Dirs = []string{
		getHomeDir(),
	}

	return &profile
}
