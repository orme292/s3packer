package pack_aws

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
	"github.com/orme292/s3packer/s3packs/provider"
)

type AwsIterator struct {
	provider *conf.Provider
	svc      *manager.Uploader
	fol      objectify.FileObjList
	stage    struct {
		index int
		fo    *objectify.FileObj
		f     *os.File
	}
	group int
	err   error
	ac    *conf.AppConfig
}

func AwsIteratorFunc(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter provider.Iterator, err error) {
	return NewIterator(ac, fol, grp)
}

func NewIterator(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter *AwsIterator, err error) {
	svc, _, err := buildUploader(ac)
	if err != nil {
		return nil, err
	}
	return &AwsIterator{
		provider: &conf.Provider{
			Is: conf.ProviderNameAWS,
		},
		fol:   fol,
		group: grp,
		ac:    ac,
		svc:   svc,
	}, nil
}

func (ai *AwsIterator) First() (err error) {
	return nil
}

func (ai *AwsIterator) Next() bool {
	if len(ai.fol) == 0 {
		return false
	}

	for {
		if ai.stage.index >= len(ai.fol) {
			return false
		}
		if ai.group != provider.DisregardGroups {
			if ai.fol[ai.stage.index].Group != ai.group {
				ai.stage.index++
				continue
			}
		}
		if ai.fol[ai.stage.index].IsUploaded || ai.fol[ai.stage.index].Ignore {
			ai.ac.Log.Warn("Skipping %q: %s", ai.fol[ai.stage.index].FKey(),
				ai.fol[ai.stage.index].IgnoreString)
			ai.stage.index++
			continue
		}
		break
	}

	f, err := os.Open(ai.fol[ai.stage.index].AbsPath)
	ai.err = err
	ai.stage.f = f
	ai.stage.fo = ai.fol[ai.stage.index]

	return ai.Err() == nil
}

func (ai *AwsIterator) Prepare() *provider.PutObject {
	f := ai.stage.f
	return &provider.PutObject{
		Before: func() error {
			ai.ac.Log.Info("Transferring (%s) %q", objectify.FileSizeString(ai.stage.fo.FileSize),
				ai.stage.fo.FKey())
			return nil
		},
		Object: func() any {
			return &s3.PutObjectInput{
				ACL:               ai.ac.Provider.AwsACL,
				Body:              f,
				Bucket:            aws.String(ai.ac.Bucket.Name),
				ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
				ChecksumSHA256:    aws.String(ai.stage.fo.ChecksumSHA256),
				Key:               aws.String(ai.stage.fo.FKey()),
				StorageClass:      ai.ac.Provider.AwsStorage,
				Tagging:           aws.String(awsTag(ai.stage.fo.TagsMap)),
			}
		},
		After: func() error {
			ai.stage.fo.IsUploaded = true
			ai.stage.index++
			return f.Close()
		},
		Output: func() provider.Object {
			return provider.Object{
				Key:      ai.stage.fo.FKey(),
				Checksum: ai.stage.fo.ChecksumSHA256,
				F:        *f,
			}
		},
		Fo: func() *objectify.FileObj {
			return ai.stage.fo
		},
	}
}

func (ai *AwsIterator) Final() (err error) {
	return nil
}

func (ai *AwsIterator) Err() (err error) {
	if ai.err != nil {
		ai.stage.fo.IsUploaded = false
		ai.ac.Log.Debug("Iterator Error: %q", ai.err.Error())
	}
	return ai.err
}

func (ai *AwsIterator) MarkIgnore(s string) {
	ai.stage.fo.IgnoreString = s
	ai.stage.fo.Ignore = true
}
