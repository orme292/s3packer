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

s3packer profiles are written in the YAML format. To set one up, you just need to fill out a few fields, and you’ll be good to go!

First, make sure you specify that you're using Version 4 of the profile format and specify OCI as the object storage provider:

```yaml
Version: 4
Provider: oci
```

---
## Authentication
> 💡 You can remove the **AWS** section from the profile.

**s3packer** handles OCI authentication that is generated with the OCI-CLI. You can specify a profile that's already set 
up in your `~/.oci/config` file.

For info on setting up the OCI-CLI, check out the [Oracle Cloud documentation][oci_cli_url].

The compartment field can be left blank. It is only required if s3packer has to create a bucket. If s3packer is creating the bucket,
and no compartment is specified, it will use the tenancy root as the compartment.

```yaml
OCI:
  Profile: "default"
  Compartment: "ocid1.compartment.oc1..aaaaaaa..."
```

Next, configure the bucket. The `name` field is required. The `region` field **must contain something**, but
it can be any string. s3packer can find the bucket by name. If it's creating the bucket, it will be 
created in the tenancy's default region.

```yaml
Bucket:
  Name: "free-data"
  Region: "eu-zurich-1"
```

Finally, tell s3packer what you want to upload. You can specify folders, directories, or individual files. (You can call
it the `Folders` section or the `Directories` section, it doesn't matter.)

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

You can also add tags to your files. Just add a `Tagging` section to your profile, like this:

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