package provider_v2

import (
	"os"

	"github.com/orme292/objectify"
	"github.com/orme292/s3packer/conf"
)

type queueJob struct {
	f       *os.File
	app     *conf.AppConfig
	details *objectify.FileObj
	root    string
	status  int
	err     error
}

func newJob(object *objectify.FileObj, root string, status int, err error) *queueJob {
	return &queueJob{
		details: object,
		root:    root,
		status:  status,
		err:     err,
	}
}

func (j *queueJob) closeFile() error {

	err := j.f.Close()
	if err != nil {
		return err
	}

	return nil

}

func (j *queueJob) openFile() error {

	f, err := os.Open(j.details.FullPath())
	if err != nil {
		return err
	}

	j.f = f

	return nil
}

func (j *queueJob) updateDetails() {
	j.details.Update()
}

func (j *queueJob) updateStatus(status int, err error) {
	j.status = status
	j.err = err
}
