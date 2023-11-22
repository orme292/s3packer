package s3pack

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/orme292/s3packer/config"
)

/*
ListedAppendPrefix takes the application configuration (c) and an array of strings (keys).
If the configuration option "prefix" (c.Options["prefix"] is set to an empty string, then the keys array is returned with nil error.
Otherwise, the prefix is prepended to each key in the keys array and a new array of strings (prefixed) is returned
with nil error.
*/
func ListedAppendPrefix(c config.Configuration, keys []string) (prefixed []string, err error) {
	if len(keys) == 0 {
		return nil, errors.New("no keys to prefix")
	}
	if c.Options["prefix"].(string) == "" {
		return keys, nil
	}
	for _, key := range keys {
		prefixed = append(prefixed, fmt.Sprintf("%s%s", c.Options["prefix"], filepath.Base(key)))
	}
	return prefixed, nil
}

/*
AppendPrefix takes the application configuration (c) and a key string (key) and returns a new string with the prefix
prepended to the key.
*/
func AppendPrefix(c config.Configuration, key string) string {
	if c.Options["prefix"] == "" {
		return key
	}
	return fmt.Sprintf("%s%s", c.Options["prefix"], filepath.Base(key))
}

/*
ListedLocalFileExists takes the application configuration (c) and an array of strings (files).
If the configuration option "overwrite" is set to true, then no changes are made and the files array
is returned with nil error.
Otherwise, we check each file in the files array to see if it exists. If it does, then we add it to a new array
of strings (exists). If it doesn't, then we log a warning and move on to the next file.
After all elements in the files array are checked, we return a new string array (exists) and a nil error.
*/
func ListedLocalFileExists(c config.Configuration, files []string) (exists []string, err error) {
	if len(files) == 0 {
		return nil, errors.New("no files to check")
	}
	for _, file := range files {
		file, _ = filepath.Abs(file)
		if _, err := os.Stat(file); err != nil {
			c.Logger.Warn(file + " not found, skipping")
		} else if !errors.Is(err, fs.ErrNotExist) {
			exists = append(exists, file)
		}
	}
	return exists, nil
}

/*
LocalFileExists takes a string (file) and returns a bool and an error.
The function checks to see if the file exists. If it does, then we return true and a nil error.
If it doesn't, then we return false and a nil error.
*/
func LocalFileExists(file string) (bool, error) {
	file, _ = filepath.Abs(file)
	if _, err := os.Stat(file); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

/*
ListedObjectsNotExisting takes the application configuration (c) and an array of strings, un-prefixed keys (objects).
If the configuration option "overwrite" is set to true, then no changes are made and the objects array
is returned with nil error.
Otherwise, we check each object key in the objects array to see if it exists in the bucket. If it does not, then
we add it to a new array of strings (remaining). If it does, then we log a warning and move on to the next object.
After all elements in the objects array are checked, we return a new string array (remaining) and a nil error.
*/
func ListedObjectsNotExisting(c config.Configuration, objects []string) (remaining []string, err error) {
	if len(objects) == 0 {
		return nil, errors.New("no objects to check")
	}
	// If overwrite is true, then we're overwriting objects that have a key with the same name
	// as the file being uploaded.
	if c.Options["overwrite"] == true {
		c.Logger.Warn("Overwrite is enabled, skipping check for existing objects.")
		return objects, nil
	}

	var awsError awserr.Error

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Bucket["region"].(string)),
		Credentials: credentials.NewStaticCredentials(c.Authentication["key"].(string), c.Authentication["secret"].(string), ""),
	})

	svc := s3.New(sess, &aws.Config{})

	for _, object := range objects {
		object = AppendPrefix(c, object)
		_, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(c.Bucket["name"].(string)),
			Key:    aws.String(object),
		})
		if err != nil {
			if errors.As(err, &awsError) {
				switch awsError.Code() {
				case s3.ErrCodeNoSuchBucket:
					return objects, errors.New("aws says bucket does not exist")
				case s3.ErrCodeNoSuchKey:
					remaining = append(remaining, object)
				default:
					return objects, errors.New("aws error isn't handled: " + awsError.Error())
				}
			}
		}
		c.Logger.Warn(fmt.Sprintf("Object with key \"%s\" exists, skipping", object))
	}
	return remaining, nil
}

/*
ObjectExists takes the application configuration (c) and a prefixed string (key).
The function checks to see if the object exists in the bucket. If it does, then we return true and a nil error.
If it doesn't, then we return false and a nil error.
*/
func ObjectExists(c config.Configuration, prefixedKey string) (bool, error) {
	if prefixedKey == "" {
		return false, errors.New("key is empty")
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Bucket["region"].(string)),
		Credentials: credentials.NewStaticCredentials(c.Authentication["key"].(string), c.Authentication["secret"].(string), ""),
	})

	svc := s3.New(sess, &aws.Config{})

	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Bucket["name"].(string)),
		Key:    aws.String(prefixedKey),
	})
	if err != nil {
		var awsError awserr.Error
		if errors.As(err, &awsError) {
			switch awsError.Code() {
			case s3.ErrCodeNoSuchBucket:
				c.Logger.Fatal("aws says bucket does not exist")
			case s3.ErrCodeNoSuchKey:
				return false, nil
			default:
				c.Logger.Fatal("aws error isn't handled: " + awsError.Error())
			}
		}

	}
	if out != nil {
		return true, nil
	}
	return false, errors.New("unable to determine if object exists")
}
