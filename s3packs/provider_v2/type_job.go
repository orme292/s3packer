package provider_v2

import (
	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
)

type Job struct {
	App *conf.AppConfig

	Object     Object
	Metadata   *objectify.FileObj
	Key        string
	SearchRoot string

	status int

	err error
}

func newJob(app *conf.AppConfig, metadata *objectify.FileObj, searchRoot string) *Job {

	j := &Job{
		App:        app,
		Metadata:   metadata,
		SearchRoot: searchRoot,
		err:        nil,
	}

	j.setKey()

	return j

}

func (j *Job) setKey() {

	key := ObjectKey{
		base: j.Metadata.Filename,
		dir:  j.Metadata.Root,

		searchRoot: j.SearchRoot,

		namePrefix: j.App.Objects.NamePrefix,
		pathPrefix: j.App.Objects.PathPrefix,
	}

	j.Key = key.String(j.App.Objects.NamingType, j.App.Objects.OmitRootDir)

}

func (j *Job) setStatus(status int, err error) {
	j.status = status
	j.err = err
}
