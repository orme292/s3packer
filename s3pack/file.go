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

/*
NewFileIterator creates a new FileIterator struct, which implements the s3manager.BatchUploadIterator
interface for use with s3manager.UploadWithIterator. It takes a config.Configuration struct as input. It returns a
pointer to a FileIterator struct.

It creates a list of files to upload from the configuration (c.Files) using the provided file's absolute path.
The resulting list is assigned to FileIterator.filePaths
*/
func NewFileIterator(c config.Configuration) s3manager.BatchUploadIterator {
	var paths []string
	for _, file := range c.Files {
		path, err := filepath.Abs(file)
		if err != nil {
			c.Logger.Fatal("Unable to get absolute path of file " + file + ":" + err.Error())
		}
		paths = append(paths, path)
	}
	return &FileIterator{
		filePaths: paths,
		bucket:    c.Bucket["name"].(string),
		config:    c,
	}
}

/*
Next returns true if there are more files to upload, or false if there are no more files to upload. It also
sets the FileIterator.next struct for processing.

 1. If there are no more files to upload, return false
 2. If there are more files to upload, open the file, fill in the FileIterator.next struct for processing and
    return true.

FileIterator.next will be used by UploadObject() to upload the file to S3.
*/
func (fi *FileIterator) Next() bool {
	if len(fi.filePaths) == 0 {
		fi.next.f = nil
		return false
	}

	f, err := os.Open(fi.filePaths[0])
	fi.err = err
	fi.next.f = f
	fi.next.CannedACL = fi.config.Options["acl"].(string)
	fi.next.storage = fi.config.Options["storage"].(string)
	fi.next.name = AppendPrefix(fi.config, filepath.Base(fi.filePaths[0]))
	fi.next.path = fi.filePaths[0]
	fi.filePaths = fi.filePaths[1:]

	return fi.Err() == nil
}

func (fi *FileIterator) Err() error {
	return fi.err
}

/*
UploadObject returns a s3manager.BatchUploadObject struct for use with s3manager.UploadWithIterator. It takes no
input and returns a s3manager.BatchUploadObject struct.

 1. Using the FileIterator.next struct, it uses a s3manager.UploadInput struct as the object data in
    s3manager.BatchUploadObject for use with s3manager.UploadWithIterator.
    It also sets the StorageClass to the value of FileIterator.next.storage
 2. It returns a s3manager.BatchUploadObject struct with the UploadInput struct and a function to close the file.
*/
func (fi *FileIterator) UploadObject() s3manager.BatchUploadObject {
	f := fi.next.f
	fi.config.Logger.Info(fmt.Sprintf("Uploading %s...", fi.next.name))
	uplCount++
	return s3manager.BatchUploadObject{
		Object: &s3manager.UploadInput{
			Bucket:       &fi.bucket,
			Key:          &fi.next.name,
			ACL:          &fi.next.CannedACL,
			StorageClass: &fi.next.storage,
			Body:         f,
		},
		After: func() error {
			return f.Close()
		},
	}
}

/*
UploadObjects takes a config.Configuration struct as input. It returns an error.

 1. It checks to make sure all the files listed in the profile exist locally, using ListedLocalFileExists().
 2. It checks whether there are already objects with the same names as the files, using ListedObjectsNotExisting().
 3. It creates a new FileIterator struct with the config.
 4. It creates a new AWS session with the authentication credentials from the config.
 5. It creates a new s3manager.Uploader with the session.
 6. It uploads the files with the iterator.
    - The iterator runs FileIterator.Next()
    - If there are more files to upload, after next is called, FileIterator.UploadObject() is called
    - UploadObject() returns a s3manager.BatchUploadObject struct, which is used by the uploader
    - The uploader uploads the file to S3, then runs the s3manager.BatchUploadObject.After() function
    - The After() function closes the file
    - Next() is called again, and the process repeats until there are no more files to upload.
*/
func UploadObjects(c config.Configuration) (error, int) {
	c.Logger.Info("Starting Individual File Upload Session...")

	// Check to make sure all the files listed in the profile exist locally
	fileList, err := ListedLocalFileExists(c, c.Files)
	if err != nil {
		return err, uplCount
	}
	c.Files = fileList

	// Check whether there are already objects keys named the same as the file names.
	objectList, err := ListedObjectsNotExisting(c, c.Files)
	if err != nil {
		return err, uplCount
	}
	c.Files = objectList

	fi := NewFileIterator(c)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(c.Bucket["region"].(string)),
		Credentials: credentials.NewStaticCredentials(c.Authentication["key"].(string), c.Authentication["secret"].(string), ""),
	})
	if err != nil {
		return errors.New("unable to create AWS session"), uplCount
	}
	svc := s3manager.NewUploader(sess)

	if err := svc.UploadWithIterator(aws.BackgroundContext(), fi); err != nil {
		return err, uplCount
	}
	return nil, uplCount
}
