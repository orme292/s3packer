package pack_oci

import (
	"context"
	"fmt"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/identity"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
	"github.com/orme292/s3packer/conf"
)

func buildClients(ac *conf.AppConfig) (
	ociIdentity identity.IdentityClient, ociClient objectstorage.ObjectStorageClient, err error) {
	var cp common.ConfigurationProvider
	if ac.Provider.OCI.Profile == OracleDefaultProfile {
		cp = common.DefaultConfigProvider()
	} else {
		cp = common.CustomProfileConfigProvider(EmptyString, ac.Provider.OCI.Profile)
	}

	ociIdentity, err = identity.NewIdentityClientWithConfigurationProvider(cp)
	if err != nil {
		return
	}

	ociClient, err = objectstorage.NewObjectStorageClientWithConfigurationProvider(cp)
	if err != nil {
		return
	}
	return
}

func getNamespace(ac *conf.AppConfig, client objectstorage.ObjectStorageClient) (namespace string, err error) {
	response, err := client.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{
		CompartmentId: &ac.Provider.OCI.Compartment,
	})
	if err != nil {
		return
	}
	return *response.Value, nil
}

func s(f string, v ...any) string {
	return fmt.Sprintf(f, v...)
}
