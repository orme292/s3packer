package conf

type ProviderGoogle struct {
	Project      string
	LocationType string
	Storage      string
	ACL          string
	ADC          string
}

func (gc *ProviderGoogle) build(inc *ProfileIncoming) error {
	gc.Project = inc.Google.Project

	gc.LocationType = inc.Google.LocationType

	gc.Storage = inc.Google.Storage

	gc.ACL = inc.Google.ACL

	return nil
}
