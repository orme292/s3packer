package config

import (
	"errors"
	"fmt"
	"strings"
)

/*
isString checks if the given interface{} is a string
*/
func isString(s any) bool {
	switch s.(type) {
	case string:
		return true
	default:
		return false
	}
}

/*
isBool checks if the given interface{} is a bool
*/
func isBool(b any) bool {
	switch b.(type) {
	case bool:
		return true
	default:
		return false
	}
}

/*
Validate validates the configuration file. It runs validation checks in the appropriate order.
  - createMissingMaps() initializes the maps in the configuration file if they are nil.
  - criticalMissingValues() checks if required fields are missing from the profile.
  - repairMissingFields() checks if a field is a different type then expected. If it is, it is cleared and set
    to a default value.
  - sanitizeACL() validates the value for the ACL field in the profile.
  - sanitizeKeyNamingMethod() validates the value for the keyNamingMethod field in the profile.
  - sanitizeTags() validates the values for the tags field in the profile. If a key or value in a tag contains an
    equal sign, then it is removed.
  - sanitizeStorageType validates the value for the storage field in the profile.
  - sanitizeDirList() removes trailing slashes from the directories in the Dirs slice.
*/
func (c *Configuration) Validate() error {
	c.createMissingMaps()
	c.repairMissingFields()
	c.sanitizeACL()
	c.sanitizeKeyNamingMethod()
	c.sanitizePrefixes()
	c.sanitizeTags()
	c.sanitizeStorageType()
	c.sanitizeDirList()
	err := c.criticalMissingValues()
	if err != nil {
		return err
	}
	return nil
}

/*
createMissingMaps initializes the maps in the configuration file if they are nil.
The maps will be nil if the fields are not present in the YAML profile.
*/
func (c *Configuration) createMissingMaps() {
	if c.Version < 0 || c.Version > 1 {
		c.Version = 1
	}
	if c.Authentication == nil {
		c.Authentication = map[string]any{}
	}
	if c.Bucket == nil {
		c.Bucket = map[string]any{}
	}
	if c.Naming == nil {
		c.Naming = map[string]any{}
	}
	if c.Options == nil {
		c.Options = map[string]any{}
	}
	if c.Logging == nil {
		c.Logging = map[string]any{}
	}
	if c.Files == nil {
		c.Files = []string{}
	}
	if c.Dirs == nil {
		c.Dirs = []string{}
	}
	if c.Tags == nil {
		c.Tags = map[string]string{}
	}
}

/*
criticalMissingValues() checks if required fields are missing from the profile.
You must run createMissingMaps() before using this function.
*/
func (c *Configuration) criticalMissingValues() error {
	if c.Authentication[ProfileAuthProfile].(string) == EmptyString {
		if c.Authentication[ProfileAuthKey].(string) == EmptyString || c.Authentication[ProfileAuthSecret].(string) == EmptyString {
			return errors.New("authentication details not provided")
		}
	}
	if c.Authentication[ProfileAuthProfile].(string) != EmptyString {
		if c.Authentication[ProfileAuthKey].(string) != EmptyString || c.Authentication[ProfileAuthSecret].(string) != EmptyString {
			return errors.New("multiple authentication method details are provided")
		}
	}
	if c.Options[ProfileOptionKeyNamingMethod].(string) == EmptyString {
		return errors.New("keyNamingMethod is empty, should be relative or absolute")
	}
	if c.Bucket[ProfileBucketName].(string) == EmptyString {
		return errors.New("bucket name not provided")
	}
	if c.Bucket[ProfileBucketRegion].(string) == EmptyString {
		return errors.New("bucket region not provided")
	}
	if c.Logging[ProfileLoggingToFile].(bool) == true && c.Logging[ProfileLoggingFilename].(string) == EmptyString {
		return errors.New("logging toFile is true but filename not provided")
	}
	return nil
}

/*
sanitizeACL() validates the value for the ACL field in the profile.

	AWS Canned ACLs
	c.Options["acl"] takes an aws canned ACL value:
	- private: Owner gets FULL_CONTROL. No one else has access rights (default).
	- public-read: Owner gets FULL_CONTROL. The AllUsers group gets READ access.
	- public-read-write: Owner gets FULL_CONTROL. The AllUsers group gets READ and WRITE access.
	- authenticated-read: Owner gets FULL_CONTROL. Amazon EC2 gets READ access to GET an AMI bundle.
	- aws-exec-read: Owner gets FULL_CONTROL. The AuthenticatedUsers group gets READ access.
	- bucket-owner-read: Object owner gets FULL_CONTROL. Bucket owner gets READ access.
	- bucket-owner-full-control: Both the object owner and the bucket owner get FULL_CONTROL over the object.
	- log-delivery-write: The LogDelivery group gets WRITE and READ_ACP permissions on the bucket.

	See: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl for more information.

See CONST values in types.go to match the constants with the canned ACL values.
*/
func (c *Configuration) sanitizeACL() {
	c.Options[ProfileOptionACL] = strings.ToLower(strings.TrimSpace(c.Options[ProfileOptionACL].(string)))
	switch c.Options[ProfileOptionACL].(string) {
	case ACLPrivate:
	case ACLPublicRead:
	case ACLPublicReadWrite:
	case ACLAuthenticatedRead:
	case ACLAwsExecRead:
	case ACLBucketOwnerRead:
	case ACLBucketOwnerFullControl:
	case ACLLogDeliveryWrite:
	case EmptyString:
		c.Logger.Warn("No ACL specified, using default.")
		c.Options[ProfileOptionACL] = ACLPrivate
	default:
		c.Logger.Warn("Invalid ACL specified, using default.")
		c.Options[ProfileOptionACL] = ACLPrivate
	}
	c.Logger.Info("ACL set to " + c.Options[ProfileOptionACL].(string))
}

func (c *Configuration) sanitizeDirList() {
	if len(c.Dirs) == 0 {
		return
	}

	var trimmedDirs []string
	for index := range c.Dirs {
		trimmedDirs = append(trimmedDirs, strings.TrimRight(c.Dirs[index], "/"))
	}
	c.Dirs = trimmedDirs
}

/*
sanitizeKeyNamingMethod() validates the value for the keyNamingMethod field in the profile.

	c.Options["keyNamingMethod"] takes one of these naming methods
	- relative: ALL objects are uploaded to the bucket root, objects will be named with the
	  pathPrefix + relative local path + objectPrefix + filename. (individual file uploads are always at the root)
	- absolute: ALL objects are uploaded to the bucket root, objects will be named with the pathPrefix +
	  absolute local path + objectPrefix + filename.

	Examples (when uploading local directories):
	- relative: (individual files)
				/home/user/my_file.doc -> /my_file.doc
				/home/user/downloads/archive.tar.gz -> /archive.tar.gz
		   		(directories)
				/home/user/ ...
					/home/user/my_file.doc -> /my_file.doc
					/home/user/downloads/archive.tar.gz -> /downloads/archive.tar.gz
	- absolute: (individual files)
				/home/user/my_file.doc -> /home/user/my_file.doc
				/home/user/downloads/archive.tar.gz -> /home/user/downloads/archive.tar.gz
				(directories)
				/home/user/ ...
					/home/user/my_file.doc -> /home/user/my_file.doc
					/home/user/downloads/archive.tar.gz -> /home/user/downloads/archive.tar.gz

See CONST values in types.go to match the constants with the naming methods.
*/
func (c *Configuration) sanitizeKeyNamingMethod() {
	c.Options[ProfileOptionKeyNamingMethod] = strings.ToLower(strings.TrimSpace(c.Options[ProfileOptionKeyNamingMethod].(string)))
	switch c.Options[ProfileOptionKeyNamingMethod].(string) {
	case NameMethodRelative:
	case NameMethodAbsolute:
	default:
		c.Options[ProfileOptionKeyNamingMethod] = EmptyString
	}
	if c.Options[ProfileOptionKeyNamingMethod] != EmptyString {
		c.Logger.Info(fmt.Sprintf("Key Naming Method set to %s", c.Options[ProfileOptionKeyNamingMethod].(string)))
	}
}

func (c *Configuration) sanitizePrefixes() {
	if c.Options[ProfileOptionPathPrefix] != EmptyString {
		for {
			newKey := strings.ReplaceAll(c.Options[ProfileOptionPathPrefix].(string), "//", "/")
			if c.Options[ProfileOptionPathPrefix].(string) == newKey {
				break
			}
			c.Options[ProfileOptionPathPrefix] = newKey
		}
		c.Options[ProfileOptionPathPrefix] = strings.Trim(c.Options[ProfileOptionPathPrefix].(string), "/")
	}

	if c.Options[ProfileOptionObjectPrefix] != EmptyString {
		for {
			newKey := strings.ReplaceAll(c.Options[ProfileOptionObjectPrefix].(string), "//", "/")
			if c.Options[ProfileOptionObjectPrefix].(string) == newKey {
				break
			}
			c.Options[ProfileOptionObjectPrefix] = newKey
		}
		c.Options[ProfileOptionObjectPrefix] = strings.Trim(c.Options[ProfileOptionObjectPrefix].(string), "/")
	}
}

/*
sanitizeTags() validates the values for the tags field in the profile. If a key or value in a tag contains an
equal sign, then it is removed.
*/
func (c *Configuration) sanitizeTags() {
	for k, v := range c.Tags {
		if strings.Contains(k, "=") || strings.Contains(v, "=") {
			c.Logger.Warn(fmt.Sprintf("Tag key pair %q: %q contains an equal sign, removing.", k, v))
			delete(c.Tags, k)
		}
	}
}

/*
sanitizeStorageType validates the value for the storage field in the profile.

	AWS S3 Storage Classes
	c.Options["storage"] takes an aws storage class:
	- STANDARD: Standard storage class for frequently accessed data.
	- STANDARD_IA: Standard-Infrequent Access storage class for infrequently accessed data.
	- ONEZONE_IA: One Zone-Infrequent Access storage class for infrequently accessed data that can be recreated.
	- INTELLIGENT_TIERING: Intelligent Tiering storage class for data with unknown or changing access patterns.
	- GLACIER: Glacier storage class for long-term data archival.
	- DEEP_ARCHIVE: Deep Archive storage class for long-term data archival with the lowest cost.

	See: https://docs.aws.amazon.com/AmazonS3/latest/userguide/storage-class-intro.html for more information.

See CONST values in types.go to match the constants with the storage classes.
*/
func (c *Configuration) sanitizeStorageType() {
	c.Options[ProfileOptionStorage] = strings.ToUpper(strings.TrimSpace(c.Options[ProfileOptionStorage].(string)))
	switch c.Options[ProfileOptionStorage].(string) {
	case StorageClassStandard:
	case StorageClassStandardIA:
	case StorageClassOneZoneIA:
	case StorageClassIntelligentTiering:
	case StorageClassGlacier:
	case StorageClassDeepArchive:
	case EmptyString:
		c.Logger.Warn("No storage class specified, using default.")
		c.Options[ProfileOptionStorage] = StorageClassStandard
	default:
		c.Logger.Warn("Invalid storage class specified, using default.")
		c.Options[ProfileOptionStorage] = StorageClassStandard
	}
	c.Logger.Info("Storage Class set to " + c.Options[ProfileOptionStorage].(string))
}

/*
repairMissingFields() checks if a field is a different type then expected. If it is, it is cleared and set
to a default value.
*/
func (c *Configuration) repairMissingFields() {
	// Validate the Authentication map: map[string]any
	_, ok := c.Authentication[ProfileAuthProfile]
	if !ok {
		if !isString(c.Authentication[ProfileAuthProfile]) {
			c.Authentication[ProfileAuthProfile] = map[string]any{}
			c.Authentication[ProfileAuthProfile] = EmptyString
		}
	}
	_, ok = c.Authentication[ProfileAuthKey]
	if !ok {
		if !isString(c.Authentication[ProfileAuthKey]) {
			c.Authentication[ProfileAuthKey] = map[string]any{}
			c.Authentication[ProfileAuthKey] = EmptyString
		}
	}
	_, ok = c.Authentication[ProfileAuthSecret]
	if !ok {
		if !isString(c.Authentication[ProfileAuthSecret]) {
			c.Authentication[ProfileAuthSecret] = map[string]any{}
			c.Authentication[ProfileAuthSecret] = EmptyString
		}
	}

	// Validate the Bucket map: map[string]any
	_, ok = c.Bucket[ProfileBucketName]
	if !ok {
		if !isString(c.Bucket[ProfileBucketName]) {
			c.Bucket[ProfileBucketName] = map[string]any{}
			c.Bucket[ProfileBucketName] = EmptyString
		}
	}
	_, ok = c.Bucket[ProfileBucketRegion]
	if !ok {
		if !isString(c.Bucket[ProfileBucketRegion]) {
			c.Bucket[ProfileBucketRegion] = map[string]any{}
			c.Bucket[ProfileBucketRegion] = EmptyString
		}
	}

	// Validate the Options map: map[string]any
	_, ok = c.Options[ProfileOptionACL]
	if !ok {
		if !isString(c.Options[ProfileOptionACL]) {
			c.Options[ProfileOptionACL] = map[string]any{}
			c.Options[ProfileOptionACL] = EmptyString
		}
	}
	_, ok = c.Options[ProfileOptionObjectPrefix]
	if !ok {
		if !isString(c.Options[ProfileOptionObjectPrefix]) {
			c.Options[ProfileOptionObjectPrefix] = map[string]any{}
			c.Options[ProfileOptionObjectPrefix] = EmptyString
		}
	}
	_, ok = c.Options[ProfileOptionOverwrite]
	if !ok {
		if !isBool(c.Options[ProfileOptionOverwrite]) {
			c.Options[ProfileOptionOverwrite] = map[string]any{}
			c.Options[ProfileOptionOverwrite] = false
		}
	}
	_, ok = c.Options[ProfileOptionPathPrefix]
	if !ok {
		if !isString(c.Options[ProfileOptionPathPrefix]) {
			c.Options[ProfileOptionPathPrefix] = map[string]any{}
			c.Options[ProfileOptionPathPrefix] = EmptyString
		}
	}
	_, ok = c.Options[ProfileOptionStorage]
	if !ok {
		if !isString(c.Options[ProfileOptionStorage]) {
			c.Options[ProfileOptionStorage] = map[string]any{}
			c.Options[ProfileOptionStorage] = EmptyString
		}
	}
	_, ok = c.Options[ProfileOptionTagOrigins]
	if !ok {
		if !isBool(c.Options[ProfileOptionTagOrigins]) {
			c.Options[ProfileOptionTagOrigins] = map[string]any{}
			c.Options[ProfileOptionTagOrigins] = true
		}
	}
	_, ok = c.Options[ProfileOptionKeyNamingMethod]
	if !ok {
		if !isString(c.Options[ProfileOptionKeyNamingMethod]) {
			c.Options[ProfileOptionKeyNamingMethod] = map[string]any{}
			c.Options[ProfileOptionKeyNamingMethod] = EmptyString
		}
	}
	_, ok = c.Options[ProfileOptionOmitOriginDir]
	if !ok {
		if !isBool(c.Options[ProfileOptionOmitOriginDir]) {
			c.Options[ProfileOptionOmitOriginDir] = map[string]any{}
			c.Options[ProfileOptionOmitOriginDir] = false
		}
	}

	// Validate the Logging map: map[string]any
	_, ok = c.Logging[ProfileLoggingToConsole]
	if !ok {
		if !isBool(c.Logging[ProfileLoggingToConsole]) {
			c.Logging[ProfileLoggingToConsole] = map[string]any{}
			c.Logging[ProfileLoggingToConsole] = true
		}
	}
	_, ok = c.Logging[ProfileLoggingToFile]
	if !ok {
		if !isBool(c.Logging[ProfileLoggingToFile]) {
			c.Logging[ProfileLoggingToFile] = map[string]any{}
			c.Logging[ProfileLoggingToFile] = false
		}
	}
	_, ok = c.Logging[ProfileLoggingFilename]
	if !ok {
		if !isString(c.Logging[ProfileLoggingFilename]) {
			c.Logging[ProfileLoggingFilename] = map[string]any{}
			c.Logging[ProfileLoggingFilename] = EmptyString
		}
	}
	_, ok = c.Logging[ProfileLoggingLevel]
	if !ok {
		if !isString(c.Logging[ProfileLoggingLevel]) {
			c.Logging[ProfileLoggingLevel] = map[string]any{}
			c.Logging[ProfileLoggingLevel] = 2
		}
	}
}
