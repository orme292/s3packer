# s3packer - A configurable profile-based S3 backup and upload tool.

Linux/MacOS | Supports AWS S3 

---
[![Go Version][go_version_img]][go_version_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Repo License][repo_license_img]][repo_license_url]

[![Code Quality](https://github.com/orme292/s3packer/actions/workflows/golang.yml/badge.svg)](https://github.com/orme292/s3packer/actions/workflows/golang.yml)

## Download 

See the [releases][releases_url] page...

## How to Use

To start a session with an existing profile, just type in the following command:

```bash
$ s3packer --profile="myprofile.yaml"
```

## Creating a new Profile

s3packer can create a base profile to help get you started. To create one, use the `--create` flag:

```bash
$ s3packer --create="my-new-profile.yaml"
```

## Setting up a Profile

s3packer profiles are written in the YAML format. To set it up, you just need to fill out a few fields, and you’ll be good to go!

First, make sure you specify that you're using Version 2 of the profile format:

```yaml
Version: 2
```

Use your AWS Key/Secret pair: 

```yaml
AWS:
  Key: "my-key"
  Secret: "my-secret"
```

Or you can specify a profile that's already set up in your `~/.aws/credentials` file:

```yaml
Version: 2
AWS:
  Profile: "my-profile"
```

Configure your bucket:

```yaml
Bucket:
  Name: "deep-freeze"
  Region: "eu-north-1"
```

And then, tell s3packer what you want to upload. You can specify folders, directories or individual files. (You can call
it the Folders section or the Directories section, it doesn't matter.)

```yaml
Uploads:
  Folders:
    - "/Users/forrest/docs/stocks/apple"
    - "/Users/jenny/docs/song_lyrics"
  Files:
    - "/Users/forrest/docs/job-application-lawn-mower.pdf"
    - "/Users/forrest/docs/dr-pepper-recipe.txt"
    - "/Users/jenny/letters/from-forrest.docx"
```

--- 

### Tags

You can also add tags to your files. Just add a `Tags` section to your profile:

```yaml
Tagging:
  Tags:
    Author: "Forrest Gump"
    Year: 1994
```
---

### Extra Options

You can also customize how your files are stored, accessed, tagged, and uploaded using these options.

---
```yaml
AWS:
  ACL: "private"
  Storage: "ONEZONE_IA"
```
**ACL** <br/>
The default is `private`, but you can use any canned ACL:
- `public-read`
- `public-read-write`
- `authenticated-read`
- `aws-exec-read`
- `bucket-owner-read`
- `bucket-owner-full-control`
- `log-delivery-write`

Read more about ACLs here: [https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html][s3_acl_url]

**Storage** <br/>
The default is `STANDARD`, but you can use any of the following storage classes:
- `STANDARD`
- `STANDARD_IA`
- `ONEZONE_IA`
- `INTELLIGENT_TIERING`
- `GLACIER`
- `DEEP_ARCHIVE`

---

```yaml
Objects:
  NamePrefix: "monthly-"
  RootPrefix: "/backups/monthly"
  Naming: "relative"
```

**NamePrefix** <br/>
This is blank by default. Any value you put here will be added before the filename when it's uploaded to S3.
Using something like `weekly-` will add that string to any file you're uploading, like `weekly-log.log` or `weekly-2021-01-01.log`.

**RootPrefix** <br/>
This is blank by default. Any value put here will be added before the file path when it's uploaded to S3.
If you use something like `/backups/monthly`, the file will be uploaded to `/backups/monthly/your-file.txt`.

**Naming** <br/>
The default is `relative`.
- `relative`: The key will be prepended with the relative path of the file on the local filesystem (individual files specified in the profile will always end up at the root of the bucket, plus the `pathPrefix` and then `objectPrefix`).
- `absolute`: The key will be prepended with the absolute path of the file on the local filesystem.

---

```yaml
Options:
  MaxUploads: 100
  Overwrite: "never"
```

**MaxUploads** <br/>
The default is `5`. This is the maximum number of files that will be uploaded at the same time. Concurrency is at the
directory level, so the biggest speed gains are seen when uploading a directory with many files.

**Overwrite**  <br/>
This is `never` by default. If you set it to `always`, s3packer will Overwrite any files in the bucket that
have the same name as what you're uploading. Useful if you're uploading a file that is updated over and over again.

---

```yaml
Tagging:
  ChecksumSHA256: false
  Origins: true
```
**ChecksumSHA256** <br/>
This is `true` by default. Every object uploaded will be tagged with the file's calculated SHA256 checksum. It'll
be used to verify file changes in the future. Whether this is `true` or `false`, the SHA256 checksum will still be
calculated and used to verify the integrity of the file after it's uploaded.

**Origins** <br/>
This is `true` by default. Every object uploaded will be tagged with the full absolute path of the file on the
local filesystem. This is useful if you want to be able to trace the origin of a file in S3.

---

### Logging Options

And if you like keeping track of things or want a paper trail, you can set up logging too:

```yaml
Logging:
  Level: 1
  Console: true
  File: true
  Filepath: "/var/log/backup.log"
 ```

**Level:**<br/>
This is `2` by default. The setting is by severity, with 0 being least severe and 5 being most severe. 0 will log 
all messages (including debug), and 5 will only log fatal messages which cause the program to exit.

**Console:**<br/>
This is `true` by default. Outputs logging messages to standard output. If you set it to `false`, s3packer
prints minimal output.

**File:**<br/>
This is `false` by default. If you set it to `true`, s3packer will write structured log (JSON) messages to 
a file. You MUST also specify a `filename`.

**Filepath:** <br/>
File to write structured log messages to. If you set `toFile` to `true`, you must specify a filename. 
The file will be created if it doesn't exist, and appended to if it does.

---

### Things to Keep in Mind...

**Individual Files**

If you’re uploading individual files, just remember that the prefix will be added to the start of the filenames and they’ll be uploaded right to the root of the bucket.
Also, if you’ve got multiple files with the same name (like if you have five ‘log.log’ files from different directories), they’ll be sequentially renamed, so log.log-0, log.log-1, log.log-2, etc etc.

**Directories**

When you’re uploading directories, all the subdirectories and files will be uploaded to the bucket as well. Processing
directories with a large number of files can take some time as the checksums are calculated. 

---

### Issues

And if you run into any issues or have any suggestions, feel free to open a new issue on [GitHub][issue_repo_url].


<!-- Links -->
[releases_url]: https://github.com/orme292/s3packer/releases
[issue_repo_url]: https://github.com/orme292/s3packer/issues/new/choose
[go_version_url]: https://golang.org/doc/go1.21
[go_report_url]: https://goreportcard.com/report/github.com/orme292/s3packer
[repo_license_url]: https://github.com/orme292/s3packer/blob/master/LICENSE
[s3_acl_url]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl

[go_version_img]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[repo_license_img]: https://img.shields.io/badge/license-MIT-orange?style=for-the-badge&logo=none
