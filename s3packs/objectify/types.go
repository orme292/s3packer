package objectify

import (
	"github.com/orme292/s3packer/conf"
)

/*
RootList is a list of FileObjLists.
When you call NewRootList and specify a list of paths, it determines
all the subdirectories of those paths and creates a FileObjList for
each subdirectory.
*/
type RootList []FileObjList

type FileObjList []*FileObj

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
	IsFailed        bool
	IsFailedString  string
	IsUploaded      bool

	Group int

	ac *conf.AppConfig
}

type Stats struct {
	Bytes    int64
	Uploaded int
	Ignored  int
	Failed   int
	Objects  int
	Discrep  int
}

func (s *Stats) Add(s2 Stats) {
	s.Bytes += s2.Bytes
	s.Uploaded += s2.Uploaded
	s.Ignored += s2.Ignored
	s.Failed += s2.Failed
	s.Objects += s2.Objects
}

const (
	EmptyString = ""
)
