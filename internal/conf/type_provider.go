package conf

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ProviderName subtype, for quickly matching providers
type ProviderName string

const (
	ProviderNameNone   ProviderName = "none"
	ProviderNameAWS    ProviderName = "aws"
	ProviderNameOCI    ProviderName = "oci"
	ProviderNameLinode ProviderName = "linode"
	ProviderNameGoogle ProviderName = "google"
)

func (pn ProviderName) String() string {
	return strings.ToLower(string(pn))
}

func (pn ProviderName) Title() string {
	caser := cases.Title(language.English)
	return caser.String(string(pn))
}

func (pn ProviderName) Match(s string) bool {
	return pn.String() == s
}

// Provider represents the configuration for a provider.
//
// Fields:
// - Is (ProviderName): The name of the provider. (e.g., "AWS", "OCI")
// - AWS (*ProviderAWS): The configuration for AWS.
// - Google (*ProviderGoogle): The configuration for Google Cloud.
// - Linode (*ProviderLinode): The configuration for Linode.
// - OCI (*ProviderOCI): The configuration for OCI.
// - Key (string): The provider key.
// - Secret (string): The provider secret.
//
// Usage examples can be found in the surrounding code.
type Provider struct {
	Is     ProviderName
	AWS    *ProviderAWS
	Google *ProviderGoogle
	Linode *ProviderLinode
	OCI    *ProviderOCI
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

	case ProviderNameGoogle:
		p.Google = &ProviderGoogle{
			ADC: inc.Provider.Profile,
		}
		return p.Google.build(inc)

	case ProviderNameLinode:
		p.Linode = &ProviderLinode{}
		return p.Linode.build(inc)

	case ProviderNameOCI:
		p.OCI = &ProviderOCI{}
		return p.OCI.build(inc)

	default:
		return fmt.Errorf("could not build profile: %v", ErrorProviderNotSpecified)

	}

}

func (p *Provider) match(s string) {

	switch tidyLowerString(s) {
	case ProviderNameAWS.String(), "amazon", "s3", "amazon s3":
		p.Is = ProviderNameAWS
	case ProviderNameGoogle.String(), "gcloud", "google cloud", "google cloud storage":
		p.Is = ProviderNameGoogle
	case ProviderNameLinode.String(), "akamai", "linode objects":
		p.Is = ProviderNameLinode
	case ProviderNameOCI.String(), "oracle", "oracle cloud":
		p.Is = ProviderNameOCI
	default:
		p.Is = ProviderNameNone
	}

}
