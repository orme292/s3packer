package conf

import (
	"fmt"
)

type Tags map[string]string

func (t *Tags) Get() map[string]string {
	return *t
}

func (t *Tags) build(tags map[string]string) error {
	for k, v := range tags {
		(*t)[sanitizeString(k)] = sanitizeString(v)
	}

	return t.validate()
}

func (t *Tags) validate() error {
	for key := range *t {
		if tidyLowerString(key) == "s3p-checksum-sha256" || tidyLowerString(key) == "s3p-origin-path" {
			return fmt.Errorf("reserved tag '%s' cannot be used", key)
		}
	}

	return nil
}
