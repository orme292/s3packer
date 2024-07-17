package oci

import (
	"fmt"
	"regexp"

	"github.com/orme292/s3packer/s3packs/provider"
)

type OracleObject struct {
	job *provider.Job

	key    string
	bucket string

	tags map[string]string
}

func NewOracleObject(job *provider.Job) provider.Object {
	return &OracleObject{
		job: job,
	}
}

func (o *OracleObject) Destroy() error {
	return o.Post()
}

func (o *OracleObject) Generate() error {

	o.tags = make(map[string]string)
	o.key = o.job.Key
	o.bucket = o.job.App.Bucket.Name

	o.setTags()
	return nil

}

func (o *OracleObject) Post() error {
	return nil
}

func (o *OracleObject) Pre() error {

	o.job.Metadata.Update()

	if !o.job.Metadata.IsExists || !o.job.Metadata.IsReadable {
		return fmt.Errorf("file no longer accessible")
	}

	return nil

}

func (o *OracleObject) setTags() {

	if len(o.job.App.Tags) == 0 {
		o.tags = make(map[string]string)
		return
	}

	cleanStr := func(s string) string {
		reg, err := regexp.Compile("[^a-zA-Z0-9_\\.\\/\\=\\+\\-\\:\\@\\s]+")
		if err != nil {
			return ""
		}
		return reg.ReplaceAllString(s, "_")
	}

	appendTags := func(input map[string]string, tags map[string]string) map[string]string {

		if len(input) == 0 {
			return tags
		}

		for k, v := range input {
			tags[cleanStr(k)] = cleanStr(v)
		}

		return tags

	}

	o.tags = appendTags(o.job.AppTags, o.tags)
	o.tags = appendTags(o.job.App.Tags, o.tags)

	o.setTagsWithWorkaround(o.job.Metadata.SizeBytes)

}

// This is a workaround for an issue with how metadata is handled when using transfer.UploadManager
// to handle uploads. When UploadManager handles a single-part upload, it automatically prefixes the
// metadata tags with "opc-meta-", which is required by OCI. When handling a multipart upload, it
// does not prefix the metadata tags which can cause the upload to fail. So, here we see if the file
// meets the multipart threshold (set in the UploadRequest above), and if it does, we prefix all the
// tags with "opc-meta-".
func (o *OracleObject) setTagsWithWorkaround(sz int64) {
	tags := make(map[string]string)
	if sz > int64(1024*1024*5) {
		for k, v := range o.tags {
			tags[fmt.Sprintf("%s%s", MultipartTagPrefix, k)] = v
		}
		o.tags = tags
	}
}
