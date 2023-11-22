package config

import (
	"errors"
	"strings"
)

/*
isString checks if the given interface{} is a string
*/
func isString(s any) bool {
	switch s.(type) {
	case string:
		return true
	}
	return false
}

/*
isBool checks if the given interface{} is a bool
*/
func isBool(b any) bool {
	switch b.(type) {
	case bool:
		return true
	}
	return false
}

func (c *Configuration) Validate() error {
	c.createMissingMaps()
	c.repairMissingFields()
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
	if c.Authentication == nil {
		c.Authentication = map[string]any{}
	}
	if c.Bucket == nil {
		c.Bucket = map[string]any{}
	}
	if c.Options == nil {
		c.Options = map[string]any{}
	}
	if c.Logging == nil {
		c.Logging = map[string]any{}
	}
}

/*
criticalMissingValues() checks if required fields are missing from the profile.
You must run createMissingMaps() before using this function.
*/
func (c *Configuration) criticalMissingValues() error {
	if c.Authentication["key"].(string) == "" {
		return errors.New("authentication key is empty")
	}
	if c.Authentication["secret"].(string) == "" {
		return errors.New("authentication secret token is empty")
	}
	if c.Bucket["name"].(string) == "" {
		return errors.New("bucket name is empty")
	}
	if c.Bucket["region"].(string) == "" {
		return errors.New("bucket region is empty")
	}
	if c.Logging["toFile"].(bool) == true && c.Logging["path"].(string) == "" {
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
*/
func (c *Configuration) sanitizeACL() {

	c.Options["acl"] = strings.ToLower(strings.TrimSpace(c.Options["acl"].(string)))
	switch c.Options["acl"].(string) {
	case "private":
	case "public-read":
	case "public-read-write":
	case "authenticated-read":
	case "aws-exec-read":
	case "bucket-owner-read":
	case "bucket-owner-full-control":
	case "log-delivery-write":
	case "":
		c.Logger.Warn("No ACL specified, using default.")
		c.Options["acl"] = "private"
	default:
		c.Logger.Warn("Invalid ACL specified, using default.")
		c.Options["acl"] = "private"
	}
	c.Logger.Info("ACL set to " + c.Options["acl"].(string))
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
*/
func (c *Configuration) sanitizeStorageType() {
	c.Options["storage"] = strings.ToUpper(strings.TrimSpace(c.Options["storage"].(string)))
	switch c.Options["storage"].(string) {
	case "STANDARD":
	case "STANDARD_IA":
	case "ONEZONE_IA":
	case "INTELLIGENT_TIERING":
	case "GLACIER":
	case "DEEP_ARCHIVE":
	case "":
		c.Logger.Warn("No storage class specified, using default.")
		c.Options["storage"] = strings.ToUpper("STANDARD")
	default:
		c.Logger.Warn("Invalid storage class specified, using default.")
		c.Options["storage"] = strings.ToUpper("STANDARD")
	}
	c.Logger.Info("Storage Class set to " + c.Options["storage"].(string))
}

/*
repairMissingFields() checks if a field is a different type then expected. If it is, it is cleared and set
to a default value.
*/
func (c *Configuration) repairMissingFields() {
	// Validate the Authentication map: map[string]any
	_, ok := c.Authentication["key"]
	if !ok {
		if !isString(c.Authentication["key"]) {
			c.Authentication["key"] = map[string]any{}
			c.Authentication["key"] = ""
		}
	}
	_, ok = c.Authentication["secret"]
	if !ok {
		if !isString(c.Authentication["secret"]) {
			c.Authentication["secret"] = map[string]any{}
			c.Authentication["secret"] = ""
		}
	}

	// Validate the Bucket map: map[string]any
	_, ok = c.Bucket["name"]
	if !ok {
		if !isString(c.Bucket["name"]) {
			c.Bucket["name"] = map[string]any{}
			c.Bucket["name"] = ""
		}
	}
	_, ok = c.Bucket["region"]
	if !ok {
		if !isString(c.Bucket["region"]) {
			c.Bucket["region"] = map[string]any{}
			c.Bucket["region"] = ""
		}
	}

	// Validate the Options map: map[string]any
	_, ok = c.Options["acl"]
	if !ok {
		if !isString(c.Options["acl"]) {
			c.Options["acl"] = map[string]any{}
			c.Options["acl"] = ""
		}
	}
	_, ok = c.Options["storage"]
	if !ok {
		if !isString(c.Options["storage"]) {
			c.Options["storage"] = map[string]any{}
			c.Options["storage"] = ""
		}
	}
	_, ok = c.Options["prefix"]
	if !ok {
		if !isString(c.Options["prefix"]) {
			c.Options["prefix"] = map[string]any{}
			c.Options["prefix"] = ""
		}
	}
	_, ok = c.Options["overwrite"]
	if !ok {
		if !isBool(c.Options["overwrite"]) {
			c.Options["overwrite"] = map[string]any{}
			c.Options["overwrite"] = false
		}
	}

	// Validate the Logging map: map[string]any
	_, ok = c.Logging["toConsole"]
	if !ok {
		if !isBool(c.Logging["toConsole"]) {
			c.Logging["toConsole"] = map[string]any{}
			c.Logging["toConsole"] = true
		}
	}
	_, ok = c.Logging["toFile"]
	if !ok {
		if !isBool(c.Logging["toFile"]) {
			c.Logging["toFile"] = map[string]any{}
			c.Logging["toFile"] = false
		}
	}
	_, ok = c.Logging["filename"]
	if !ok {
		if !isString(c.Logging["filename"]) {
			c.Logging["filename"] = map[string]any{}
			c.Logging["filename"] = ""
		}
	}
	_, ok = c.Logging["level"]
	if !ok {
		if !isString(c.Logging["level"]) {
			c.Logging["level"] = map[string]any{}
			c.Logging["level"] = 2
		}
	}
}
