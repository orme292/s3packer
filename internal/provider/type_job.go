package provider

import (
	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/internal/conf"
)

type Job struct {
	App *conf.AppConfig

	Object     Object
	Metadata   *objectify.FileObj
	Key        string
	SearchRoot string

	AppTags map[string]string

	status int

	err error
}

func newJob(app *conf.AppConfig, metadata *objectify.FileObj, searchRoot string) *Job {

	j := &Job{
		App:        app,
		Metadata:   metadata,
		SearchRoot: searchRoot,
		AppTags:    make(map[string]string),
		status:     JobStatusQueued,
		err:        nil,
	}

	j.setKey()
	j.setAppTags()

	return j

}

func (j *Job) setAppTags() {

	if j.App.TagOpts.ChecksumSHA256 && j.Metadata != nil {
		if j.Metadata.ChecksumSHA256 != "" {
			j.AppTags["ChecksumSHA256"] = j.Metadata.ChecksumSHA256
		}
	}

	if j.App.TagOpts.OriginPath && j.Metadata != nil {
		j.AppTags["Origin"] = j.Metadata.FullPath()
	}

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
