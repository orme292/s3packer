package provider

type Supports struct {
	BucketCreate  bool
	BucketDelete  bool
	ObjectDelete  bool
	GetObjectTags bool
}

func NewSupports(bucketCreate, bucketDelete, objectDelete, getObjectTags bool) *Supports {
	return &Supports{
		BucketCreate:  bucketCreate,
		BucketDelete:  bucketDelete,
		ObjectDelete:  objectDelete,
		GetObjectTags: getObjectTags,
	}
}
