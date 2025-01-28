package conf

import (
	"fmt"
	"strings"
)

// Naming type is a string enum of the supported object naming methods.
type Naming string

const (
	NamingRelative Naming = "relative"
	NamingAbsolute Naming = "absolute"
	NamingNone     Naming = "none"
)

// String returns the string representation of the Naming object.
// It converts the Naming object to a string by using the underlying string value.
func (n Naming) String() string {
	return string(n)
}

// Objects contain the object naming configuration.
type Objects struct {
	NamingType Naming
	NamePrefix string
	PathPrefix string

	// OmitRootDir is used to remove the root directory name from the object's final FormattedKey.
	OmitRootDir bool
}

func (o *Objects) build(inc *ProfileIncoming) error {
	switch tidyLowerString(inc.Objects.NamingType) {

	case NamingAbsolute.String(), "abs":
		o.NamingType = NamingAbsolute

	case NamingRelative.String(), "rel":
		o.NamingType = NamingRelative

	default:
		o.NamingType = NamingNone

	}

	o.NamePrefix = strings.TrimPrefix(inc.Objects.NamePrefix, "/")
	o.PathPrefix = strings.TrimPrefix(inc.Objects.PathPrefix, "/")
	o.PathPrefix = strings.TrimSuffix(inc.Objects.PathPrefix, "/")
	o.OmitRootDir = inc.Objects.OmitRootDir

	return o.validate()

}

func (o *Objects) validate() error {

	if o.NamingType == NamingNone {
		return fmt.Errorf("bad objects config: %v", InvalidNamingType)
	}

	return nil

}
