package gcloud

import (
	"fmt"
	"os"

	"github.com/orme292/s3packer/internal/provider"
)

type CloudObject struct {
	job *provider.Job

	f *os.File

	key    string
	tagMap map[string]string
}

func NewCloudObject(job *provider.Job) provider.Object {
	return &CloudObject{
		job:    job,
		tagMap: make(map[string]string),
	}
}

func (o *CloudObject) Destroy() error { return o.Post() }

func (o *CloudObject) Generate() error {

	o.key = o.job.Key

	o.setTags(o.job.App.Tags)

	return nil

}

func (o *CloudObject) Post() error { return o.f.Close() }

func (o *CloudObject) Pre() error {

	o.job.Metadata.Update()

	if !o.job.Metadata.IsExists || !o.job.Metadata.IsReadable {
		return fmt.Errorf("file no longer accessible")
	}

	f, err := os.Open(o.job.Metadata.FullPath())
	if err != nil {
		return err
	}

	o.f = f

	return nil

}

func (o *CloudObject) setTags(input map[string]string) {

	if len(input) != 0 {
		o.tagMap = input
	}

	for key := range o.job.AppTags {
		o.tagMap[key] = o.job.AppTags[key]
	}

}
