package objectify

import (
	"github.com/orme292/s3packer/conf"
)

type DirObjList []*DirObj

func NewDirObjList(ac *conf.AppConfig, path string) (dol DirObjList, err error) {
	dirs, err := getSubDirs(path)
	if err != nil {
		return nil, err
	}
	dol = make(DirObjList, len(dirs))
	for i, p := range dirs {
		do, err := NewDirObj(ac, p, path)
		if err != nil {
			return nil, err
		}
		dol[i] = do
	}
	return dol, err
}
