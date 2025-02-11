## CHANGELOG

This is the Changelog. Between each version I'll document changes, bug fixes, feature additions, and minor tweaks.

---

### **1.6.5** (2025-02-11)

- project: repo now uses the internal/pkg project layout
- project: s3packer now referred to as s3p, repo name will change later
- profiles: moved to docs/samples/
- conf: moved to internal/conf
- tuipack: moved to internal/distlog
- s3packs/provider: moved to internal/provider
- s3packs/pack_{aws,gcloud,linode,oracle}: moved to internal/providers/{aws,gcloud,linode,oracle}
- main: removed and merged with cmd/s3p
- cmd: added
- cmd/s3p: added
- docs: refer to s3packer as s3p
- pkgs: added
- internal: added
- internal/conf: added
- internal/distlog: added
- internal/providers/aws: added
- internal/providers/gcloud: added
- internal/providers/linode: added
- internal/providers/oracle: added
- cmd/s3p: now uses cobra to allow for future expansion
- internal/conf: adds methods to the ProviderName type
- internal/conf: exports Filename from Builder
- internal/conf: all references to tui/tuipack change to log/distlog
- internal/providers/{aws,gcloud,linode,oracle}: use internal/conf and internal/distlog
- README: s3packer now referred to as s3p
- README: update docs to reflect the new cli structure
- README: updated links to sample profiles

#### Up Next:

- testing
- s3p will support running from command line without a profile
- download from buckets
- compare bucket to local paths
- simplify the configuration
- support system wide default configuration

---

### **1.6.0** (2025-01-21)

- conf: added gcloud provider and defaults
- conf: added gcloud options, types, and globals
- s3packs: removed types.go (unused)
- s3packs/provider: added support for the Google Cloud provider
- s3packs/providers/gcloud: added new provider with experience support for Google Cloud
- docs: add google cloud config doc
- docs: update README with google cloud information
- profiles: add sample google profile

---
### **1.5.3** (2024-12-27)

- conf: removed the TUI interface and bubbletea
- aws-go-sdk-v2 updated to latest
- oci-go-sdk updated to latest (still v65)

---
### **1.5.0** (2024-07-18)

- conf: restructured and cleaned up the package
- conf: some yaml fields are adjusted
- logbot: removed the package
- tuipack: added new package to handle the TUI and logging output (to replace logbot)
- main: add --noscreen option to force s3p to not use the TUI
- main: support standard logging output or new TUI output
- readme: updated to be a general config document, moved provider specific readmes to the docs/ dir
- docs: removed old readmes and added new docs by provider
- s3packs/objectify: removed in factor of github.com/orme292/objectify and github.com/orme292/symwalker
- s3packs/pack_akamai: removed
- s3packs/pack_aws: removed
- s3packs/pack_oci: removed
- s3packs/provider: rewritten
- s3packs/providers/aws: added new provider for aws based on the new provider package
- s3packs/providers/oci: added new provider for oci *
- s3packs/providers/linode: added new provider for linode *
- aws-go-sdk-v2 updated to latest
- oci-go-sdk updated to v65

---

### **1.4.0** (2024-06-14)

- conf: package rebuilt to be modular and readable.
- conf: Akamai renamed to Linode because Linode is better.
- conf: Directories renamed 'Dirs'
- main: Update --help text
- main: support new conf package
- profiles: update for new conf package
- READMEs: updated with a slightly new format
- s3packs/objectify: support new conf package
- s3packs/pack_akamai: fatal error if bucket cannot be created.
- CHANGELOG: CHANGES LOGGED

---

### **1.3.4** (2024-02-13)
- conf: Added support for the Akamai provider
- conf: Renamed provider-specific files like: provider_aws.go
- conf: Better whitespace trimming from profile fields.
- conf: Fixed bugs in error text for reading the profile file.
- main: Fixed a seg fault that occurred when trying to write to the logger after an error occurred when reading the profile file.
- s3packs/provider: moved bucket exists check to the provider initializer, to reduce the number of times it's called
- s3packs/pack_akamai: core support for Akamai (Linode) Object Storage)
- s3packs/objectify: move types to the type file.
- profiles: added the yaml header "---"
- profiles: added example3.yaml for Akamai
- README: updated with Akamai information, header updated, go-version updated
- README_OCI: header updated
- README_AKAMAI: added
- GITIGNORE: added local CI dev files
- GHA: Updated formats, names, triggers, etc.
- CHANGELOG: CHANGES LOGGED

---

### **1.3.3a** (2024-02-12)
- Use Go 1.22.0
- Update Github Actions to use Go 1.22.0
- Update Dependencies: 
    - aws-sdk-go-v2/feature/s3/manager v1.15.14 -> v1.15.15
    - aws-sdk-go-v2/service/s3 v1.48.0 -> v1.48.1
    - rs/zerolog v1.31.0 -> v1.32.0

### **1.3.3** (2024-02-12)
- conf: Added support for the OCI provider
- conf: Fixed a bug where ChecksumSHA256 was never read from the profile
- s3packs/pack_oci: full support for OCI Object Storage (Oracle Cloud)
- s3packs/pack_oci: workaround OCI SDK's broken metadata handling when using the UploadManager.
- s3packs/pack_aws: fixed broken stats for failed uploads
- s3packs/objectify: fixed broken tagging for ChecksumSHA256 and Origins
- profiles: updated and added example2.yaml
- README: updated with OCI information
- README_OCI: added

### **1.3.2** (2024-01-12)
- s3packs/objectify: removed DirObjList and DirObj. RootList is now a slice of FileObjLists.

### **1.3.1** (2024-01-10)
- replaced old example profiles with a new one that's up to date
- s3packs/objectify: comment update

### **1.3.0** (2024-01-07)
- s3pack: Removed s3pack
- s3packs: Added s3packs, which has modular support for multiple providers.
- s3packs/objectify: added objectify, that has an object-models for directory trees
- s3packs/objectify: a lot less code than s3pack used to be, but with a ton of for loops, which might not be a good thing...
- s3packs/objectify: more robust and resilient file tree builder.
- s3packs/objectify: don't automatically generate checksums, unless the option to tag them is set.
- s3packs/provider: added provider, which is the start of a modular provider system.
- s3packs/provider: add interface for creating a bucket.
- s3packs/provider: stats generation and population done with provider, calculated by objectify.
- s3packs/pack_aws: added pack_aws, which is the first provider, AWS S3.
- s3packs/pack_aws: add support for creating a bucket.
- s3packs/pack_aws: added support for multipart parallel uploads with integrity checks.
- s3packs/pack_aws: lets aws automatically calculate checksums, except for multipart uploads.

### **1.2.0** (2024-12-29)
- config: Remove config module
- conf: Add conf module with new AppConfig model
- conf: Profiles are not versioned, only version 2 will be supported
- conf: Adding conf support for the checksum overwrite method and multipart upload, but neither are supported yet
- conf: Add feature to write out a sample profile, `s3packer --create="file.yaml"`
- logbot: Logging now has fmt.Sprintf style formatting
- s3pack: started using the new conf.AppConfig model, removed old config.Configuration model. Much cleaner.
- README updated to reflect new config format and `--create` feature

### **1.1.0** (2024-12-21)
- Upgrade to AWS SDK for Go V2
- Move to Go 1.21.5
- s3pack: Checksum matching on successful upload
- s3pack: Dropped multipart upload support (for now) in favor of checksum matching
- s3pack: AWS SDK for Go V2 dropped the iterator model, so I wrote my own iterator implementation.

### **1.0.3** (2024-12-17)
- s3pack: concurrency for checksum calculations, more speed
- s3pack: concurrency for checking for dupe objects, more speed
- s3pack: counting uploads and ignored files is done on the fly
- s3pack: display total uploaded bytes

### **1.0.2** (2024-12-13)
- config: add new options 'maxConcurrentUploads'
- s3pack: add upload concurrency (handled at ObjectList level)
- s3pack: config references changed to 'c'
- s3pack: FileIterator overhaul, group and index tracking used for concurrency
- s3pack: FileObject has new individual Upload option, but it's unused.
- s3pack: BucketExists checks are done once before processing any files/dirs (See main.go)

### **1.0.1** (2024-12-04)
- use gocritic suggestions
- resolve gosec scan issues
- fix ineffectual assignment
- correct version number

### **1.0.0** (2024-12-03)
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

