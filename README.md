# s3p - A configurable profile-based object backup and upload tool.

**CLI for Linux/MacOS**  **supports** Amazon S3 **|** Google Cloud Storage **|** Linode (Akamai) Object Storage **|**
Oracle Cloud Object Storage

---
[![Go Version][go_version_img]][go_version_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Repo License][repo_license_img]][repo_license_url]

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/orme292/s3packer/tests.yml?style=for-the-badge&label=Tests&labelColor=blue)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/orme292/s3packer/quality.yml?style=for-the-badge&label=Vulnerabilities&labelColor=blue)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/orme292/s3packer/compile.yml?style=for-the-badge&label=Compiles&labelColor=blue)





---

[![Jetbrains_OSS][jetbrains_logo]][jetbrains_oss_url] [![Jetbrains_GoLand][jetbrains_goland_logo]][jetbrains_goland_url]

Special thanks to JetBrains! **s3p** was developed with help from JetBrains' Open Source program.

---
## About

s3p is an S3 / Object Storage upload and backup tool. It uses YAML-based configs that tell it what to upload, where to
upload, how to name, and how to tag the objects. s3p makes backup redundancy easier by using separate profiles
for buckets and providers. Currently, it supports AWS, Google Cloud, Linode (Akamai), OCI (Oracle Cloud).

---

## Download

See the [releases][releases_url] page...

---

## Service Provider Support

**s3p** supports AWS S3, Google Cloud Storage, Oracle Cloud Object Storage (OCI), and Linode (Akamai) Object
Storage.

- AWS: [using_aws.md][s3packer_aws_readme_url]
- Google: [using_gcloud.md][s3packer_gcloud_readme_url] **(experimental)**
- Linode: [using_linode.md][s3packer_akamai_readme_url]
- OCI: [using_oci.md][s3packer_oci_readme_url]

See the example profiles:

- [example_aws.yaml][example_aws_url] (AWS)
- [example_gcloud.yaml][example_gcloud_url] (Google Cloud)
- [example_linode.yaml][example_linode_url] (Linode)
- [example_oci.yaml][example_oci_url] (OCI)

---

## How to Use

To start a session with an existing profile, just type in the following command:

```bash
~$ s3p use -f "my-custom-profile.yml"
```

---

## Creating a new Profile

s3p can create a base profile to help get you started. To create one, use the `--create` flag:

```bash
~$ s3p profile sample --filename "new-profile.yml"
```

---

## Setting up a Profile

s3p profiles are written in YAML. To set one up, you just need to fill out a few fields before you can get started.

### **Provider**

Tell s3p which service you're using

| PROVIDER | Acceptable Values        | Required | Description                                |
|:---------|:-------------------------|:---------|:-------------------------------------------|
| Use      | aws, google, linode, oci | Y        | name of service provider you will be using |

```yaml
Provider:
  Use: aws
```

---
**Each provider needs their own special fields to configure.**<br/>
SEE: [docs/general_config.md](https://github.com/orme292/s3packer/blob/master/docs/)

### **Bucket**

Tell s3p where the bucket is and whether to create it

| BUCKET | Acceptable Values | Default | Required | Description                                                  |
|:-------|:------------------|:--------|:---------|:-------------------------------------------------------------|
| Create | boolean           | false   | F        | Whether s3p should create the bucket if it is not found      |
| Name   | any string        |         | Y        | The name of the bucket                                       |
| Region | any string        |         | Y        | The region that the bucket is or will be in, e.g. eu-north-1 |

```yaml
Bucket:
  Create: true
  Name: "deep-freeze"
  Region: "eu-north-1"
```

---

### **Options**

s3p's configurable options

| OPTIONS          | Acceptable Values | Default | Required | Description                                                |
|:-----------------|:------------------|:--------|:---------|:-----------------------------------------------------------|
| MaxUploads       | any integer       | 1       | N        | The number of simultaneous uploads, at least 1.            |
| FollowSymlinks   | boolean           | false   | N        | Whether to follow symlinks under dirs provided             |
| WalkDirs         | boolean           | true    | N        | Whether s3p will walk subdirectories of dirs provided      |
| OverwriteObjects | always, never     | never   | N        | Whether overwrite objects that already exist in the bucket |

```yaml
Options:
  MaxUploads: 1
  FollowSymlinks: false
  WalkDirs: true
  OverwriteObjects: "never"
```

### **MaxUploads (Experimental)**

**The suggested setting for `MaxUploads` is `1`.**

Be careful when you set `MaxUploads` as some services struggle with anything more than `1`. The notable exception being
AWS which seems fine with 50-100 -- though, whether its faster to have a high `MaxUploads` value depends on your upload
job.

Large files can be broken up into many parts which are then simultaneously uploaded. s3p uses default SDK values
for part count, part size, and the large file threshold values, unless otherwise called out.

An example of this would be: If you specify a MaxUploads value of 5, and s3p tries to upload 5 large files that are
each split into 20 parts, then there would be 100 simultaneous uploads happening. If you specify a MaxUpload value of 50
and there are 50 large files each split into 20 parts, then you could potentially have as many as 1,000 simultaneous
uploads.

---

### **Objects**

s3p's configurable options for object name and renaming

| OBJECTS     | Acceptable Values  | Default | Required | Description                                                                              |
|:------------|:-------------------|:--------|:---------|:-----------------------------------------------------------------------------------------|
| NamingType  | absolute, relative |         | Y        | the method s3p uses to name objects that it uploads                                      |
| NamePrefix  | any string         |         | N        | The string that will be prefixed to the object's "file" name                             |
| PathPrefix  | any string         |         | N        | a string path that will be prefixed to the object's "file" name and "path" name          |
| OmitRootDir | boolean            | True    | N        | whether the relative root of a provided directory will be added to the objects path name |

```yaml
Objects:
  NamingType: absolute
  NamePrefix: backup-
  PathPrefix: /backups/april/2023
  OmitRootDir: true
```

**NamingType** <br/>
The default is `relative`.

- `relative`: The key will be prepended with the relative path of the file on the local filesystem (individual files
  specified in the profile will always end up at the root of the bucket, plus the `PathPrefix` and then `NamePrefix`).
- `absolute`: The key will be prepended with the absolute path of the file on the local filesystem.

**NamePrefix** <br/>
This is blank by default. Any value you put here will be added before the filename when it's uploaded to S3. Using
something like `weekly-` will add that string to any file you're uploading, like `weekly-log.log`
or `weekly-2021-01-01.log`.

**PathPrefix** <br/>
This is blank by default. Any value put here will be added before the file path when it's uploaded to S3. If you use
something like `/backups/monthly`, the file will be uploaded to `/backups/monthly/your-file.txt`.

---

### **Files, Dirs**

Tells s3p what you want to upload. You can specify directories or individual files. When you specify a directory, s3p
will **NOT** traverse subdirectories, unless configured to. You must specify one or the other.

| FILES | Required | Description                                         |
|:------|:---------|:----------------------------------------------------|
| path  | Y        | the absolute path to the file that will be uploaded |

| DIRS | Required | Description                                              |
|:-----|:---------|:---------------------------------------------------------|
| path | Y        | the absolute path to the directory that will be uploaded |

```yaml
Files:
  - "/Users/forrest/docs/stocks/apple"
  - "/Users/jenny/docs/song_lyrics"
Dirs:
  - "/Users/forrest/docs/objJob-application-lawn-mower.pdf"
  - "/Users/forrest/docs/dr-pepper-recipe.txt"
  - "/Users/jenny/letters/from-forrest.docx"
```

--- 

### **Tags**

Add tags to each uploaded object (if supported by the provider)

| TAGS | Acceptable Values | Required | Description                                       |
|:-----|:------------------|:---------|:--------------------------------------------------|
| Key  | any value         | N        | key:value tag pair, will be converted to a string |


```yaml
Tags:
  Author: "Forrest Gump"
  Year: 1994
```

---

### **TagOptions**

Options related to object tagging (dependent on whether the provider supports object tagging)

| TAGOPTIONS     | Acceptable Values | Default | Required | Description                                                                      |
|:---------------|:------------------|:--------|:---------|:---------------------------------------------------------------------------------|
| OriginPath     | boolean           | False   | N        | Whether s3p will tag the object with the original absolute path of the file      |
| ChecksumSHA256 | boolean           | False   | N        | Whether s3p will tag the object with the sha256 checksum of the file as uploaded |

```yaml
Tagging:
  OriginPath: true
  ChecksumSHA256: false
```

**Note on Checksum Tagging**<br/>
Some providers have checksum validation on objects to verify that uploads are completed correctly. This checksum is
calculated separately from that process and is only for your reference.

---

### Logging

Options for logging output

| LOGGING | Acceptable Values | Default            | Required | Description                                                                            |
|:--------|:------------------|:-------------------|:---------|:---------------------------------------------------------------------------------------|
| Level   | 1-5               | 2                  | N        | The severity level a log message must be to output to the console or file              |
| Console | boolean           | True               | N        | Whether logging message will be output to stdout.                                      |
| File    | boolean           | False              | N        | Whether logging output will be written to a file. Output is structured in JSON format. |
| Logfile | path              | "/var/log/s3p.log" | N        | The name of the file that output logging will be appended to.                          |


```yaml
Logging:
  Level: 3
  Console: true
  File: true
  Logfile: "/var/log/backup.log"
 ```

**Notes on Level**<br/>
Level is set to `2` (WARN) by default. The setting is by severity, with 1 being least severe (INFO) and 5 being most
severe (PANIC).

---

### Things to Keep in Mind...

**Individual Files**

If you’re uploading individual files, just remember that the prefix will be added to the start of each filename and
they’ll be uploaded right to the root of the bucket (unless you specify custom _PathPrefix_. Note that if you have
multiple files with the same name (like if you have five ‘log.log’ files from different directories), they could be
overwritten as you upload.

**Directories**

When you’re uploading directories, when WalkDirs is set to true, then all the subdirectories and files will be uploaded
to the bucket as well. Processing directories with a large number of files can take some time as the checksums are
calculated and each directory entry is read.

---

### Issues & Suggestions

If you run into problems, errors, or have feature suggestions, it would be great if you took the time to open a new
issue on [GitHub][issue_repo_url].

---


<!-- Links -->
[releases_url]: https://github.com/orme292/s3packer/releases
[issue_repo_url]: https://github.com/orme292/s3packer/issues/new/choose

[go_version_url]: https://golang.org/doc/go1.23
[go_report_url]: https://goreportcard.com/report/github.com/orme292/s3packer
[repo_license_url]: https://github.com/orme292/s3packer/blob/master/LICENSE
[go_tests_url]: https://img.shields.io/github/actions/workflow/status/orme292/s3packer/tests.yml?style=for-the-badge&label=Tests

[s3_acl_url]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl

[s3packer_aws_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_aws.md

[s3packer_gcloud_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_gcloud.md
[s3packer_oci_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_oci.md
[s3packer_akamai_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_linode.md

[example_aws_url]:https://github.com/orme292/s3packer/blob/master/docs/sample/example_aws.yaml

[example_gcloud_url]:https://github.com/orme292/s3packer/blob/master/docs/sample/example_gcloud.yaml

[example_linode_url]:https://github.com/orme292/s3packer/blob/master/docs/sample/example_linode.yaml

[example_oci_url]:https://github.com/orme292/s3packer/blob/master/docs/sample/example_oci.yaml


[go_version_img]: https://img.shields.io/github/go-mod/go-version/orme292/s3packer?style=for-the-badge&logo=go
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[repo_license_img]: https://img.shields.io/badge/license-MIT-orange?style=for-the-badge&logo=none

[jetbrains_logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/jb_square.svg
[jetbrains_oss_url]: https://www.jetbrains.com/community/opensource/#support
[jetbrains_goland_logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.svg
[jetbrains_goland_url]: https://www.jetbrains.com/goland