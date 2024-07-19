package oci

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/orme292/s3packer/s3packs/provider"
)

type OracleObject struct {
	job *provider.Job

	key    string
	bucket string

	tags map[string]string

	multipartThreshold int64
}

func NewOracleObject(job *provider.Job) provider.Object {
	return &OracleObject{
		job:                job,
		multipartThreshold: 1024 * 1024 * 5,
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
		reg := regexp.MustCompile(`[^a-zA-Z0-9_./=\+\-:@\s]+`)
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

// This is a workaround for an issue with how metadata is handled when using transfer.UploadManager.
// When UploadManager handles a single-part upload, it automatically prefixes the
// metadata tags with the required prefix "opc-meta-". When handling a multipart upload, it
// does not prefix the metadata tags and causes the upload to fail. Here we see if the file
// meets the multipart threshold (50MiB unless set otherwise in the UploadRequest), and if it does,
// we prefix all the tags with "opc-meta-", or the value of the const MetadataTagPrefix.
func (o *OracleObject) setTagsWithWorkaround(sz int64) {
	tags := make(map[string]string)
	if sz > o.multipartThreshold {
		for k, v := range o.tags {
			if strings.HasPrefix(strings.ToLower(k), MetadataTagPrefix) == false {
				tags[fmt.Sprintf("%s%s", MetadataTagPrefix, k)] = v
			}
		}
		o.tags = tags
	}
}
