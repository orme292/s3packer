package conf

import (
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

// ProviderOCI represents the OCI provider configuration.
type ProviderOCI struct {
	Profile     string
	Compartment string
	Storage     objectstorage.StorageTierEnum

	PutStorage objectstorage.PutObjectStorageTierEnum
}

func (oci *ProviderOCI) build(inc *ProfileIncoming) error {

	err := oci.matchStorage(inc.OCI.Storage)
	if err != nil {
		return err
	}

	oci.Profile = inc.Provider.Profile
	if tidyUpperString(oci.Profile) == OciDefaultProfile {
		oci.Profile = OciDefaultProfile
	}

	oci.Compartment = inc.OCI.Compartment

	return oci.validate()

}

// ociMatchStorage will match the Storage string to the OCI Storage Tier type.
// The constant values above are used to match the string.
func (oci *ProviderOCI) matchStorage(tier string) error {

	tier = strings.ToLower(strings.TrimSpace(tier))

	ociStorageTiersMap := map[string]objectstorage.StorageTierEnum{
		OracleStorageTierStandard:         objectstorage.StorageTierStandard,
		OracleStorageTierInfrequentAccess: objectstorage.StorageTierInfrequentAccess,
		OracleStorageTierArchive:          objectstorage.StorageTierArchive,
	}

	ociPutStorageTiersMap := map[string]objectstorage.PutObjectStorageTierEnum{
		OracleStorageTierStandard:         objectstorage.PutObjectStorageTierStandard,
		OracleStorageTierInfrequentAccess: objectstorage.PutObjectStorageTierInfrequentaccess,
		OracleStorageTierArchive:          objectstorage.PutObjectStorageTierArchive,
	}

	storeTier, ok := ociStorageTiersMap[tier]
	if !ok {
		oci.Storage = objectstorage.StorageTierStandard
		oci.PutStorage = objectstorage.PutObjectStorageTierStandard
		return fmt.Errorf("%s %q", ErrorOCIStorageNotSpecified, tier)
	}

	putTier, _ := ociPutStorageTiersMap[tier]

	oci.Storage = storeTier
	oci.PutStorage = putTier

	return nil

}

func (oci *ProviderOCI) validate() error {

	if oci.Profile == Empty {
		return fmt.Errorf("bad OCI configuration: %v", ErrorOCIAuthNotSpecified)
	}
	return nil

}
