package conf

import (
	"fmt"
)

// ProviderName subtype, for quickly matching providers
type ProviderName string

const (
	ProviderNameNone   ProviderName = "none"
	ProviderNameAWS    ProviderName = "aws"
	ProviderNameOCI    ProviderName = "oci"
	ProviderNameLinode ProviderName = "linode"
)

func (pn ProviderName) String() string {
	return string(pn)
}

// Provider represents the configuration for a provider.
//
// Fields:
// - Is (ProviderName): The name of the provider. (e.g., "AWS", "OCI")
// - AWS (*ProviderAWS): The configuration for AWS.
// - OCI (*ProviderOCI): The configuration for OCI.
// - Key (string): The provider key.
// - Secret (string): The provider secret.
//
// Usage examples can be found in the surrounding code.
type Provider struct {
	Is     ProviderName
	AWS    *ProviderAWS
	OCI    *ProviderOCI
	Linode *ProviderLinode
}

func (p *Provider) build(inc *ProfileIncoming) error {

	p.match(inc.Provider.Use)
	if p.Is == ProviderNameNone {
		return fmt.Errorf("error loading profile: %v", ErrorProviderNotSpecified)
	}

	switch p.Is {

	case ProviderNameAWS:
		p.AWS = &ProviderAWS{}
		return p.AWS.build(inc)

	case ProviderNameOCI:
		p.OCI = &ProviderOCI{}
		return p.OCI.build(inc)

	case ProviderNameLinode:
		p.Linode = &ProviderLinode{}
		return p.Linode.build(inc)

	default:
		return fmt.Errorf("could not build profile: %v", ErrorProviderNotSpecified)

	}

}

func (p *Provider) match(s string) {

	switch tidyLowerString(s) {
	case "aws", "amazon", "s3", "amazon s3":
		p.Is = ProviderNameAWS
	case "oci", "oracle", "oracle cloud":
		p.Is = ProviderNameOCI
	case "akamai", "linode", "linode objects":
		p.Is = ProviderNameLinode
	default:
		p.Is = ProviderNameNone
	}

}
