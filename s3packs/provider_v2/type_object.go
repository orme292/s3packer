package provider_v2

import (
	"os"
)

type Object struct {
	F        os.File
	Key      string
	Checksum struct {
		SHA256 string
	}
}
