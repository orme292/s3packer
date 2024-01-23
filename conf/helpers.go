package conf

import (
	"strings"
)

func formatPath(p string) string {
	p = strings.TrimPrefix(p, "/")
	// Trimming ending slash if exists
	p = strings.TrimSuffix(p, "/")
	return p
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
