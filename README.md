# s3packer - A configurable profile-based S3 backup and upload tool.

**CLI for Linux/MacOS**  **supports** Amazon S3 **|** Oracle Cloud Object Storage **|** Linode (Akamai) Object Storage

---
[![Go Version][go_version_img]][go_version_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Repo License][repo_license_img]][repo_license_url]

[![Code Quality](https://github.com/orme292/s3packer/actions/workflows/golang.yml/badge.svg)](https://github.com/orme292/s3packer/actions/workflows/golang.yml)

---

[![Jetbrains_OSS][jetbrains_logo]][jetbrains_oss_url] [![Jetbrains_GoLand][jetbrains_goland_logo]][jetbrains_goland_url]

Special thanks to JetBrains! </br>
**s3packer** is developed with support from the JetBrains Open Source program.

---
## About

s3packer is a configurable yaml-based S3 storage upload and backup tool. Use YAML-based configs with s3packer that tell
it what to upload, where to upload, how to name, and how to tag the objects. Redundancy is easier by using separate
profiles for each provider. s3packer supports several AWS, OCI (Oracle Cloud), and Linode (Akamai).

---

## Add Support for new Providers

**You can build support for a custom provider by using the _s3packs/provider_ package interfaces. To implement your own
provider:**

- Create a new package under s3packs/providers (e.g. s3packs/providers/azure). Use a simple name prefixed with s3,
  like "s3azure", as the package name.
- Implement the operator and object interface (
  see: [s3packs/provider/interfaces.go](https://github.com/orme292/s3packer/blob/master/s3packs/provider/interfaces.go))
- Implement the generator function interfaces (OperGenFunc, ObjectGenFunc)
- Add the required configuration code
  - see: [conf/type_provider.go](https://github.com/orme292/s3packer/blob/master/conf/type_provider.go)
  - see example: [conf/provider_aws.go](https://github.com/orme292/s3packer/blob/master/conf/provider_aws.go)
- Add new provider to the getProviderFunctions
  fn [s3packs/main.go](https://github.com/orme292/s3packer/blob/master/s3packs/main.go)

See current provider code for implementation
examples: [aws](https://github.com/orme292/s3packer/blob/master/s3packs/providers/aws), [oci](https://github.com/orme292/s3packer/blob/master/s3packs/providers/oci), [linode](https://github.com/orme292/s3packer/blob/master/s3packs/providers/linode)

---

## Download

See the [releases][releases_url] page...

---
## Providers

**s3packer** supports AWS S3, Oracle Cloud Object Storage (OCI), and Linode (Akamai) Object Storage.

- AWS: [using_aws.md][s3packer_aws_readme_url]
- OCI: [using_oci.md][s3packer_oci_readme_url]
- Linode: [using_linode.md][s3packer_akamai_readme_url]

See the example profiles:
- [example1.yaml][example1_url] (AWS)
- [example2.yaml][example2_url] (OCI)
- [example3.yaml][example3_url] (Linode/Akamai)
---

## How to Use

To start a session with an existing profile, just type in the following command:

```bash
$ s3packer --profile="myprofile.yaml"
```

---

## Creating a new Profile

s3packer can create a base profile to help get you started. To create one, use the `--create` flag:

```bash
$ s3packer --create="my-new-profile.yaml"
```

---

## Setting up a Profile

s3packer profiles are written in YAML. To set one up, you just need to fill out a few fields before you can get started.

### **Version**<br/>

```yaml
Version: 6
```

---

### **Provider**

Tell s3packer which service you're using

| PROVIDER | Acceptable Values | Required | Description                        |
|:---------|:------------------|:---------|:-----------------------------------|
| Use      | aws, oci, linode  | Y        | name of provider you will be using |

```yaml
Provider:
  Use: aws
```

---
Each provider needs their own special fields filled out.<br/>
SEE: [docs/general_config.md](https://github.com/orme292/s3packer/blob/master/docs/)

### **Bucket**

Tell s3packer where the bucket is and whether to create it

| BUCKET | Acceptable Values | Default | Required | Description                                                  |
|:-------|:------------------|:--------|:---------|:-------------------------------------------------------------|
| Create | boolean           | false   | F        | Whether s3packer should create the bucket if it is not found |
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

s3packer's configurable options

| OPTIONS          | Acceptable Values | Default | Required | Description                                                |
|:-----------------|:------------------|:--------|:---------|:-----------------------------------------------------------|
| MaxUploads       | any integer       | 1       | N        | The number of simultaneous uploads, at least 1.            |
| FollowSymlinks   | boolean           | false   | N        | Whether to follow symlinks under dirs provided             |
| WalkDirs         | boolean           | true    | N        | Whether s3packer will walk subdirectories of dirs provided |
| OverwriteObjects | always, never     | never   | N        | Whether overwrite objects that already exist in the bucket |

<details>
<summary>
  MaxUpload considerations
</summary>

Some providers can struggle with a high number of simultaneous uploads. Generally, anywhere between 1 and 5 is safe,
however providers like AWS have demonstrated the ability to handle up to 50, or even more.

It's important to note that large files can be broken up into many parts which are then simultaneously uploaded. Part
count, part size, and the large file threshold values are not configured by s3packer, unless otherwise called out.

For example, if you specify a MaxUploads value of 5, and s3packer tries to upload 5 large files that are each split into
20 parts, then there would be 100 simultaneous uploads happening. If you specify a MaxUpload value of 50 and there are
50 large files each split into 20 parts, then you could potentially have as many as 1,000 simultaneous uploads.
</details>

```yaml
Options:
  MaxUploads: 1
  FollowSymlinks: false
  WalkDirs: true
  OverwriteObjects: "never"
```

---

### **Objects**

s3packer's configurable options for object name and renaming

| OBJECTS     | Acceptable Values  | Default | Required | Description                                                                              |
|:------------|:-------------------|:--------|:---------|:-----------------------------------------------------------------------------------------|
| NamingType  | absolute, relative |         | Y        | the method s3packer uses to name objects that it uploads                                 |
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
  specified in the profile will always end up at the root of the bucket, plus the `pathPrefix` and then `objectPrefix`).
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

Tells s3packer what you want to upload. You can specify directories or individual files. When you specify a directory,
s3packer will **NOT** traverse subdirectories, unless configured to. You must specify one or the other.

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

| TAGOPTIONS     | Acceptable Values | Default | Required | Description                                                                           |
|:---------------|:------------------|:--------|:---------|:--------------------------------------------------------------------------------------|
| OriginPath     | boolean           | False   | N        | Whether s3packer will tag the object with the original absolute path of the file      |
| ChecksumSHA256 | boolean           | False   | N        | Whether s3packer will tag the object with the sha256 checksum of the file as uploaded |

```yaml
Tagging:
  OriginPath: true
  ChecksumSHA256: false
```

**Note on Checksum Tagging**<br/>
Some providers have checksum validation on objects to verify that uploads are completed correctly. This checksum is
calculated separately from that process and is only for your future reference.

---

### Logging

Options for logging output

| LOGGING | Acceptable Values | Default            | Required | Description                                                                               |
|:--------|:------------------|:-------------------|:---------|:------------------------------------------------------------------------------------------|
| Screen  | boolean           | False              | N        | Whether s3packer will run in pretty mode, with an AltScreen display                       |
| Level   | 1-5               | 2                  | N        | The severity level a log message must be to output to the console or file                 |
| Console | boolean           | True               | N        | Whether logging message will be output to stdout (set to false if screen is set to true). |
| File    | boolean           | False              | N        | Whether logging output will be written to a file. Output is structured in JSON format.    |
| Logfile | path              | "/var/log/s3p.log" | N        | The name of the file that output logging will be appended to.                             |


```yaml
Logging:
  Screen: false
  Level: 3
  Console: true
  File: true
  Logfile: "/var/log/backup.log"
 ```

**Notes on Level**<br/>
This is `2` WARN by default. The setting is by severity, with 1 being least severe (INFO) and 5 being most severe (
PANIC).

---

### Things to Keep in Mind...

**Individual Files**

If you’re uploading individual files, just remember that the prefix will be added to the start of the filenames and they’ll be uploaded right to the root of the bucket.
Note that if you have multiple files with the same name (like if you have five ‘log.log’ files from different
directories), they could be overwritten as you upload.

**Directories**

When you’re uploading directories, all the subdirectories and files will be uploaded to the bucket as well. Processing
directories with a large number of files can take some time as the checksums are calculated.

---

### Issues & Suggestions

If you run into any problems, errors, or have feature suggestions PLEASE feel free to open a new issue on
[GitHub][issue_repo_url].

---


<!-- Links -->
[releases_url]: https://github.com/orme292/s3packer/releases
[issue_repo_url]: https://github.com/orme292/s3packer/issues/new/choose
[go_version_url]: https://golang.org/doc/go1.22
[go_report_url]: https://goreportcard.com/report/github.com/orme292/s3packer
[repo_license_url]: https://github.com/orme292/s3packer/blob/master/LICENSE
[s3_acl_url]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl

[s3packer_aws_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_aws.md

[s3packer_oci_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_oci.md

[s3packer_akamai_readme_url]: https://github.com/orme292/s3packer/blob/master/docs/using_linode.md

[example1_url]:https://github.com/orme292/s3packer/blob/master/profiles/example1.yaml
[example2_url]:https://github.com/orme292/s3packer/blob/master/profiles/example2.yaml
[example3_url]:https://github.com/orme292/s3packer/blob/master/profiles/example3.yaml

[go_version_img]: https://img.shields.io/badge/Go-1.22-00ADD8?style=for-the-badge&logo=go
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[repo_license_img]: https://img.shields.io/badge/license-GPL%202.0-orange?style=for-the-badge&logo=none

[jetbrains_logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/jb_square.svg
[jetbrains_oss_url]: https://www.jetbrains.com/community/opensource/#support
[jetbrains_goland_logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.svg
[jetbrains_goland_url]: https://www.jetbrains.com/goland