## CHANGELOG
Hello Changelog. This is where I document all changes to s3packer. I'll try to record every bug fix, feature addition, 
or minor tweak. 

---
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

