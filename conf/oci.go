package conf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oracle/oci-go-sdk/v49/objectstorage"
)

const (
	OciDefaultProfile = "DEFAULT"
)

const (
	ErrorOCICompartmentNotSpecified = "oracle cloud compartment will be tenancy root"
	ErrorOCIAuthNotSpecified        = "oracle cloud auth not specified"
	ErrorOCIStorageNotSpecified     = "oracle cloud storage tier is not valid"
)

const (
	OracleStorageTierStandard         = "standard"
	OracleStorageTierInfrequentAccess = "infrequentaccess" // the case is strange because of the
	OracleStorageTierArchive          = "archive"
)

// ociMatchStorage will match the Storage string to the OCI Storage Tier type.
// The constant values above are used to match the string.
func ociMatchStorage(tier string) (ociTier objectstorage.StorageTierEnum, putTier objectstorage.PutObjectStorageTierEnum, err error) {
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

	ociTier, ok := ociStorageTiersMap[tier]
	putTier, _ = ociPutStorageTiersMap[tier]
	if !ok {
		return objectstorage.StorageTierStandard, objectstorage.PutObjectStorageTierStandard, errors.New(fmt.Sprintf("%s %q", ErrorOCIStorageNotSpecified, tier))
	}
	return ociTier, putTier, nil
}
