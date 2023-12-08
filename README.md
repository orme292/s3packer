# s3packer - A profile-based S3 backup and upload tool.

---
    
[![Go Version][go_version_img]][go_version_url]
[![Go Report Card][go_report_img]][go_report_url]
[![Repo License][repo_license_img]][repo_license_url]

## How to Use

Just type in the following command:

```bash
$ s3packer -profile myprofile.yaml
```

## Setting up a Profile

s3packer uses YAML to define profiles. You just need to fill out a few fields and you’re good to go!

Here's what you need:

```yaml
Authentication:
  key: "my-key"
  secret: "my-secret"
```

Or you can use a profile that's already set up in your `~/.aws/credentials` file:

```yaml
Authentication:
  aws-profile: "my-profile"
```

Configure your bucket:

```yaml
Bucket:
  name: "deep-freeze"
  region: "eu-north-1"
```

And then, tell s3packer what you want to upload. You can specify directories or individual files. 

```yaml
Dirs:
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
Tags:
  - "Author": "Forrest Gump"
  - "Year": 1994
```
---

### Extra Options

You can also customize how your files are stored and accessed with these options:

```yaml
Options:
  acl: "private"
  storage: "ONEZONE_IA"
  objectPrefix: "monthly-"
  pathPrefix: "/backups/monthly"
  tagOrigins: true
  keyNamingMethod: "relative""
  overwrite: false
```

**acl:** <br/>
The default is `private`, but you can use any canned ACL: 
- `public-read`
- `public-read-write`
- `authenticated-read`
- `aws-exec-read`
- `bucket-owner-read`
- `bucket-owner-full-control`
- `log-delivery-write`

See more details on the AWS site, [https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html][s3_acl_url]

**storage:** <br/>
The default is `STANDARD`, but you can use any of the following storage classes:
- `STANDARD`
- `STANDARD_IA`
- `ONEZONE_IA`
- `INTELLIGENT_TIERING`
- `GLACIER`
- `DEEP_ARCHIVE`

**pathPrefix** <br/>
This is blank by default. Any value put here will be added before the file path when it's uploaded to S3.
If you use something like `/backups/monthly`, the file will be uploaded to `/backups/monthly/your-file.txt`.

**objectPrefix:** <br/>
This is blank by default. Any value you put here will be added before the filename when it's uploaded to S3.
Using something like `weekly-` will add that string to any file you're uploading, like `weekly-log.log` or `weekly-2021-01-01.log`.

**tagOrigins** <br/>
This is `true` by default. Every object uploaded will be tagged with the full absolute path of the file on the
local filesystem. This is useful if you want to be able to trace the origin of a file in S3.

**keyNamingMethod** <br/>
The default is `relative`.
- `relative`: The key will be prepended with the relative path of the file on the local filesystem.
- `absolute`: The key will be prepended with the absolute path of the file on the local filesystem.


**overwrite:**  <br/>
This is `false` by default. If you set it to `true`, s3packer will overwrite any files in the bucket that 
have the same name as what you're uploading. Useful if you're uploading a file that is updated over and over again.

---

### Logging Options

And if you like keeping track of things or want a paper trail, you can set up logging too:

```yaml
Logging:
  toConsole: true
  toFile: true
  filename: "/var/log/backup.log"
 ```

**toConsole:**<br/>
This is `true` by default. Outputs logging messages to standard output. If you set it to `false`, s3packer
prints minimal output.

**toFile:**<br/>
This is `false` by default. If you set it to `true`, s3packer will write structured log (JSON) messages to 
a file. You MUST also specify a `filename`.

**filename:** <br/>
File to write structured log messages to. If you set `toFile` to `true`, you must specify a filename. 
The file will be created if it doesn't exist, and appended to if it does.

---

### Things to Keep in Mind...

**Individual Files**

If you’re uploading individual files, just remember that the prefix will be added to the start of the filenames and they’ll be uploaded right to the root of the bucket.
Also, if you’ve got multiple files with the same name (like if you have five ‘log.log’ files from different directories), they’ll overwrite each other as they’re uploaded. So, you’ll end up with just one ‘log.log’ file in the bucket.

**Directories**

When you’re uploading directories, all the subdirectories and files will be uploaded to the bucket as well.

---

### Issues

And if you run into any issues or have any suggestions, feel free to open a new issue on [GitHub][issue_repo_url].


<!-- Links -->
[issue_repo_url]: https://github.com/orme292/s3packer/issues/new/choose
[go_version_url]: https://golang.org/doc/go1.21
[go_report_url]: https://goreportcard.com/report/github.com/orme292/s3packer
[repo_license_url]: https://github.com/orme292/s3packer/blob/master/LICENSE
[s3_acl_url]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/acl-overview.html#canned-acl

[go_version_img]: https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go
[go_report_img]: https://img.shields.io/badge/Go_report-A+-success?style=for-the-badge&logo=none
[repo_license_img]: https://img.shields.io/badge/license-MIT-orange?style=for-the-badge&logo=none
