package s3pack

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/orme292/s3packer/conf"
)

/*
AppendObjectPrefix takes the application configuration (c) and a key string (key) and returns a new string with the prefix
prepended to the key.
*/
func AppendObjectPrefix(a *conf.AppConfig, key string) string {
	if a.Objects.NamePrefix == EmptyString {
		return key
	}
	return fmt.Sprintf("%s%s", a.Objects.NamePrefix, filepath.Base(key))
}

/*
AppendPathPrefix takes the application configuration (c) and a key string (key) and returns a new string with the
pathPrefix prepended to the key.

This is the last step in the process of building the object key (PrefixedKey)
*/
func AppendPathPrefix(a *conf.AppConfig, key string) string {
	if a.Objects.RootPrefix == EmptyString {
		return key
	}
	return path.Clean(fmt.Sprintf("/%s/%s", a.Objects.RootPrefix, key))
}

/*
BucketExists takes the application configuration (c) and returns a bool and an error. It confirms that the bucket
provided in the configuration exists and is accessible.
*/
func BucketExists(a *conf.AppConfig) (exists bool, err error) {
	client, _ := BuildClient(a)
	_, err = client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(a.Bucket.Name),
	})
	if err != nil {
		if errors.As(err, &s3Error) {
			switch {
			case errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchBucket):
				return false, errors.New(fmt.Sprintf("aws says bucket %q does not exist",
					a.Bucket.Name))
			default:
				return false, errors.New(fmt.Sprintf("aws error when checking if %q exists: %q",
					a.Bucket.Name, err))
			}
		}
	}
	return true, nil
}

/*
CalcChecksumSHA256 takes a path string as input and returns a checksum string and an error, if there is one.
You should check if the path exists and is readable before using this function.
*/
func CalcChecksumSHA256(p string) (checksum string, err error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return
	}

	f, err := os.Open(filepath.Clean(absPath))
	if err != nil {
		return
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			return
		}
	}(f)

	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

/*
FileSizeString takes an int64 (size) and returns a human-readable file size string.
*/
func FileSizeString(size int64) string {
	switch {
	case size < 1024:
		return fmt.Sprintf("%d bytes", size)
	case size < 1024*1024:
		return fmt.Sprintf("%d KB", size/1024)
	case size < 1024*1024*1024:
		return fmt.Sprintf("%.2f MB", float32(size)/(1024*1024))
	case size < 1024*1024*1024*1024:
		return fmt.Sprintf("%.2f GB", float32(size)/(1024*1024*1024))
	default:
		return fmt.Sprintf("%.2f TB", float32(size)/(1024*1024*1024*1024))
	}
}

/*
GetFiles returns a list of files in a directory. It takes a path string as input and returns a slice of strings and an
error, if there is one.

It does NOT walk subdirectories.
*/
func GetFiles(p string) (files []string, err error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return nil, errors.New("Error getting absolute path: " + err.Error())
	}

	objects, err := os.ReadDir(absPath)
	if err != nil {
		return nil, errors.New("Error reading directory: " + err.Error())
	}
	for _, file := range objects {
		if !file.IsDir() {
			files = append(files, fmt.Sprintf("%s/%s", absPath, file.Name()))
		}
	}
	return
}

/*
GetFileSize returns the size of a file in bytes. It takes a path string as input and returns an int64 and an error
if there is one.
*/
func GetFileSize(p string) (size int64, err error) {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return 0, err
	}
	fInfo, err := os.Stat(absPath)
	if err != nil {
		return 0, err
	}
	size = fInfo.Size()
	return size, nil
}

/*
GetSubDirs returns a list of subdirectories at all depths in a given directory. It takes a path string as input and
returns a slice of strings and an error, if there is one.

TODO: Add depth option
*/
func GetSubDirs(p string) (subDirs []string, err error) {
	absPath, err := filepath.Abs(filepath.Clean(p))
	if err != nil {
		return nil, errors.New("Error getting absolute path: " + err.Error())
	}

	err = filepath.Walk(absPath, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			subDirs = append(subDirs, path)
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

/*
LocalFileExists takes a string (file) and returns a bool and an error.
The function checks to see if the file exists. If it does, then it returns true and a nil error.
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
ObjectExists takes the application configuration (c) and a prefixed string (key).
The function checks to see if an identically named object exists in the bucket. If it does, then it returns true and
a nil error. If it doesn't, then false and a nil error are returned.
*/
func ObjectExists(a *conf.AppConfig, objectKey string) (exists bool, err error) {
	if objectKey == "" {
		return false, errors.New(ErrIgnoreObjectKeyEmpty)
	}

	client, _ := BuildClient(a)
	exists, err = ObjectExistsWithClient(a, objectKey, client)
	return
}

func ObjectExistsWithClient(a *conf.AppConfig, objectKey string, client *s3.Client) (exists bool, err error) {
	if objectKey == "" {
		return false, errors.New(ErrIgnoreObjectKeyEmpty)
	}
	a.Log.Debug("Checking if object %q already exists", objectKey)
	_, err = client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(a.Bucket.Name),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		if errors.As(err, &s3Error) {
			switch {
			case errors.As(s3Error, &s3NotFound) || errors.As(s3Error, &s3NoSuchKey):
				return false, nil
			case errors.As(s3Error, &s3NoSuchBucket):
				return false, errors.New(fmt.Sprintf("aws says bucket %q does not exist", a.Bucket.Name))
			default:
				return false, errors.New(fmt.Sprintf("aws error: %q", err))
			}
		}
	}
	return true, nil
}

/* DEBUG */

/*
PrintMemUsage is a debug function
*/
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

/*
ExecutionTime is a debug function.
*/
func ExecutionTime(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", name, time.Since(start))
	}
}
