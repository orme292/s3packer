package conf

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// ProviderLinode represents the Linode/Akamai provider configuration
type ProviderLinode struct {
	Key       string
	Secret    string
	Endpoint  string
	BucketACL types.BucketCannedACL
	ObjectACL types.ObjectCannedACL
}

var linodeEndpointsMap = map[string]string{
	LinodeRegionAmsterdam:  LinodeClusterAmsterdam,
	LinodeRegionAtlanta:    LinodeClusterAtlanta,
	LinodeRegionChennai:    LinodeClusterChennai,
	LinodeRegionChicago:    LinodeClusterChicago,
	LinodeRegionFrankfurt:  LinodeClusterFrankfurt,
	LinodeRegionJakarta:    LinodeClusterJakarta,
	LinodeRegionLosAngeles: LinodeClusterLosAngeles,
	LinodeRegionMiami:      LinodeClusterMiami,
	LinodeRegionMilan:      LinodeClusterMilan,
	LinodeRegionNewark:     LinodeClusterNewark,
	LinodeRegionOsaka:      LinodeClusterOsaka,
	LinodeRegionParis:      LinodeClusterParis,
	LinodeRegionSaoPaulo:   LinodeClusterSaoPaulo,
	LinodeRegionSeattle:    LinodeClusterSeattle,
	LinodeRegionSingapore:  LinodeClusterSingapore,
	LinodeRegionStockholm:  LinodeClusterStockholm,
	LinodeRegionAshburn:    LinodeClusterAshburn,
}

func (l *ProviderLinode) build(inc *ProfileIncoming) error {

	err := l.matchRegion(inc.Linode.Region)
	if err != nil {
		return err
	}

	l.BucketACL = types.BucketCannedACLPrivate
	l.ObjectACL = types.ObjectCannedACLPublicRead

	l.Key = inc.Provider.Key
	l.Secret = inc.Provider.Secret

	return l.validate()

}

func (l *ProviderLinode) matchRegion(region string) error {

	endpoint, ok := linodeEndpointsMap[tidyLowerString(region)]
	if !ok {
		l.Endpoint = LinodeClusterAshburn
		return fmt.Errorf("%s, %q", LinodeInvalidRegion, region)
	}
	l.Endpoint = endpoint

	return nil

}

func (l *ProviderLinode) validate() error {

	if l.Secret == Empty || l.Key == Empty {
		return fmt.Errorf("bad Linode config: %v", LinodeAuthNeeded)
	}

	if l.Endpoint == "" {
		return fmt.Errorf("bad Linode config: %v", LinodeInvalidRegion)
	}

	return nil

}
