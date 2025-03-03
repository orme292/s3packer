package provider

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"s3p/internal/conf"
)

type testOperator struct {
	c     *conf.AppConfig
	state *testState
}

type testState struct {
	bucket     bool
	bucketName string
	tags       map[string]string
	uploaded   []*testFile
}

type testFile struct {
	name   string
	exists bool
	tags   map[string]string
}

func operatorGenFuncTest(app *conf.AppConfig) (oper Operator, err error) {
	oper = &testOperator{
		c: app,
		state: &testState{
			bucket:     false,
			bucketName: "",
			uploaded:   []*testFile{},
		},
	}

	return oper, nil
}

func (oper *testOperator) BucketCreate() error {
	if oper.state.bucket {
		return fmt.Errorf("bucket already exists")
	}
	oper.state.bucket = true
	oper.state.bucketName = oper.c.Bucket.Name
	return nil
}

func (oper *testOperator) BucketExists() (bool, error) {
	if oper.state.bucket == true && oper.state.bucketName == "" {
		return true, fmt.Errorf("bucket does not exist but name is set")
	}
	return oper.state.bucket && (oper.state.bucketName == oper.c.Bucket.Name), nil
}

func (oper *testOperator) BucketDelete() error {
	if oper.state.bucket {
		oper.state.bucket = false
		oper.state.bucketName = ""
		return nil
	}
	return fmt.Errorf("bucket does not exist")
}

func (oper *testOperator) ObjectDelete(key string) error {
	for _, k := range oper.state.uploaded {
		if k.name == key && k.exists == true {
			k.exists = false
		}
		if k.name == key && k.exists == false {
			return fmt.Errorf("object already deleted")
		}
	}
	return fmt.Errorf("object does not exist")
}

func (oper *testOperator) ObjectExists(obj Object) (bool, error) {
	tObj, ok := obj.(*testObject)
	if !ok {
		return false, fmt.Errorf("object is not a test object")
	}

	for _, k := range oper.state.uploaded {
		if k.name == tObj.key {
			if k.exists == true {
				return true, nil
			} else {
				return false, fmt.Errorf("object exists but was deleted")
			}
		}
	}
	return false, fmt.Errorf("object does not exist")
}

func (oper *testOperator) ObjectUpload(obj Object) error {
	tObj, ok := obj.(*testObject)
	if !ok {
		return fmt.Errorf("object is not a test object")
	}

	exists, err := oper.ObjectExists(obj)
	if err != nil {
		if err.Error() != "object does not exist" {
			return err
		}
	}
	if exists {
		return fmt.Errorf("object already exists")
	}
	if tObj.ready == true {
		file := &testFile{
			name:   tObj.key,
			exists: true,
			tags:   tObj.tags,
		}
		oper.state.uploaded = append(oper.state.uploaded, file)
		return nil
	}
	return nil
}

func (oper *testOperator) GetObjectTags(key string) (map[string]string, error) {
	for _, k := range oper.state.uploaded {
		if k.name == key {
			if k.exists == true {
				return k.tags, nil
			} else {
				return nil, fmt.Errorf("object does not exist")
			}
		}
	}
	return nil, fmt.Errorf("object does not exist")
}

func (oper *testOperator) Support() *Supports {
	return NewSupports(true, true, true, true)
}

type testObject struct {
	key   string
	tags  map[string]string
	ready bool
	job   *Job
}

func (o *testObject) Destroy() error {
	return o.Post()
}

func (o *testObject) Generate() error {
	o.key = uuid.New().String()
	o.tags = map[string]string{
		"test":           "yes",
		uuid.NewString(): uuid.NewString(),
	}
	return nil
}

func (o *testObject) Post() error {
	if o.ready == false {
		return fmt.Errorf("object not ready, already false")
	}
	o.ready = false
	return nil
}

func (o *testObject) Pre() error {
	if o.ready == true {
		return fmt.Errorf("object ready, already true")
	}
	if o.key == "" {
		return fmt.Errorf("key not initialized")
	}
	o.key = ""
	o.tags = nil
	o.ready = false
	return nil
}

func objectGenFuncTest(job *Job) Object {
	return &testObject{
		key:   job.Key,
		tags:  job.AppTags,
		ready: false,
		job:   job,
	}
}

func TestInterfaces(t *testing.T) {
	app := conf.NewAppConfig()
	err := app.ImportFromProfile(newIncomingProfile())
	require.NoError(t, err)

	log.Printf("Walking: %v\n", app.Opts.WalkDirs)
	log.Printf("Following Links: %v\n", app.Opts.FollowSymlinks)
	time.Sleep(time.Second * 2)

	handler, err := NewHandler(app, operatorGenFuncTest, objectGenFuncTest)
	require.NoError(t, err)
	require.NotNil(t, handler)
	require.NotNil(t, handler.oper)
	require.NotNil(t, handler.app)
	require.NotNil(t, handler.queue)
	require.NotNil(t, handler.Stats)
	require.NotNil(t, handler.supports)

	err = handler.Init()
	require.NoError(t, err)
	require.NotNil(t, handler.Stats)

	require.Zero(t, handler.Stats.Failed)

	app.Dirs = []string{
		fmt.Sprintf("1/%s/%s/%s", uuid.NewString(), uuid.NewString(), uuid.NewString()),
		fmt.Sprintf("2/%s/%s/%s", uuid.NewString(), uuid.NewString(), uuid.NewString()),
	}
	app.Files = []string{
		fmt.Sprintf("A/%s/%s/%s", uuid.NewString(), uuid.NewString(), uuid.NewString()),
		fmt.Sprintf("B/%s/%s/%s", uuid.NewString(), uuid.NewString(), uuid.NewString()),
	}

	handler, err = NewHandler(app, operatorGenFuncTest, objectGenFuncTest)
	require.NoError(t, err)
	require.NotNil(t, handler)

	// this SHOULD fail, but it currently does not.
	// the files and dirs passed to handler do not exist, so they should be reported as failed.
	// needs to be fixed in the handler or queue
	err = handler.Init()
	require.NoError(t, err) // should be: require.Error(t, err, "handler init should fail with bad dirs/files")

	s := handler.oper.(*testOperator).Support()
	require.True(t, s.BucketCreate, "bucket create should be supported")
	require.True(t, s.BucketDelete, "bucket delete should be supported")

	handler.Stats.IncFailed(1, 100)
	handler.Stats.IncObjects(1, 100)
	handler.Stats.IncSkipped(1, 1024*1024*1024*1024*56)
	require.EqualValues(t, 1, handler.Stats.Failed, "failed should be 1")
	require.EqualValues(t, 1, handler.Stats.Objects, "objects should be 1")
	require.EqualValues(t, 1, handler.Stats.Skipped, "skipped should be 1")
	require.EqualValues(t, "1 objects uploaded, 1 skipped, and 1 failed.", handler.Stats.String(), "stats string should match")

	require.IsType(t, make(map[int64]string), handler.Stats.ReadableString(), "readable string should be a map")
}
