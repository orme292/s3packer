package pack_oci

import (
	"os"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/objectstorage/transfer"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/objectify"
	"github.com/orme292/s3packer/s3packs/provider"
)

func OracleIteratorFunc(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter provider.Iterator, err error) {
	return NewIterator(ac, fol, grp)
}

func NewIterator(ac *conf.AppConfig, fol objectify.FileObjList, grp int) (iter *OracleIterator, err error) {
	_, ociClient, err := buildClients(ac)
	if err != nil {
		return nil, err
	}

	namespace, err := getNamespace(ac, ociClient)
	if err != nil {
		return nil, err
	}

	return &OracleIterator{
		provider: &conf.Provider{
			Is: conf.ProviderNameOCI,
		},
		fol:       fol,
		group:     grp,
		ac:        ac,
		client:    ociClient,
		um:        transfer.NewUploadManager(),
		namespace: namespace,
	}, nil
}

func (oi *OracleIterator) First() (err error) {
	return nil
}

func (oi *OracleIterator) Next() bool {
	if len(oi.fol) == 0 {
		return false
	}

	for {
		if oi.stage.index >= len(oi.fol) {
			return false
		}
		if oi.group != provider.DisregardGroups {
			if oi.fol[oi.stage.index].Group != oi.group {
				oi.stage.index++
				continue
			}
		}
		if oi.fol[oi.stage.index].IsUploaded || oi.fol[oi.stage.index].Ignore {
			oi.ac.Log.Warn("Skipping %q: %s", oi.fol[oi.stage.index].FKey(),
				oi.fol[oi.stage.index].IgnoreString)
			oi.stage.index++
			continue
		}
		break
	}

	f, err := os.Open(oi.fol[oi.stage.index].AbsPath)
	oi.err = err
	oi.stage.f = f

	oi.stage.fo = oi.fol[oi.stage.index]

	return oi.Err() == nil
}

func (oi *OracleIterator) Prepare() *provider.PutObject {
	f := oi.stage.f
	return &provider.PutObject{
		Before: func() error {
			oi.ac.Log.Info("Transferring (%s) %q", objectify.FileSizeString(oi.stage.fo.FileSize),
				oi.stage.fo.FKey())
			return nil
		},
		Object: func() any {
			//tags := make(map[string]string)
			//for k, v := range oi.ac.Tags {
			//	tags[s("%s%s", "opc-meta-", k)] = v
			//}
			return transfer.UploadRequest{
				NamespaceName:                       common.String(oi.namespace),
				BucketName:                          common.String(oi.ac.Bucket.Name),
				ObjectName:                          common.String(oi.stage.fo.FKey()),
				ObjectStorageClient:                 &oi.client,
				EnableMultipartChecksumVerification: common.Bool(true),
				StorageTier:                         oi.ac.Provider.OCI.PutStorage,
				PartSize:                            common.Int64(1024 * 1024 * 5),
				Metadata:                            oi.ac.Tags,
			}
		},
		After: func() error {
			if !oi.stage.fo.IsFailed && !oi.stage.fo.Ignore {
				oi.stage.fo.IsUploaded = true
			}
			oi.stage.index++
			return f.Close()
		},
		Output: func() provider.Object {
			return provider.Object{
				Key:      oi.stage.fo.FKey(),
				Checksum: oi.stage.fo.ChecksumSHA256,
				F:        *f,
			}
		},
		Fo: func() *objectify.FileObj {
			return oi.stage.fo
		},
	}
}

func (oi *OracleIterator) Final() (err error) {
	return nil
}

func (oi *OracleIterator) Err() (err error) {
	if oi.err != nil {
		oi.stage.fo.IsUploaded = false
		oi.ac.Log.Debug("Iterator Error: %q", oi.err.Error())
	}
	return oi.err
}

func (oi *OracleIterator) MarkIgnore(s string) {
	oi.stage.fo.IsUploaded = false
	oi.stage.fo.IgnoreString = s
	oi.stage.fo.Ignore = true
}

func (oi *OracleIterator) MarkFailed(s string) {
	oi.stage.fo.IsUploaded = false
	oi.stage.fo.IsFailedString = s
	oi.stage.fo.IsFailed = true
}
