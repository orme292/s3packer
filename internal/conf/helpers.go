package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func canCreate(path string) (bool, error) {
	filename := expandHome(path)

	filename, err := filepath.Abs(filename)
	if err != nil {
		return false, err
	}

	// Resolve G304: Potential file inclusion via variable
	if strings.Contains(filename, "..") {
		return false, fmt.Errorf("invalid filename: %s", filename)
	}

	_, err = os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	return false, fmt.Errorf("file %s already exists", filename)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
		}
		return strings.Replace(path, "~", home, 1)
	}
	return path
}

func sanitizeString(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9-+_]+")
	return reg.ReplaceAllString(s, "")
}

// tidyLowerString takes a string and performs two operations on it: trimming any leading/trailing whitespace and converting it to lowercase.
// It then returns the resulting modified string.
func tidyLowerString(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

func tidyUpperString(s string) string {
	return strings.TrimSpace(strings.ToUpper(s))
}
