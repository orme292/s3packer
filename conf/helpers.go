package conf

import (
	"regexp"
	"strings"
)

func formatPath(p string) string {
	p = strings.TrimPrefix(p, "/")
	// Trimming ending slash if exists
	p = strings.TrimSuffix(p, "/")
	return p
}

func alphaNumericString(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(s, "")
}

func tidyString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	return s
}

func whichProvider(s string) ProviderName {
	s = tidyString(s)
	switch s {
	case "aws", "amazon", "s3", "amazon s3":
		return ProviderNameAWS
	case "oci", "oracle", "oraclecloud", "oracle cloud":
		return ProviderNameOCI
	default:
		return ProviderNameNone
	}
}
