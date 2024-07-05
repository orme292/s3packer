package provider_v2

import (
	"fmt"
	"path"
	"strings"

	"github.com/orme292/s3packer/conf"
)

type ObjectKey struct {
	base string
	dir  string

	searchRoot string

	namePrefix string
	pathPrefix string

	key, path string
}

func (k *ObjectKey) String(method conf.Naming, omit bool) string {

	// TODO: Remove
	// fmt.Println(
	// 	fmt.Sprintf("\nbase: %s\n\tdir: %s\n\tsearchRoot: %s\n\tnamePrefix: %s\n\tpathPrefix: %s\n\t", k.base, k.dir, k.searchRoot, k.namePrefix, k.pathPrefix))
	//
	// // base: file.txt
	// // dir: /users/aorme/documents
	// // searchRoot: /users
	// // namePrefix: backup-
	// // pathPrefix: /new/backups
	// // RELATIVE RESULT: new/backups/aorme/documents/backup-file.txt
	// // ABSOLUTE RESULT: /new/backups/users/aorme/documents/backup-file.txt

	strip := func(s string) string {
		invalid := []string{":", "*", "?", "\"", "<", ">", "|"}
		s = strings.TrimSpace(path.Clean(s))
		for _, char := range invalid {
			s = strings.ReplaceAll(s, char, "")
		}
		s = strings.TrimPrefix(s, "/")
		s = strings.TrimSuffix(s, "/")
		return s
	}

	k.key = fmt.Sprintf("%s%s", strip(k.namePrefix), k.base)
	k.key = strip(k.key)

	if k.searchRoot != EmptyPath {

		if omit {
			k.path = strings.TrimPrefix(strings.TrimPrefix(k.dir, k.searchRoot), "/")
		} else {
			k.path = strings.Join(strings.Split(k.dir, "/")[strings.Count(k.searchRoot, "/"):], "/")
		}

	}

	k.pathPrefix = strip(k.pathPrefix)

	switch method {
	case conf.NamingRelative:
		k.path = strip(fmt.Sprintf("%s/%s", k.pathPrefix, k.path))
	default:
		k.path = fmt.Sprintf("%s", k.pathPrefix)
		k.path = fmt.Sprintf("%s/%s", k.pathPrefix, strings.TrimPrefix(k.dir, "/"))
	}

	return fmt.Sprintf("%s/%s", k.path, k.key)

}
