package gcloud

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider"
)

type GoogleOperator struct {
	App   *conf.AppConfig
	Cloud *GoogleClient
}

func NewGCloudOperator(app *conf.AppConfig) (oper provider.Operator, err error) {

	client := &GoogleClient{
		cfg: &googleCfg{},
	}

	client.cfg.adc = app.Provider.Google.ADC

	err = client.getClient()
	if err != nil {
		return nil, err
	}

	client.getBucket(app.Bucket.Name)

	return &GoogleOperator{
		App:   app,
		Cloud: client,
	}, nil

}

func (oper *GoogleOperator) BucketCreate() error {

	if strings.TrimSpace(oper.App.Provider.Google.Project) == EmptyString {
		return fmt.Errorf("project name is required")
	}

	attrs := &storage.BucketAttrs{
		Name:                       oper.App.Bucket.Name,
		PredefinedACL:              oper.App.Provider.Google.BucketACL,
		PredefinedDefaultObjectACL: oper.App.Provider.Google.BucketACL,
		Location:                   oper.App.Bucket.Region,
		StorageClass:               oper.App.Provider.Google.Storage,
	}

	if oper.App.Provider.Google.LocationType == conf.GCLocationTypeDual {
		var locStr []string
		locations := strings.Split(oper.App.Bucket.Region, ",")
		for str := range locations {
			locStr = append(locStr, strings.TrimSpace(locations[str]))
		}

		if len(locStr) != 2 {
			return fmt.Errorf("dual location requires 2 locations")
		}

		attrs.CustomPlacementConfig = &storage.CustomPlacementConfig{
			DataLocations: locStr,
		}
		attrs.Location = EmptyString
	}

	if err := oper.Cloud.Bucket.Create(oper.Cloud.Ctx, oper.App.Provider.Google.Project, attrs); err != nil {
		return fmt.Errorf("error creating bucket: %s", err.Error())
	}

	oper.Cloud.refreshBucket()

	acl := oper.Cloud.Bucket.ACL()
	if err := acl.Set(oper.Cloud.Ctx, storage.AllAuthenticatedUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("error setting ACL: %s", err.Error())
	}

	return nil

}

func (oper *GoogleOperator) BucketExists() (bool, error) {

	// Why is there no function to check if a bucket exists?
	_, err := oper.Cloud.Bucket.Attrs(oper.Cloud.Ctx)
	if errors.Is(err, storage.ErrBucketNotExist) {
		return false, fmt.Errorf("bucket not found")
	}
	if err != nil {
		return false, fmt.Errorf("error trying to find bucket: %s", err.Error())
	}

	return true, nil

}

func (oper *GoogleOperator) BucketDelete() error {
	if err := oper.Cloud.Bucket.Delete(oper.Cloud.Ctx); err != nil {
		return fmt.Errorf("error deleting bucket: %s", err.Error())
	}

	return nil
}

func (oper *GoogleOperator) ObjectDelete(key string) error {
	// not supported
	return nil
}

func (oper *GoogleOperator) ObjectExists(obj provider.Object) (bool, error) {

	gcObj, ok := obj.(*CloudObject)
	if !ok {
		return true, fmt.Errorf("trouble building object to check")
	}

	_, err := oper.Cloud.Bucket.Object(gcObj.key).Attrs(oper.Cloud.Ctx)
	if errors.Is(err, storage.ErrObjectNotExist) {
		return false, fmt.Errorf("object not found")
	}
	if err != nil {
		oper.App.Tui.Error(err.Error())
		return true, fmt.Errorf("error trying to find object: %s", err.Error())
	}

	return true, nil

}

func (oper *GoogleOperator) ObjectUpload(obj provider.Object) error {

	gcObj, ok := obj.(*CloudObject)
	if !ok {
		return fmt.Errorf("trouble building object to upload")
	}

	if gcObj.job.Metadata.HasChanged() {
		return fmt.Errorf("file changed during upload: %s", gcObj.job.Metadata.FullPath())
	}

	var attempt int
	maxAttempts := 5
	backoff := time.Second
	for {

		wc := oper.Cloud.Bucket.Object(gcObj.key).NewWriter(oper.Cloud.Ctx)
		defer func() {
			if err := wc.Close(); err != nil {
				oper.App.Tui.Warn(err.Error())
			}
		}()

		wc.ObjectAttrs = storage.ObjectAttrs{
			Name:          gcObj.key,
			PredefinedACL: oper.App.Provider.Google.ObjectACL,
			Metadata:      gcObj.tagMap,
		}

		_, err := io.Copy(wc, gcObj.f)
		if err == nil {
			break
		}

		if attempt >= maxAttempts {
			return fmt.Errorf("error uploading [%s]: %s", err.Error(), gcObj.key)
		}

		attempt++
		dur := backoff * time.Duration(1<<attempt)
		oper.App.Tui.Info("Retrying upload in ", dur, " seconds. Attempt ", attempt, " of ", maxAttempts, "")
		time.Sleep(dur)
	}

	return nil

}

func (oper *GoogleOperator) GetObjectTags(key string) (map[string]string, error) {
	return make(map[string]string), nil
}

func (oper *GoogleOperator) Support() *provider.Supports {
	return provider.NewSupports(true, true, false, false)
}
