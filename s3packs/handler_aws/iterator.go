package handler_aws

import (
	"github.com/orme292/s3packer/conf"
	"github.com/orme292/s3packer/s3packs/provider_v2"
)

type AwsIterator struct {
	app *conf.AppConfig
}

func NewAwsIterator(app *conf.AppConfig) (provider_v2.Iterator, error) {

	iter := &AwsIterator{
		app: app,
	}
	return iter, nil

}

func (iter *AwsIterator) Pre() error {
	return nil
}

func (iter *AwsIterator) Next() bool {
	return false
}

func (iter *AwsIterator) Prepare() error {
	return nil
}

func (iter *AwsIterator) Post() error {
	return nil
}
