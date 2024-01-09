package objectify

import (
	"fmt"
	"reflect"

	"github.com/orme292/s3packer/conf"
)

type DirObj struct {
	Path         string
	RelativeRoot string
	Fol          FileObjList
	Ac           *conf.AppConfig
}

func NewDirObj(ac *conf.AppConfig, p, rr string) (do *DirObj, err error) {
	do = &DirObj{
		Path:         p,
		RelativeRoot: rr,
		Ac:           ac,
	}
	files, err := getFiles(ac, p)
	if err != nil {
		return nil, err
	}

	ac.Log.Info("Processing Directory: (%d objects) %q", len(files), p)
	do.Fol, err = NewFileObjList(ac, files, rr)
	if err != nil {
		return nil, err
	}
	return
}

/* DEBUG */

func (do *DirObj) Values() {
	val := reflect.ValueOf(do)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typeOfT := val.Type()

	fmt.Println("--------------------")
	fmt.Printf("Type: %s\n", typeOfT)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.CanInterface() {
			fmt.Printf("Field: %s\tValue: %v\n", typeOfT.Field(i).Name, field.Interface())
		}
	}
	fmt.Printf("Has: %d objects\n", len(do.Fol))
}
