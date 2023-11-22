package s3pack

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/orme292/s3packer/config"
)

var uplCount int

/*
NewDirectoryIterator creates a new DirectoryIterator struct, which implements the s3manager.BatchUploadIterator
interface for use with s3manager.UploadWithIterator. It takes a config.Configuration struct and a directory
path as input. It returns a pointer to a DirectoryIterator struct.

It walks the specified directory (dir) and creates a list of files to upload, and is assigned
to DirectoryIterator.filePaths
*/
func NewDirectoryIterator(c config.Configuration, dir string) s3manager.BatchUploadIterator {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		c.Logger.Fatal(err.Error())
	}

	return &DirectoryIterator{
		filePaths: paths,
		bucket:    c.Bucket["name"].(string),
		config:    c,
	}
}

/*
Next returns true if there are more files to upload, or false if there are no more files to upload. It also
sets the DirectoryIterator.next struct for processing.

 1. Check if an object exists with the same name as the file.
    If it does, and the overwrite option is not set, we skip to the next file and recheck. Otherwise, we move forward.
 2. If there are no more files to upload, return false
 3. If there are more files to upload, open the file, fill in the DirectoryIterator.next struct for processing and
    return true.

DirectoryIterator.next will be used by UploadObject() to upload the file to S3.
*/
func (di *DirectoryIterator) Next() bool {
	for {
		if len(di.filePaths) == 0 {
			di.next.f = nil
			return false
		}
		// If overwrite is true, just break out of the loop and continue.
		if di.config.Options["overwrite"].(bool) == true {
			break
		}
		exists, err := ObjectExists(di.config, AppendPrefix(di.config, di.filePaths[0]))
		if err != nil {
			di.config.Logger.Fatal(err.Error())
		}
		if exists {
			di.config.Logger.Warn(fmt.Sprintf("Object with key \"%s\" exists, skipping", filepath.Base(di.filePaths[0])))
			di.filePaths = di.filePaths[1:]
		} else {
			break
		}
	}

	// If all is good, open the file, fill in the next struct for processing and return true
	f, err := os.Open(di.filePaths[0])
	di.err = err
	di.next.f = f
	di.next.CannedACL = di.config.Options["acl"].(string)
	di.next.storage = di.config.Options["storage"].(string)
	di.next.name = AppendPrefix(di.config, f.Name())
	di.next.path = di.filePaths[0]
	di.filePaths = di.filePaths[1:]

	return di.Err() == nil
}

func (di *DirectoryIterator) Err() error {
	return di.err
}

/*
UploadObject takes the DirectoryIterator.next struct and returns a s3manager.BatchUploadObject struct to be used by
s3manager.UploadWithIterator. We set After: to a function that closes the file after it's uploaded.
*/
func (di *DirectoryIterator) UploadObject() s3manager.BatchUploadObject {
	f := di.next.f
	di.config.Logger.Info(fmt.Sprintf("Uploading %s...", di.next.name))
	uplCount++
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket:       &di.bucket,
			Key:          &di.next.name,
			ACL:          &di.next.CannedACL,
			StorageClass: &di.next.storage,
			Body:         f,
		},
		After: func() error {
			return f.Close()
		},
	}
}

/*
UploadDirectory takes a config.Configuration struct and a directory path as input. It returns an error.

1. Create a new DirectoryIterator struct (di) with the config and directory path
2. Create a new AWS session with the authentication credentials from the config
3. Create a new s3manager.Uploader with the session
4. Upload the directory with the iterator
  - The iterator runs DirectoryIterator.Next()
  - If there are more files to upload, after next is called, DirectoryIterator.UploadObject() is called
  - UploadObject() returns a s3manager.BatchUploadObject struct, which is used by the uploader
  - The uploader uploads the file to S3, then runs the s3manager.BatchUploadObject.After() function
  - After() closes the file
  - The iterator runs DirectoryIterator.Next() and the process repeats until it Next() returns false.

5. Log the completion of the upload and return nil.
*/
func UploadDirectory(c config.Configuration, dir string) (error, int) {
	// If the directory does not exist, we just stop here and return the error
	dirExists, err := LocalFileExists(dir)
	if !dirExists {
		return errors.New(fmt.Sprintf("directory %s does not exist, skipping", dir)), uplCount
	} else if err != nil {
		return err, uplCount
	}

	c.Logger.Info(fmt.Sprintf("Starting Directory Upload Session for %s", dir))

	di := NewDirectoryIterator(c, dir)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Bucket["region"].(string)),
		Credentials: credentials.NewStaticCredentials(c.Authentication["key"].(string), c.Authentication["secret"].(string), ""),
	})
	if err != nil {
		return err, uplCount
	}

	uploader := s3manager.NewUploader(sess)

	if c.Options["overwrite"] == true {
		c.Logger.Warn("Overwrite is enabled, skipping check for existing objects.")
	}

	err = uploader.UploadWithIterator(aws.BackgroundContext(), di)
	if err != nil {
		return err, uplCount
	}
	c.Logger.Info(fmt.Sprintf("Finished uploading %q to %q", dir, c.Bucket["name"]))
	return nil, uplCount
}
