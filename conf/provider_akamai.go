package conf

import (
	"errors"
	"strings"
)

const (
	AkamaiClusterAmsterdam  = "nl-ams-1.linodeobjects.com"
	AkamaiClusterAtlanta    = "us-southeast-1.linodeobjects.com"
	AkamaiClusterChennai    = "in-maa-1.linodeobjects.com"
	AkamaiClusterChicago    = "us-ord-1.linodeobjects.com"
	AkamaiClusterFrankfurt  = "eu-central-1.linodeobjects.com"
	AkamaiClusterJakarta    = "id-cgk-1.linodeobjects.com"
	AkamaiClusterLosAngeles = "us-lax-1.linodeobjects.com"
	AkamaiClusterMiami      = "us-mia-1.linodeobjects.com"
	AkamaiClusterMilan      = "it-mil-1.linodeobjects.com"
	AkamaiClusterNewark     = "us-east-1.linodeobjects.com"
	AkamaiClusterOsaka      = "jp-osa-1.linodeobjects.com"
	AkamaiClusterParis      = "fr-par-1.linodeobjects.com"
	AkamaiClusterSaoPaulo   = "br-gru-1.linodeobjects.com"
	AkamaiClusterSeattle    = "us-sea-1.linodeobjects.com"
	AkamaiClusterSingapore  = "ap-south-1.linodeobjects.com"
	AkamaiClusterStockholm  = "se-sto-1.linodeobjects.com"
	AkamaiClusterAshburn    = "us-iad-1.linodeobjects.com"
)

const (
	AkamaiRegionAmsterdam  = "nl-ams-1"
	AkamaiRegionAtlanta    = "us-southeast-1"
	AkamaiRegionChennai    = "in-maa-1"
	AkamaiRegionChicago    = "us-ord-1"
	AkamaiRegionFrankfurt  = "eu-central-1"
	AkamaiRegionJakarta    = "id-cgk-1"
	AkamaiRegionLosAngeles = "us-lax-1"
	AkamaiRegionMiami      = "us-mia-1"
	AkamaiRegionMilan      = "it-mil-1"
	AkamaiRegionNewark     = "us-east-1"
	AkamaiRegionOsaka      = "jp-osa-1"
	AkamaiRegionParis      = "fr-par-1"
	AkamaiRegionSaoPaulo   = "br-gru-1"
	AkamaiRegionSeattle    = "us-sea-1"
	AkamaiRegionSingapore  = "ap-south-1"
	AkamaiRegionStockholm  = "se-sto-1"
	AkamaiRegionAshburn    = "us-iad-1"
)

func akamaiMatchRegion(region string) (endpoint string, err error) {
	region = strings.ToLower(strings.TrimSpace(region))
	akamaiEndpointsMap := map[string]string{
		AkamaiRegionAmsterdam:  AkamaiClusterAmsterdam,
		AkamaiRegionAtlanta:    AkamaiClusterAtlanta,
		AkamaiRegionChennai:    AkamaiClusterChennai,
		AkamaiRegionChicago:    AkamaiClusterChicago,
		AkamaiRegionFrankfurt:  AkamaiClusterFrankfurt,
		AkamaiRegionJakarta:    AkamaiClusterJakarta,
		AkamaiRegionLosAngeles: AkamaiClusterLosAngeles,
		AkamaiRegionMiami:      AkamaiClusterMiami,
		AkamaiRegionMilan:      AkamaiClusterMilan,
		AkamaiRegionNewark:     AkamaiClusterNewark,
		AkamaiRegionOsaka:      AkamaiClusterOsaka,
		AkamaiRegionParis:      AkamaiClusterParis,
		AkamaiRegionSaoPaulo:   AkamaiClusterSaoPaulo,
		AkamaiRegionSeattle:    AkamaiClusterSeattle,
		AkamaiRegionSingapore:  AkamaiClusterSingapore,
		AkamaiRegionStockholm:  AkamaiClusterStockholm,
		AkamaiRegionAshburn:    AkamaiClusterAshburn,
	}

	endpoint, ok := akamaiEndpointsMap[region]
	if !ok {
		return AkamaiClusterAshburn, errors.New(S("invalid akamai region: %q", region))
	}
	return
}

const (
	ErrorAkamaiKeyOrSecretNotSpecified = "akamai access keys not specified"
)
