# s3packer

**Linux/MacOS |||||** Amazon S3 **|** Oracle Cloud Object Storage **|** Akamai (Linode) Object Storage

---
## Using s3packer with Akamai (Linode) Object Storage

◀️ [Back to s3packer][s3packer_readme_url]

---

## Creating a new Profile

s3packer can create a base profile to help get you started. To create one, use the `--create` flag:

```bash
$ s3packer --create="my-new-akamai-profile.yaml"
```

## Setting up a Profile

s3packer profiles are written in YAML. To set one up, you just need to fill out a few fields, and you’ll be good to go!

First, make sure you specify that you're using Version 4 of the profile format:

```yaml
Version: 5
```

Be sure to specify a provider:

```yaml
Version: 5
Provider:
  Use: Linode
```

---
## Authentication

**s3packer** uses object storage access keys to authenticate with Linode. You can find and generate access keys in the
Cloud Manager. For info on generating new access keys, check out the [Linode Object Storage Guide][akamai_auth_url].

```yaml
Version: 5
Provider:
  Use: Linode
  Key: "zzzyyyyxxxxx1111222"
  Secret: "aabbbcccddddeeeffff999988888"
```

Configure the `region` to generate a Linode object storage endpoint name.

```yaml
Linode:
  Region: se-sto-1
```

Next, configure the bucket. The `name` and `region` fields are required. If the `region` field isn't correct,
s3packer won't find the bucket and (if configured to) will create a new one in the specified region.

`Create` defaults to `false`. If `true`, s3packer will create the bucket in the specified region if it doesn't exist.

`Region` should contain the region short-code. When you create a bucket in the Cloud Manager, the short code will be
listed in the region dropdown. You can also check Linode's documentation for a list of region short-codes:
[Linode Region List][akamai_region_list_url].

```yaml
Bucket:
  Create: false
  Name: "s3packer-bucket"
  Region: "se-sto-1"
```

And then, tell s3packer what you want to upload. You can specify directories or individual files. When you specify a
directory, s3packer will traverse all subdirectories.

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

### Tags

Unfortunately, tags are not supported by Linode Object Storage.

---

### Linode Specific Options

Unfortunately, s3packer does not support assigning ACLs to Linode objects. Linode does not support storage tiers.

---

### Extra Options

You can customize how your files are named and uploaded to Linode object storage.

---

```yaml
Objects:
  NamingType: "relative"
  NamePrefix: "monthly-"
  PathPrefix: "/backups/monthly"
```

**NamingType** <br/>
The default is `relative`.

- `relative`: The key will be prepended with the relative path of the file on the local filesystem (individual files
  specified in the profile will always end up at the root of the bucket, plus the `pathPrefix` and then `objectPrefix`).
- `absolute`: The key will be prepended with the absolute path of the file on the local filesystem.

**NamePrefix** <br/>
This is blank by default. Any value you put here will be added before the filename when it's uploaded to S3.
Using something like `weekly-` will add that string to any file you're uploading, like `weekly-log.log` or `weekly-2021-01-01.log`.

**PathPrefix** <br/>
This is blank by default. Any value put here will be added before the file path when it's uploaded to S3.
If you use something like `/backups/monthly`, the file will be uploaded to `/backups/monthly/your-file.txt`.

---

```yaml
Options:
  OverwriteObjects: "never"
```

**MaxParts** <br/>
The default depends on the provider. The AWS default is `100`. MaxParts specifies the number of pieces a large file will
be broken up into before uploading and reassembling.

**MaxUploads** <br/>
The default is `5`. This is the maximum number of files that will be uploaded at the same time. Concurrency is at the
directory level, so the biggest speed gains are seen when uploading a directory with many files.

**OverwriteObjects**  <br/>
This is `never` by default. If you set it to `always`, s3packer will overwrite any files in the bucket that
have the same name as what you're uploading. Useful if you're uploading a file that is updated over and over again.

---

### Logging Options

And if you like keeping track of things or want a paper trail, you can set up logging too:

```yaml
Logging:
  Level: 1
  OutputToConsole: true
  OutputToFile: true
  Path: "/var/log/backup.log"
 ```

**Level:**<br/>
This is `2` by default. The setting is by severity, with 0 being least severe and 5 being most severe. 0 will log
all messages (including debug), and 5 will only log fatal messages which cause the program to exit.

**OutputToConsole:**<br/>
This is `true` by default. Outputs logging messages to standard output. If you set it to `false`, s3packer
prints minimal output.

**OutputToFile:**<br/>
This is `false` by default. If you set it to `true`, s3packer will write structured log (JSON) messages to a file. You
MUST also specify a `Path`.

**Path:** <br/>
Path of the file to write structured log messages to. If you set `OutputToFile` to `true`, you must specify a filename.
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

### Issues & Suggestions

If you run into any problems, errors, or have feature suggestions PLEASE feel free to open a new issue on
[GitHub][issue_repo_url].

---

### Development

s3packer was built with Jetbrains GoLand

[![Jetbrains_OSS][GoLand logo]][jetbrains_goland_url]

---

◀️ [Back to s3packer][s3packer_readme_url]

<!-- Links -->
[s3packer_readme_url]: https://github.com/orme292/s3packer/blob/master/README.md
[akamai_auth_url]: https://www.linode.com/docs/products/storage/object-storage/guides/access-keys/
[akamai_region_list_url]: https://www.linode.com/docs/products/storage/object-storage/#availability

[GoLand logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.svg
[jetbrains_goland_url]: https://www.jetbrains.com/go/