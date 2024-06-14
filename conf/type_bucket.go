package conf

import (
	"fmt"
)

// Bucket contains all details related to the bucket, for any provider. Create is not implemented.
type Bucket struct {
	Create bool
	Name   string
	Region string
}

func (b *Bucket) build(inc *ProfileIncoming, pn ProviderName) error {

	b.Name = inc.Bucket.Name
	b.Create = inc.Bucket.Create
	b.Region = inc.Bucket.Region

	return b.validate(pn)
}

func (b *Bucket) validate(pn ProviderName) error {

	if b.Name == Empty {
		return fmt.Errorf("bad bucket config: %v", ErrorBucketInfoNotSpecified)
	}
	if b.Region == Empty && (pn == ProviderNameAWS || pn == ProviderNameLinode) {
		return fmt.Errorf("bad bucket config: %v", ErrorBucketInfoNotSpecified)
	}

	return nil

}
