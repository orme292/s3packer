package provider_v2

import (
	"os"
	"sync"

	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
)

type QueueJob struct {
	F       *os.File
	Obj     Object
	Details *objectify.FileObj

	Key        string
	SearchRoot string

	App *conf.AppConfig

	status int
	err    error

	mu *sync.RWMutex
}

func newJob(object *objectify.FileObj, app *conf.AppConfig, searchRoot string, status int, err error) *QueueJob {

	j := &QueueJob{
		Details:    object,
		SearchRoot: searchRoot,
		App:        app,
		status:     status,
		err:        err,
		mu:         &sync.RWMutex{},
	}

	j.setKey()

	return j

}

func (j *QueueJob) CloseFile() error {

	j.mu.Lock()
	defer j.mu.Unlock()

	err := j.F.Close()
	if err != nil {
		return err
	}

	return nil

}

func (j *QueueJob) OpenFile() error {

	j.mu.Lock()
	defer j.mu.Unlock()

	f, err := os.Open(j.Details.FullPath())
	if err != nil {
		return err
	}

	j.F = f

	return nil
}

func (j *QueueJob) setKey() {

	key := ObjectKey{
		base: j.Details.Filename,
		dir:  j.Details.Root,

		searchRoot: j.SearchRoot,

		namePrefix: j.App.Objects.NamePrefix,
		pathPrefix: j.App.Objects.PathPrefix,
	}

	j.Key = key.String(j.App.Objects.NamingType, j.App.Objects.OmitRootDir)

}

func (j *QueueJob) updateDetails() {

	j.mu.Lock()
	defer j.mu.Unlock()

	j.Details.Update()

}

func (j *QueueJob) updateStatus(status int, err error) {

	j.mu.Lock()
	defer j.mu.Unlock()

	j.status = status
	j.err = err

}
