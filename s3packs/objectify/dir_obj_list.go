package objectify

import (
	"fmt"

	"github.com/orme292/s3packer/conf"
)

type DirObjList []*DirObj

func NewDirObjList(ac *conf.AppConfig, path string) (dol DirObjList, err error) {
	fmt.Printf("Processing Directory: %q\n", path)
	dirs, err := getSubDirs(path)
	if err != nil {
		return nil, err
	}
	dol = make(DirObjList, len(dirs))
	for i, p := range dirs {
		fmt.Printf("Processing Sub-Directory: %q\n", p)
		do, err := NewDirObj(ac, p, path)
		if err != nil {
			return nil, err
		}
		dol[i] = do
	}
	return dol, err
}
