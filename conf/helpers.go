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
