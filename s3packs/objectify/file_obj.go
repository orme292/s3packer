package objectify

import (
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/orme292/s3packer/conf"
)

type FileObj struct {
	OriginPath string
	OriginDir  string
	Base       string
	RelRoot    string
	AbsPath    string

	FileSize       int64
	ChecksumSHA256 string

	FName    string
	FPseudoP string
	TagsMap  map[string]string

	Ignore          bool
	IgnoreString    string
	IsDirectoryPart bool
	IsUploaded      bool

	Group int

	ac *conf.AppConfig
}

func NewFileObj(ac *conf.AppConfig, p string, rel string, grp int) (fo *FileObj, err error) {
	ap, err := filepath.Abs(filepath.Clean(p))
	if err != nil {
		return nil, err
	}

	fo = &FileObj{
		OriginPath:      p,
		OriginDir:       filepath.Dir(ap),
		Base:            filepath.Base(p),
		RelRoot:         rel,
		AbsPath:         ap,
		FName:           EmptyString,
		TagsMap:         make(map[string]string),
		IsDirectoryPart: rel == EmptyString,
		Group:           grp,
		ac:              ac,
	}
	exists, err := fileExists(fo.AbsPath)
	if err != nil {
		fo.setIgnore(s("file access issue: %q", err))
	}
	if !exists {
		fo.setIgnore("file does not exist")
	} else {
		fo.FileSize, err = getFileSize(fo.AbsPath)
		if err != nil {
			fo.setIgnore(s("could not get file size: %q", err))
		}
		fo.ChecksumSHA256, err = GetChecksumSHA256(fo.AbsPath)
		if err != nil {
			fo.setIgnore(s("could not get checksum: %q", err))
		}
		for k, v := range ac.Tags {
			fo.addTag(k, v)
		}
		fo.addTag("ChecksumSHA256", fo.ChecksumSHA256)
		fo.FName, fo.FPseudoP = formatFullKey(ac, fo.Base, fo.OriginDir, fo.RelRoot)
	}
	ac.Log.Debug("Processed file: %q", fo.FKey())
	return fo, err
}

// setIgnore sets the Ignore field of the `FileObj` object to true.
// It also sets the IgnoreString field to the string argument s.
// s should be used to specify the reason why this file is ignored.
func (fo *FileObj) setIgnore(s string) {
	fo.Ignore = true
	fo.IgnoreString = s
}

func (fo *FileObj) setRelRoot(p string) {
	fo.RelRoot = p
}

// addTag adds a net tag to the FileObj.TagsMap map. K is the tag key
// and v is the tag value.  Some providers, like AWS, want tags to be
// URL encoded and combined to a single string -- that will need to
// be done separately.
// TODO: duplicate tag key checking
func (fo *FileObj) addTag(k, v string) {
	fo.TagsMap[k] = v
}

// FKey method concatenates the `FPseudoP` and `FName` fields of the
// `FileObj` instance with a '/' in between. `FPseudoP` represents
// the pseudo path of the file, while `FName` is the sanitized file name.
func (fo *FileObj) FKey() string {
	return s("%s/%s", fo.FPseudoP, fo.FName)
}

/* DEBUG */

func (fo *FileObj) Values() {
	val := reflect.ValueOf(fo)
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
	fmt.Printf("FKey: %s\n", fo.FKey())
}
