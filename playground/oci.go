package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/oracle/oci-go-sdk/v49/common"
	"github.com/oracle/oci-go-sdk/v49/identity"
	"github.com/oracle/oci-go-sdk/v49/objectstorage"
	"github.com/oracle/oci-go-sdk/v49/objectstorage/transfer"
)

const (
	oraclePrintOut = true
	oracleDefault  = "DEFAULT"

	oracleTestFile1   = "/Users/aorme/Downloads/test/wc058.zip"
	oracleTestDir1    = "/Users/aorme/Downloads/test"
	oracleOcidBuckets = "ocid1.compartment.oc1..aaaaaaaaa2qfwzy3ec6js4cpsslvpda2fkzf5cjcqcrua2ybtyyh3m5636qa"
)

type oracle struct {
	profile  string
	tenancy  string
	provider common.ConfigurationProvider
	client   struct {
		identity identity.IdentityClient
		obj      objectstorage.ObjectStorageClient
	}
	bucket struct {
		name        string
		namespace   *string
		compartment string
		ocid        string
	}
}

func tryOCI() {
	oci := oracle{}
	err := oci.build(oracleDefault)
	if err != nil {
		panic(err)
	}

	_, err = oci.getADs()
	if err != nil {
		panic(err)
	}

	err = oci.setBucketDetails("s3packer-test", oracleOcidBuckets, true)
	if err != nil {
		panic(err)
	}

	err = oci.createBucket()
	if err != nil {
		panic(err)
	}

	_, err = oci.getBuckets()
	if err != nil {
		panic(err)
	}

	files1, err := returnFileNames(oracleTestDir1)
	for _, file := range files1 {
		err = oci.uploadObject(filepath.Base(file), fmt.Sprintf("%s/%s", oracleTestDir1, file))
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}
	for _, file := range files1 {
		err = oci.uploadWithManager(fmt.Sprintf("%s/%s", "parallel", filepath.Base(file)), fmt.Sprintf("%s/%s", oracleTestDir1, file))
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	_, err = oci.getObjects()

	for _, name := range files1 {
		err = oci.deleteObject(name)
		if err != nil {
			panic(err)
		}
	}

	for _, name := range files1 {
		err = oci.deleteObject(fmt.Sprintf("%s/%s", "parallel", filepath.Base(name)))
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
	}

	_, err = oci.getObjects()
	if err != nil {
		panic(err)
	}

	err = oci.deleteBucket()
	if err != nil {
		panic(err)
	}

	_, err = oci.getBuckets()
	if err != nil {
		panic(err)
	}

	fmt.Println("Finished with OCI")
}

func (o *oracle) build(profile string) (err error) {
	o.profile = strings.TrimSpace(strings.ToUpper(profile))

	if o.profile == oracleDefault || o.profile == empty {
		o.profile = oracleDefault
	}

	if o.profile != oracleDefault {
		o.provider = common.CustomProfileConfigProvider(empty, o.profile)
	} else {
		o.provider = common.DefaultConfigProvider()
	}

	o.client.identity, err = identity.NewIdentityClientWithConfigurationProvider(o.provider)
	if err != nil {
		return err
	}
	o.client.obj, err = objectstorage.NewObjectStorageClientWithConfigurationProvider(o.provider)
	if err != nil {
		return err
	}

	o.tenancy, err = o.provider.TenancyOCID()
	if err != nil {
		return err
	}

	response, err := o.client.obj.GetNamespace(context.Background(), objectstorage.GetNamespaceRequest{
		CompartmentId: &o.tenancy,
	})
	if err != nil {
		return err
	}
	o.bucket.namespace = response.Value

	return
}

func (o *oracle) getADs() (ads []identity.AvailabilityDomain, err error) {
	var compartment string
	if o.bucket.compartment == empty {
		compartment = o.tenancy
	} else {
		compartment = o.bucket.compartment
	}

	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: &compartment,
	}

	r, err := o.client.identity.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		return
	}

	ads = r.Items

	if oraclePrintOut {
		fmt.Println("List of availability domains:")
		for _, ad := range ads {
			println("AD: ", *ad.Name)
		}
		fmt.Println(empty)
	}
	return
}

func (o *oracle) setBucketDetails(name string, compartment string, randomizeName bool) (err error) {
	if name == empty {
		return errors.New("bucket name cannot be empty")
	}
	if compartment == empty {
		compartment = o.tenancy
	}

	o.bucket.name = name
	o.bucket.compartment = compartment

	if randomizeName {
		o.bucket.name = fmt.Sprintf("%s_%s", name, randomString(8))
	}
	return
}

func (o *oracle) createBucket() (err error) {
	request := objectstorage.CreateBucketRequest{
		NamespaceName: o.bucket.namespace,
		CreateBucketDetails: objectstorage.CreateBucketDetails{
			CompartmentId: &o.bucket.compartment,
			Name:          &o.bucket.name,
			Metadata: map[string]string{
				"creator": "s3packer-playground",
			},
		},
	}

	r, err := o.client.obj.CreateBucket(context.Background(), request)
	if err != nil {
		return err
	}

	if oraclePrintOut {
		fmt.Print("Create bucket: ", *r.Name, " - ", *r.CompartmentId, " - ", *r.Namespace, "...")
	}
	o.bucket.ocid = *r.CompartmentId
	if oraclePrintOut {
		fmt.Print("done\n")
	}
	return
}

func (o *oracle) getBuckets() (buckets []objectstorage.BucketSummary, err error) {
	request := objectstorage.ListBucketsRequest{
		NamespaceName: o.bucket.namespace,
		CompartmentId: &o.bucket.compartment,
	}

	r, err := o.client.obj.ListBuckets(context.Background(), request)
	if err != nil {
		return
	}

	buckets = r.Items
	if oraclePrintOut {
		fmt.Println("\nList of buckets:")
		for _, b := range buckets {
			println("Bucket: ", *b.Name)
		}
		fmt.Println("Total Buckets: ", len(buckets))
		fmt.Println(empty)
	}
	return
}

func (o *oracle) deleteBucket() (err error) {
	request := objectstorage.DeleteBucketRequest{
		NamespaceName: o.bucket.namespace,
		BucketName:    &o.bucket.name,
	}

	if oraclePrintOut {
		fmt.Print("Deleting bucket: ", o.bucket.name, "...")
	}
	_, err = o.client.obj.DeleteBucket(context.Background(), request)
	if oraclePrintOut {
		fmt.Print("done\n")
	}
	return
}

func (o *oracle) uploadObject(name string, file string) (err error) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return
	}

	request := objectstorage.PutObjectRequest{
		NamespaceName: o.bucket.namespace,
		BucketName:    &o.bucket.name,
		PutObjectBody: f,
		ObjectName:    &name,
		StorageTier:   objectstorage.PutObjectStorageTierStandard,
	}

	if oraclePrintOut {
		fmt.Print("Uploading object: ", name, "...")
	}
	_, err = o.client.obj.PutObject(context.Background(), request)
	if oraclePrintOut {
		fmt.Print("done\n")
	}
	return
}

func (o *oracle) getObjects() (objects []objectstorage.ObjectSummary, err error) {
	r_fields := "name,size,etag,timeCreated,md5,timeModified,storageTier,archivalState"
	request := objectstorage.ListObjectsRequest{
		NamespaceName: o.bucket.namespace,
		BucketName:    &o.bucket.name,
		Fields:        &r_fields,
	}

	r, err := o.client.obj.ListObjects(context.Background(), request)
	if err != nil {
		return
	}

	objects = r.Objects
	if oraclePrintOut {
		fmt.Println("\nList of objects:")
		for _, obj := range objects {
			fmt.Printf("Object: %s\n", obj.String())
		}
		fmt.Println("Total Objects: ", len(objects))
		fmt.Println(empty)
	}
	return
}

func (o *oracle) deleteObject(name string) (err error) {
	request := objectstorage.DeleteObjectRequest{
		NamespaceName: o.bucket.namespace,
		BucketName:    &o.bucket.name,
		ObjectName:    &name,
	}

	if oraclePrintOut {
		fmt.Print("Deleting object: ", name, "...")
	}
	_, err = o.client.obj.DeleteObject(context.Background(), request)
	if oraclePrintOut {
		fmt.Print("done\n")
	}

	return
}

func (o *oracle) uploadWithManager(name string, path string) (err error) {
	path, _ = filepath.Abs(path)
	um := transfer.NewUploadManager()

	fmt.Print("Uploading file: ", name, "...")
	r, err := um.UploadFile(context.Background(), transfer.UploadFileRequest{
		UploadRequest: transfer.UploadRequest{
			NamespaceName:                       o.bucket.namespace,
			BucketName:                          &o.bucket.name,
			ObjectName:                          common.String(name),
			EnableMultipartChecksumVerification: common.Bool(true),
			StorageTier:                         objectstorage.PutObjectStorageTierStandard,
			PartSize:                            common.Int64(1024 * 1024 * 5),
		},
		FilePath: path,
	})
	fmt.Print("done\n")
	if r.SinglepartUploadResponse != nil {
		fmt.Println("Singlepart Upload Response: ", r.SinglepartUploadResponse.ETag)
	}
	if r.MultipartUploadResponse != nil {
		fmt.Println("Multipart Upload Response: ", r.MultipartUploadResponse.UploadID, ":", r.MultipartUploadResponse.ETag)
	}

	fmt.Println("\n")
	return
}
