package conf

import (
	"fmt"
	"log"
)

// Overwrite type is a string enum of the supported overwrite methods. OverwriteChecksum is not implemented.
// Overwrite.String() will return the string representation of the enum for convenience, either in output or logging.
type Overwrite string

const (
	OverwriteChecksum Overwrite = "checksum"
	OverwriteNever    Overwrite = "never"
	OverwriteAlways   Overwrite = "always"
)

func (o Overwrite) String() string {
	return string(o)
}

// Opts contains application level configuration options.
type Opts struct {
	MaxUploads     int
	FollowSymlinks bool
	WalkDirs       bool
	Overwrite      Overwrite
}

func (o *Opts) build(inc *ProfileIncoming) error {

	o.MaxUploads = inc.Options.MaxUploads

	switch tidyLowerString(inc.Options.OverwriteObjects) {

	case OverwriteAlways.String(), "yes", "true":
		o.Overwrite = OverwriteAlways

	case OverwriteNever.String(), "no", "false":
		o.Overwrite = OverwriteNever

	// checksum not supported yet
	case Empty:
		return fmt.Errorf("bad options config: %s", InvalidOverwriteMethod)

	default:
		o.Overwrite = OverwriteNever

	}

	return o.validate()

}

func (o *Opts) validate() error {

	if o.MaxUploads <= 0 {
		log.Printf("here")
		return fmt.Errorf("MaxUploads must be at least 1")
	}
	if o.Overwrite != OverwriteChecksum && o.Overwrite != OverwriteNever && o.Overwrite != OverwriteAlways {
		return fmt.Errorf("OverwriteObjects value should be \"never\" or \"always\": %q", o.Overwrite)
	}
	return nil

}
