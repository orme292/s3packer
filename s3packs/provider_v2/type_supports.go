package provider_v2

type Supports struct {
	BucketCreate bool
	BucketDelete bool
	ObjectDelete bool
}

func NewSupports(bucketCreate, bucketDelete, objectDelete bool) *Supports {
	return &Supports{
		BucketCreate: bucketCreate,
		BucketDelete: bucketDelete,
		ObjectDelete: objectDelete,
	}
}
