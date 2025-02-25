package conf

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

const testInvalidString = "iNvAlIdStRiNg"
const testSomeString = "someValue"
const testEmptyString = ""

func failMsg(name, expect, actual string, err ...error) string {
	msg := fmt.Sprintf("%s -- EXPECT: [%s] GOT: [%s]", name, expect, actual)
	if len(err) != 0 {
		for i := range err {
			msg = fmt.Sprintf("%s ERROR: [%v]", msg, err[i])
		}
	}
	return msg
}

func camelCaseByChar(s string) string {
	s = strings.ToLower(s)
	var result []rune
	letterCount := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			if letterCount%2 == 0 {
				result = append(result, unicode.ToUpper(r))
			} else {
				result = append(result, r)
			}
			letterCount++
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func createTestFile() (string, error) {
	file, err := os.CreateTemp("", "builder_test.yaml")
	if err != nil {
		log.Printf("Could not create temp file: %v\n", file.Name())
		return "", err
	}
	log.Println("Created test file ", file.Name())

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("Could not close temp file.")
		}
	}(file)

	err = writeTestFile(file)
	if err != nil {
		log.Println("Could not write to test file.")
		return "", err
	}

	return file.Name(), nil
}

func newIncomingProfile() *ProfileIncoming {
	profile := ProfileIncoming{
		Version: 6,
	}
	profile.loadSampleData()

	profile.Provider.Use = "aws"
	profile.Provider.Profile = "default"
	profile.Provider.Key = ""
	profile.Provider.Secret = ""
	profile.AWS.ACL = "private"
	profile.AWS.Storage = "standard"
	profile.Bucket.Name = "s3p_builder_test_bucket"
	profile.Bucket.Create = true
	profile.Bucket.Region = "us-east-1"
	profile.Options.MaxUploads = 10
	profile.Options.FollowSymlinks = true
	profile.Options.WalkDirs = true
	profile.Options.OverwriteObjects = "always"
	profile.TagOptions.OriginPath = true
	profile.TagOptions.ChecksumSHA256 = true
	profile.Tags = map[string]string{
		"test?": "yes",
	}
	profile.Objects.NamingType = "relative"
	profile.Objects.NamePrefix = "builder_test_"
	profile.Objects.PathPrefix = "gotest"
	profile.Objects.OmitRootDir = true
	profile.Logging.Level = 5
	profile.Logging.Screen = false
	profile.Logging.Console = true
	profile.Logging.File = false
	profile.Logging.Logfile = ""
	profile.Files = []string{
		".",
	}
	profile.Dirs = []string{
		"~",
	}

	return &profile
}

func writeTestFile(file *os.File) error {

	profile := newIncomingProfile()
	yamlData, err := yaml.Marshal(&profile)
	if err != nil {
		log.Println("Could not marshal test profile.")
		return err
	}

	n, err := file.WriteString("---\n")
	if err != nil || n == 0 {
		log.Println("Could not write YAML header to test file.")
		return err
	}

	n, err = file.Write(yamlData)
	if err != nil || n == 0 {
		log.Println("Could not write YAML data to test file.")
		return err
	}

	return nil
}
