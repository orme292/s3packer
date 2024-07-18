package oci

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
	"github.com/oracle/oci-go-sdk/v65/objectstorage/transfer"
)

type details struct {
	profile     string
	compartment string
	namespace   string
}

type OracleClient struct {
	obj     *objectstorage.ObjectStorageClient
	id      *identity.IdentityClient
	cp      *common.ConfigurationProvider
	manager *transfer.UploadManager
	details details
}

func NewOracleClient(profile, compartment string) (*OracleClient, error) {

	client := &OracleClient{
		details: details{
			profile:     profile,
			compartment: compartment,
		},
	}

	err := client.Init()
	return client, err

}

func (client *OracleClient) Init() error {

	var (
		cp  common.ConfigurationProvider
		id  identity.IdentityClient
		obj objectstorage.ObjectStorageClient
		err error
	)

	if client.details.profile == OracleDefaultProfile {
		cp = common.DefaultConfigProvider()
	} else {
		cp = common.CustomProfileConfigProvider(EmptyString, client.details.profile)
	}

	id, err = identity.NewIdentityClientWithConfigurationProvider(cp)
	if err != nil {
		return err
	}

	obj, err = objectstorage.NewObjectStorageClientWithConfigurationProvider(cp)
	if err != nil {
		return err
	}

	if client.details.compartment == EmptyString {
		client.details.compartment, err = getTenancyRootOCID(cp)
		if err != nil {
			return err
		}
	}

	client.cp = &cp
	client.id = &id
	client.obj = &obj
	client.manager = transfer.NewUploadManager()

	err = client.getNamespace()
	if err != nil {
		return err
	}

	return nil

}

func (client *OracleClient) getNamespace() error {

	resp, err := client.obj.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{
		CompartmentId: &client.details.compartment,
	})
	if err != nil {
		return err
	}

	client.details.namespace = *resp.Value
	return nil

}

func getTenancyRootOCID(cp common.ConfigurationProvider) (string, error) {
	return cp.TenancyOCID()
}
