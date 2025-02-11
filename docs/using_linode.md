# s3p - A configurable profile-based S3 backup and upload tool.

**CLI for Linux/MacOS**  **supports** Amazon S3 **|** Google Cloud Storage **|** Linode (Akamai) Object Storage
**|** Oracle Cloud Object Storage

---

## Linode specific options

### Providers settings

**s3p** uses object storage access keys to authenticate with Linode. You can find and generate access keys in the
Cloud Manager. For info on generating new access keys, check out the [Linode Object Storage Guide][akamai_auth_url].

#### Use a LinodeObjects API Key/Secret pair to authenticate

| Provider | Acceptable Values | Required | Description                                            |
|:---------|:------------------|:---------|:-------------------------------------------------------|
| Use      | linode            | Y        | Tell s3p to use Linode                                 |
| Key      | any string        | Y        | The api key used to authenticate with LinodeObjects    |
| Secret   | any string        | Y        | The api secret used to authenticate with LinodeObjects |

```yaml
Provider:
  Use: linode
  Key: "my-key-value"
  Secret: "my-key-secret-value"
```

---

### LinodeObjects Region

Specify a region (separate from the Bucket heading) that s3p can use to generate the LinodeObjects endpoint

```yaml
Linode:
  Region: "se-sto-1"
```

`Region` should contain the region short-code. When you create a bucket in the Cloud Manager, the short code will be
listed in the region dropdown. You can also check Linode's documentation for a list of region short-codes:
[Linode Region List][akamai_region_list_url].

---

### ACLs

LinodeObjects has support for ACLs, but they aren't supported by s3p. All created buckets and objects will use the
default canned AWS ACL `public-read`

---

### Tags

LinodeObjects doesn't provide info on its support for object tagging (as far as I could find)

[akamai_auth_url]: https://www.linode.com/docs/products/storage/object-storage/guides/access-keys/

[akamai_region_list_url]: https://www.linode.com/docs/products/storage/object-storage/#availability