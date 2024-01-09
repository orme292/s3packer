package pack_aws

import (
	"bytes"
	"context"
	"errors"
	"os"
	"sort"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/orme292/s3packer/s3packs/objectify"
	"github.com/orme292/s3packer/s3packs/provider"
)

func (op *AwsOperator) SupportsMultipartUploads() bool { return true }

func (op *AwsOperator) UploadMultipart(po provider.PutObject) (err error) {
	var start int64 = 0
	var bufferSize int64 = 10000000
	var i int32 = 1
	file, _ := os.Open(po.Fo().AbsPath)
	fi, _ := file.Stat()
	size := fi.Size()

	op.ctl = &MultipartControl{
		max:   5,
		retry: 5,
	}
	op.ctl.upload = make(map[int]*mpu)
	op.ctl.ctx, op.ctl.cancel = context.WithCancel(context.Background())

	obj, ok := po.Object().(*s3.PutObjectInput)
	if !ok {
		return errors.New(ErrorCouldNotAssertObject)
	}
	op.ctl.obj = obj

	op.ctl.cmo, _ = op.client.CreateMultipartUpload(op.ctl.ctx, &s3.CreateMultipartUploadInput{
		Bucket:            aws.String(op.ac.Bucket.Name),
		Key:               obj.Key,
		ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
		ACL:               obj.ACL,
		StorageClass:      obj.StorageClass,
		Tagging:           obj.Tagging,
	})

	op.ctl.uploadId = *op.ctl.cmo.UploadId

	for start < size {
		op.ctl.upload[int(i)] = &mpu{}
		bufferSize = min(size-start, bufferSize)
		op.ctl.upload[int(i)].data = make([]byte, bufferSize)
		_, _ = file.Read(op.ctl.upload[int(i)].data)

		op.ctl.upload[int(i)].cs, err = objectify.GetChecksumSHA256Reader(bytes.NewReader(op.ctl.upload[int(i)].data))
		if err != nil {
			_ = op.MultipartAbort()
			return err
		}
		op.ctl.upload[int(i)].index = int(i)
		op.ctl.upload[int(i)].input = &s3.UploadPartInput{
			Bucket:            aws.String(op.ac.Bucket.Name),
			Key:               obj.Key,
			PartNumber:        aws.Int32(i),
			Body:              bytes.NewReader(op.ctl.upload[int(i)].data),
			UploadId:          aws.String(op.ctl.uploadId),
			ChecksumAlgorithm: types.ChecksumAlgorithmSha256,
			ChecksumSHA256:    aws.String(op.ctl.upload[int(i)].cs),
			ContentLength:     aws.Int64(bufferSize),
		}
		start += bufferSize
		i++
	}

	if len(op.ctl.upload) < op.ctl.max {
		op.ctl.max = len(op.ctl.upload)
	}

	err = op.MultipartParallelize()
	if err != nil {
		op.ac.Log.Error("Error uploading %q: %q", obj.Key, err.Error())
	}
	return op.MultipartFinalize()
}

func (op *AwsOperator) MultipartParallelize() (err error) {
	op.ac.Log.Info("Performing Multipart Parallel Upload...")
	var wg sync.WaitGroup
	errChan := make(chan error)
	for i := 0; i < op.ctl.max; i++ {
		for index := range op.ctl.upload {
			op.ctl.upload[index].group = index % op.ctl.max
		}
	}

	for i := 0; i < op.ctl.max; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, group int) {
			defer wg.Done()
			for index := range op.ctl.upload {
				if op.ctl.upload[index].group == group {
					complete, output, oErr := op.MultipartUpload(index)
					op.ctl.upload[index].complete = complete
					op.ctl.upload[index].output = output
					op.ctl.upload[index].err = oErr
					if oErr != nil {
						errChan <- oErr
						op.ctl.cancel()
					}
				}
			}
		}(&wg, i)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(errChan)
	}(&wg)

	if err = <-errChan; err != nil {
		return op.MultipartAbort()
	}
	return
}

func (op *AwsOperator) MultipartUpload(i int) (complete bool, output *s3.UploadPartOutput, err error) {
	op.ac.Log.Debug("Uploading part %d", i)
	retry := op.ctl.retry + 1
	for retry > 0 {
		if retry != op.ctl.retry+1 {
			op.ac.Log.Debug("Retrying part %d", i)
		}
		output, err = op.client.UploadPart(op.ctl.ctx, op.ctl.upload[i].input)
		if err != nil {
			op.ac.Log.Error("Problem: %q", i, err.Error())
			retry--
			continue
		} else {
			break
		}
	}
	if retry == 0 {
		return false, output, errors.New("couldn't complete upload")
	}
	op.ctl.upload[i].etag = *output.ETag
	return true, output, nil
}

func (op *AwsOperator) MultipartFinalize() (err error) {
	keys := make([]int, 0, len(op.ctl.upload))
	for k := range op.ctl.upload {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	parts := make([]types.CompletedPart, len(keys))
	for i, k := range keys {
		parts[i] = types.CompletedPart{
			ETag:           aws.String(op.ctl.upload[k].etag),
			PartNumber:     aws.Int32(int32(op.ctl.upload[k].index)),
			ChecksumSHA256: aws.String(op.ctl.upload[k].cs),
		}
	}

	_, err = op.client.CompleteMultipartUpload(op.ctl.ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(op.ac.Bucket.Name),
		Key:      op.ctl.obj.Key,
		UploadId: aws.String(op.ctl.uploadId),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	op.ac.Log.Debug("Completed Multipart Upload: %q", *op.ctl.obj.Key)
	return
}
func (op *AwsOperator) MultipartAbort() (err error) {
	op.ac.Log.Error("Aborting multipart upload: %q" + op.ctl.uploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(op.ac.Bucket.Name),
		Key:      op.ctl.obj.Key,
		UploadId: aws.String(op.ctl.uploadId),
	}
	_, err = op.client.AbortMultipartUpload(op.ctl.ctx, abortInput)
	if err != nil {
		op.ac.Log.Fatal("Error aborting multipart upload: %q", err.Error())
	}
	return err
}
