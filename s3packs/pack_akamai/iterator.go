package pack_akamai

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
	"github.com/orme292/s3packer/s3packs/provider"
)

func AkamaiIteratorFunc(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter provider.Iterator, err error) {
	return NewIterator(ac, fol, grp)
}

func NewIterator(ac *conf.AppConfig, list objectify.FileObjList, grp int) (iter *AkamaiIterator, err error) {
	return &AkamaiIterator{
		provider: &conf.Provider{
			Is: conf.ProviderNameAkamai,
		},
		fol:   list,
		group: grp,
		ac:    ac,
	}, nil
}

func (ai *AkamaiIterator) First() (err error) {
	return nil
}

func (ai *AkamaiIterator) Next() bool {
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

func (ai *AkamaiIterator) Prepare() *provider.PutObject {
	f := ai.stage.f
	return &provider.PutObject{
		Before: func() error {
			ai.ac.Log.Info("Transferring (%s) %q", objectify.FileSizeString(ai.stage.fo.FileSize),
				ai.stage.fo.FKey())
			return nil
		},
		Object: func() any {
			return &s3.PutObjectInput{
				Body:    f,
				Bucket:  aws.String(ai.ac.Bucket.Name),
				Key:     aws.String(ai.stage.fo.FKey()),
				Tagging: aws.String(awsTag(ai.stage.fo.TagsMap)),
			}
		},
		After: func() error {
			if !ai.stage.fo.IsFailed && !ai.stage.fo.Ignore {
				ai.stage.fo.IsUploaded = true
			}
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

func (ai *AkamaiIterator) Final() (err error) {
	return nil
}

func (ai *AkamaiIterator) Err() (err error) {
	if ai.err != nil {
		ai.stage.fo.IsUploaded = false
		ai.ac.Log.Debug("Iterator Error: %q", ai.err.Error())
	}
	return ai.err
}

func (ai *AkamaiIterator) MarkIgnore(s string) {
	ai.stage.fo.IgnoreString = s
	ai.stage.fo.Ignore = true
}

func (ai *AkamaiIterator) MarkFailed(s string) {
	ai.stage.fo.IsUploaded = false
	ai.stage.fo.IsFailed = true
	ai.stage.fo.IsFailedString = s
}
