package s3pack

import (
	"os"

	"github.com/orme292/s3packer/config"
)

/*
DirectoryIterator is used with the s3manager.BatchUploadIterator to process a directory of files for upload to s3.
*/
type DirectoryIterator struct {
	filePaths []string
	bucket    string
	next      struct {
		path      string
		f         *os.File
		name      string
		storage   string
		CannedACL string
	}
	err    error
	config config.Configuration
}

/*
FileIterator is used with the s3manager.BatchUploadIterator to process a list of files for upload to s3.
*/
type FileIterator struct {
	filePaths []string
	bucket    string
	next      struct {
		path      string
		f         *os.File
		name      string
		storage   string
		CannedACL string
	}
	err    error
	config config.Configuration
}
