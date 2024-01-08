## CHANGELOG
This is the Changelog. Between each version, major or minor, I'll document all changes, record every 
bug fix, feature addition, or minor tweak. 

---
### **1.3.0** (2023-01-07)
- s3pack: Removed s3pack
- s3packs: Added s3packs, which has modular support for multiple providers.
- s3packs/objectify: added objectify, that has an object-models for directory trees
- s3packs/objectify: a lot less code than s3pack used to be, but with a ton of for loops, which might not be a good thing...
- s3packs/provider: added provider, which is the start of a modular provider system.
- s3packs/provider: add interface for creating a bucket.
- s3packs/pack_aws: added pack_aws, which is the first provider, AWS S3.
- s3packs/pack_aws: add support for creating a bucket.
- s3packs/pack_aws: added support for multipart uploads with integrity checks.

### **1.2.0** (2023-12-29)
- config: Remove config module
- conf: Add conf module with new AppConfig model
- conf: Profiles are not versioned, only version 2 will be supported
- conf: Adding conf support for the checksum overwrite method and multipart upload, but neither are supported yet
- conf: Add feature to write out a sample profile, `s3packer --create=\"file.yaml\"`
- logbot: Logging now has fmt.Sprintf style formatting
- s3pack: started using the new conf.AppConfig model, removed old config.Configuration model. Much cleaner.
- README updated to reflect new config format and `--create` feature

### **1.1.0** (2023-12-21)
- Upgrade to AWS SDK for Go V2
- Move to Go 1.21.5
- s3pack: Checksum matching on successful upload
- s3pack: Dropped multipart upload support (for now) in favor of checksum matching
- s3pack: AWS SDK for Go V2 dropped the iterator model, so I wrote my own iterator implementation.

### **1.0.3** (2023-12-17)
- s3pack: concurrency for checksum calculations, more speed
- s3pack: concurrency for checking for dupe objects, more speed
- s3pack: counting uploads and ignored files is done on the fly
- s3pack: display total uploaded bytes

### **1.0.2** (2023-12-13)
- config: add new options 'maxConcurrentUploads'
- s3pack: add upload concurrency (handled at ObjectList level)
- s3pack: config references changed to 'c'
- s3pack: FileIterator overhaul, group and index tracking used for concurrency
- s3pack: FileObject has new individual Upload option, but it's unused.
- s3pack: BucketExists checks are done once before processing any files/dirs (See main.go)

### **1.0.1** (2023-12-04)
- use gocritic suggestions
- resolve gosec scan issues
- fix ineffectual assignment
- correct version number

### **1.0.0** (2023-12-03)
- config: More config profile validation occurs.
- config: Added 'level' option to control the logging level (0 debug, 5 Panic)
- config: console and file logging disabled by default
- config: added support for using aws cli profiles instead of secrets/keys
- logbot: fixed an issue where it was impossible to set the logging level
- s3pack: rewrite the whole module
- s3pack: add an explicit bucket check before starting uploads
- s3pack: new structure types: rootlist => dirlist => dirobject => objectlist => fileobject
- s3pack: added keyNamingMethod - relative/absolute
- s3pack: separated prefix options - objectPrefix, pathPrefix
- s3pack: added redundant key check for individual file uploads
- s3pack: added checksum tracking and tagging
- s3pack: add origin tagging
- s3pack: added Tagging support
- s3pack: use HeadObject and HeadBucket to check metadata instead of GetObject
- s3pack: added filesize tracking
- s3pack: removed directory iterator in favor of file iterator
- s3pack: added total upload/ignore tracking and counter
- s3pack: more lines of code, but overall, it's cleaner. Just as slow/fast.
- s3pack: upgrade to AWS SDK 1.48.13

### **0.0.1a** (2023-11-27)

- Adds a README file and fixes a typo in an example profile

### **0.0.1**

- Initial release

