package handler_aws

type AwsIterator struct{}

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
