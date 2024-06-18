# s3packer

**Linux/MacOS |||||** Amazon S3 **|** Oracle Cloud Object Storage **|** Akamai (Linode) Object Storage

---
## Using s3packer with Oracle Cloud (OCI)

◀️ [Back to s3packer][s3packer_readme_url]

---

## Creating a new Profile

s3packer can create a base profile to help get you started. To create one, use the `--create` flag:

```bash
$ s3packer --create="my-new-oci-profile.yaml"
```

## Setting up a Profile

s3packer profiles are written in YAML. To set one up, you just need to fill out a few fields, and you’ll be good to go!

First, make sure you specify that you're using Version 4 of the profile format:

```yaml
Version: 5
Provider:
  Use: oracle
```

---
## Authentication

**s3packer** handles OCI authentication using the config generated with the OCI-CLI. You can specify a profile that's
already set up in your `~/.oci/config` file.

For info on setting up the OCI-CLI, check out the [Oracle Cloud documentation][oci_cli_url].

```yaml
Version: 5
Provider:
  Use: oracle
  Profile: DEFAULT
```

Under the OCI field, specify a compartment. It is only required if s3packer has to create a bucket. If s3packer is
creating the bucket and no compartment is specified, it will create the bucket in the tenancy's root compartment.

```yaml
OCI:
  Compartment: "ocid1.compartment.oc1..aaaaaaa..."
```

Next, configure the bucket. The `name` field is required. The `region` field **must contain something**, but
it can be any string. s3packer can find the bucket by name. If it's creating the bucket, it will be 
created in the tenancy's default region.

```yaml
Bucket:
  Create: true
  Name: "free-data"
  Region: "eu-zurich-1"
```

Finally, tell s3packer what you want to upload. You can specify directories or individual files. When you specify a
directory, s3packer will traverse all subdirectories.

```yaml
Dirs:
  - "/Users/forrest/docs/stocks/apple"
  - "/Users/jenny/docs/song_lyrics"
Files:
  - "/Users/forrest/docs/objJob-application-lawn-mower.pdf"
  - "/Users/forrest/docs/dr-pepper-recipe.txt"
  - "/Users/jenny/letters/from-forrest.docx"
```

--- 

### Tags

You can also add tags to your files. Just add a `Tags` section to your profile:

```yaml
Tags:
  Author: "Forrest Gump"
  Year: 1994
```
---

### OCI Specific Options

Configure your object storage tier.

```yaml
OCI:
  Storage: "standard"
```

**Storage** <br/>
The default is `STANDARD`, but you can use any of the following storage classes:
- `Standard`
- `InfrequentAccess`
- `Archive`

Read more about OCI's storage tiers here: [https://docs.oracle.com/en-us/iaas/Content/Object/Concepts/understandingstoragetiers.htm][oci_tier_url]

---

### Extra Options

You can also customize how your files are stored, accessed, tagged, and uploaded using these options.

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

```yaml
Tagging:
  OriginPath: true
  ChecksumSHA256: false
```

**OriginPath** <br/>
This is `true` by default. Every object uploaded will be tagged with the full absolute path of the file on the local
filesystem. This is useful if you want to be able to trace the origin of a file in S3. The tag name will be
`s3packer-origin-path`. Oracle may add a prefix to the tag name.

**ChecksumSHA256** <br/>
This is `true` by default. Every object uploaded will be tagged with the file's calculated SHA256 checksum. The tag name
will be `s3packer-checksum-sha256`. Oracle may add a prefix to the tag name.

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
[oci_tier_url]: https://docs.oracle.com/en-us/iaas/Content/Object/Concepts/understandingstoragetiers.htm
[oci_cli_url]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm#InstallingCLI__macos_homebrew

[GoLand logo]: https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.svg
[jetbrains_goland_url]: https://www.jetbrains.com/go/