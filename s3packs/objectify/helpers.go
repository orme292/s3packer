package objectify

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/orme292/s3packer/conf"
)

// getFileSize returns the size of a file in bytes.
// It takes an `ap` string as an argument which should be the absolute path of the file.
// It uses the `os.Stat` method to retrieve the file's information and extracts the file size from there.
// If there is an error (for instance, the file doesn't exist or there are insufficient permissions),
// the function will return 0 for the size and carry the error. Otherwise, it will return the file size and `nil` for the error.
func getFileSize(ap string) (size int64, err error) {
	fi, err := os.Stat(ap)
	if err != nil {
		return 0, err
	}
	size = fi.Size()
	return size, nil
}

func getFiles(ac *conf.AppConfig, p string) (files []string, err error) {
	ap, err := filepath.Abs(filepath.Clean(p))
	if err != nil {
		return nil, errors.New("Error getting absolute path: " + err.Error())
	}

	ap, err = filepath.EvalSymlinks(ap)
	if err != nil {
		return nil, errors.New("Error evaluating link: " + err.Error())
	}

	objs, err := os.ReadDir(ap)
	if err != nil {
		return nil, errors.New("Error reading directory: " + err.Error())
	}
	for _, file := range objs {
		info, _ := file.Info()
		mode := info.Mode()
		if mode.IsRegular() && !info.IsDir() {
			files = append(files, filepath.Join(ap, file.Name()))
		} else if !mode.IsRegular() && !info.IsDir() {
			ac.Log.Warn("Skipping non-regular file: %q", filepath.Join(ap, file.Name()))
		}
	}
	return
}

func isRegular(ap string) (bool, error) {
	fi, err := os.Stat(ap)
	if err != nil {
		return false, err
	}
	return fi.Mode().IsRegular(), nil
}

func getSubDirs(p string) (dirs []string, err error) {
	ap, err := filepath.Abs(filepath.Clean(p))
	if err != nil {
		return nil, errors.New("Error getting absolute path: " + err.Error())
	}

	ap, err = filepath.EvalSymlinks(ap)
	if err != nil {
		return nil, errors.New("Error evaluating link: " + err.Error())
	}

	err = filepath.WalkDir(ap, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if d.Type()&fs.ModeSymlink != 0 {
				path, err = os.Readlink(path)
				if err != nil {
					return err
				}

				rInfo, err := os.Stat(path)
				if err != nil {
					return err
				}

				if rInfo.IsDir() {
					dirs = append(dirs, path)
				}
			} else {
				dirs = append(dirs, path)
			}
		}
		return nil
	})

	return
}

// fileExists verifies the existence of a file at the provided absolute path (`ap`).
// It uses the `os.Stat` function to retrieve the file's information.
// If the file does not exist, it returns `false` and `nil` error.
// Any other error during the `os.Stat` call is returned as `false` and the respective error.
// If the file exists and there are no errors, it returns `true` and `nil` error.
func fileExists(ap string) (exists bool, err error) {
	if _, err := os.Stat(ap); errors.Is(err, fs.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// GetChecksumSHA256 takes an absolute path (`ap`) to a file as a string and returns
// the SHA-256 checksum of the file content as a base64 encoded string.
// It opens the file, calculates the checksum by reading the file content,
// then base64 encodes the checksum and returns it along with any error encountered
// during the process. If an error occurs at any step (like during file opening, reading,
// or closing), it returns the error and an empty checksum string.
func GetChecksumSHA256(ap string) (cs string, err error) {
	f, err := os.Open(ap)
	if err != nil {
		fmt.Println("os.Open error")
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			fmt.Println("Checksum Error, closing file:", closeErr)
		}
	}()

	hash := sha256.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

func GetChecksumSHA256Reader(r io.Reader) (cs string, err error) {
	hash := sha256.New()
	_, err = io.Copy(hash, r)
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(hash.Sum(nil)), nil
}

// formatFullKey`function comprises two string parameters: `base`, which is the base file name,
// and `od`, which is the original directory path, rr represents the relative root directory path.
// The function also takes the `AppConfig` object as a parameter.
// First, the function concatenates the `NamePrefix` defined in `AppConfig` and the base file name,
// and then sanitizes it with `stripSafePath`. If a non-empty relative root is provided, the function
// formats the original directory by invoking `formatPseudoPath`. Then, it concatenates the
// `RootPrefix` defined in `AppConfig` and the formatted directory (`fPseudo`), and sanitizes it.
// The function finally returns the sanitized `fName` and `fPseudo`, which can be used directly in
// object storage operations.
func formatFullKey(ac *conf.AppConfig, base string, od string, rr string) (fName string, fPseudo string) {
	fName = s("%s%s", ac.Objects.NamePrefix, base)
	fName = stripSafePath(fName)
	if rr != EmptyString {
		fPseudo = formatPseudoPath(ac, od, rr)
	}
	switch ac.Objects.Naming {
	case conf.NamingRelative:
		fPseudo = s("%s/%s", ac.Objects.RootPrefix, fPseudo)
		fPseudo = stripSafePath(fPseudo)
	default:
		fPseudo = s("%s/%s", ac.Objects.RootPrefix, strings.TrimPrefix(od, "/"))
	}
	return fName, fPseudo
}

// formatPseudoPath takes an AppConfig object, origin dir path (od), and relative root path (rr)
// and returns a formatted string representing a path.
// It checks the `OmitRootDir` field of the `Objects` struct in the `AppConfig`.
// If `OmitRootDir` is true, the function trims the relative path prefix and any trailing
// slash from the original path.
// If `OmitRootDir` is false, the function trims the leading directory paths from the original
// path up to the depth of the relative path and returns the resulting path. The depth is
// determined by the count of slashes in the relative path.
func formatPseudoPath(ac *conf.AppConfig, od string, rr string) string {
	if ac.Objects.OmitRootDir {
		return strings.TrimPrefix(strings.TrimPrefix(od, rr), "/")
	} else {
		return strings.Join(strings.Split(od, "/")[strings.Count(rr, "/"):], "/")
	}
}

// stripSafePath takes a file path as parameter and returns a sanitized version of it.
// It first trims whitespace from the start/end of the string and cleans the path
// (removes any unnecessary slashes) using the path.Clean method.
// It removes invalid characters, such as (`:`, `*`, `?`, `"`, `<`, `>`, `|`),
// After that, any leading and trailing slashes are removed.
func stripSafePath(p string) string {
	invalidChars := []string{":", "*", "?", "\"", "<", ">", "|"}
	p = strings.TrimSpace(path.Clean(p))
	for _, char := range invalidChars {
		p = strings.ReplaceAll(p, char, "")
	}
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimSuffix(p, "/")
	return p
}

func s(format string, s ...any) string {
	return fmt.Sprintf(format, s...)
}

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
