// Package s3pack provides functions for uploading files to s3.
// This file implements functions for creating a session.Session object and an s3manager.Uploader object.
// https://github.com/orme292/s3packer is licensed under the MIT License.
package s3pack

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/orme292/s3packer/config"
)

/*
BuildUploader builds and returned a s3manager.Uploader object. It takes a config.Configuration object. The func creates
a session by calling NewSession and passes it to s3manager.NewUploader.
*/
func BuildUploader(c *config.Configuration) (uploader *s3manager.Uploader, err error) {
	sess, err := NewSession(c)
	uploader = s3manager.NewUploader(sess)
	return
}

/*
NewSession creates a new session.Session object using the authentication information in the profile. It takes a
config.Configuration object. It determines whether to use a profile or a keypair based on the presence of a profile
name in the profile configuration.

NewSession calls NewSessionWithProfile or NewSessionWithKeypair
*/
func NewSession(c *config.Configuration) (sess *session.Session, err error) {
	if c.Authentication[config.ProfileAuthProfile].(string) != EmptyString {
		sess, err = NewSessionWithProfile(c)
	} else {
		sess, err = NewSessionWithKeypair(c)
	}
	if err != nil {
		c.Logger.Fatal(fmt.Sprintf("Unable to create session: %q", err.Error()))
	}
	return
}

/*
NewSessionWithKeypair creates a new session.Session object using a key/secret pair. It takes a config.Configuration
object as input and returns a session.Session object and an error, if there is one.
*/
func NewSessionWithKeypair(c *config.Configuration) (sess *session.Session, err error) {
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(c.Bucket[config.ProfileBucketRegion].(string)),
		Credentials: credentials.NewStaticCredentials(
			c.Authentication[config.ProfileAuthKey].(string), c.Authentication[config.ProfileAuthSecret].(string), ""),
	})
	if err != nil {
		return
	}
	return
}

/*
NewSessionWithProfile creates a new session.Session object using a profile. It takes a config.Configuration object
as input and returns a session.Session object and an error, if there is one.
*/
func NewSessionWithProfile(c *config.Configuration) (sess *session.Session, err error) {
	sess, err = session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{},
		Profile:           c.Authentication[config.ProfileAuthProfile].(string),
		SharedConfigState: session.SharedConfigEnable,
	})
	return
}
