# s3packer - A configurable profile-based S3 backup and upload tool.

**CLI for Linux/MacOS**  **supports** Amazon S3 **|** Oracle Cloud Object Storage **|** Linode (Akamai) Object Storage

---

## AWS Specific Options

### Providers settings

AWS requires some authentication fields to be filled in before files can be uploaded.

#### Use an AWS CLI Profile to authenticate

| Provider | Acceptable Values | Required | Description                                                        |
|:---------|:------------------|:---------|:-------------------------------------------------------------------|
| Use      | aws               | Y        | Tell s3packer to use AWS                                           |
| Profile  | any string        | Y        | The profile name to pull auth details from, see ~/.aws/credentials |

```yaml
Provider:
  Use: aws
  Profile: "BackupProfile"
```

#### Use an AWS API Key/Secret pair to authenticate

| Provider | Acceptable Values | Required | Description                                  |
|:---------|:------------------|:---------|:---------------------------------------------|
| Use      | aws               | Y        | Tell s3packer to use AWS                     |
| Key      | any string        | Y        | The api key used to authenticate with AWS    |
| Secret   | any string        | Y        | The api secret used to authenticate with AWS |

```yaml
Provider:
  User: aws
  Key: "my-key-value"
  Secret: "my-key-secret-value"
```

---

### Storage and ACL Options

Configure your object ACLs and the storage type.

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