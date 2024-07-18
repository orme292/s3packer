# s3packer - A configurable profile-based S3 backup and upload tool.

**CLI for Linux/MacOS**  **supports** Amazon S3 **|** Oracle Cloud Object Storage **|** Linode (Akamai) Object Storage

---

## Oracle OCI specific options

### Providers settings

**s3packer** handles OCI authentication using the config generated with the OCI-CLI. You can specify a profile that's
already set up in your `~/.oci/config` file.

For info on installing and configuring the OCI-CLI, see the [Oracle Cloud documentation][oci_cli_url].

| Provider | Acceptable Values | Required | Description                                          |
|:---------|:------------------|:---------|:-----------------------------------------------------|
| Use      | oci               | Y        | Tell s3packer to use Oracle OCI                      |
| Profile  | any string        | Y        | The oci-cli profile name to use, see `~/.oci/config` |

```yaml
Provider:
  Use: oci
  Profile: "DEFAULT"
```

---

### OCI Compartment

Under the OCI heading, specify a compartment. It is only required if s3packer has to create a bucket. If s3packer is
creating the bucket and no compartment is specified, it will create the bucket in the tenancy's root compartment.

| OCI         | Acceptable Values | Required | Description                                                                          |
|:------------|:------------------|:---------|:-------------------------------------------------------------------------------------|
| Compartment | OCID              | N        | The OCID of the compartment that s3packer would create a bucket in, if configured to |
| Storage     | see below         | N        | The storage tier that will be specified when uploading objects                       |

```yaml
OCI:
  Compartment: "ocid1.compartment.oc1..aaaaaaa..."
  Storage: "standard"
```

**Storage** <br/>
The default is `STANDARD`, but you can use any of the following storage classes:

- `Standard`
- `InfrequentAccess`
- `Archive`

Be sure to read about OCI's storage tiers
here: [https://docs.oracle.com/en-us/iaas/Content/Object/Concepts/understandingstoragetiers.htm][oci_tier_url]

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

**SPECIAL NOTE ABOUT OCI TAGS**<br/>
All OCI tags are prefixed with the string "opc-meta-". s3packer does and the OCI SDK do this automatically, so DO NOT
add this to your tag keys.


[oci_tier_url]: https://docs.oracle.com/en-us/iaas/Content/Object/Concepts/understandingstoragetiers.htm

[oci_cli_url]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm#InstallingCLI__macos_homebrew