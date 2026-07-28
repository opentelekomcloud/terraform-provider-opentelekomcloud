---
subcategory: ""
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: S3 Backends Guide"
description: |-
  Additional documentation about s3 backend configuration within OpenTelekomCloud.
---

> Starting with Terraform v1.11.2, changes in the bundled AWS SDK can cause state uploads to fail with some third-party S3 implementations, potentially resulting in divergent `.tfstate` files. The SDK defaults to calculating and validating checksums whenever an S3 API operation supports them. This can add optional CRC32 checksums to uploads and request checksum validation for downloads, even when `skip_s3_checksum = true`. Some S3-compatible APIs do not fully support these newer checksum headers and can reject the request with errors such as `XAmzContentSHA256Mismatch`.
>
> Configure the SDK to use checksums only when an operation requires them:
>
> ```shell
> export AWS_REQUEST_CHECKSUM_CALCULATION=when_required
> export AWS_RESPONSE_CHECKSUM_VALIDATION=when_required
> ```
>
> `AWS_REQUEST_CHECKSUM_CALCULATION=when_required` prevents the SDK from automatically adding an optional CRC32 checksum to requests such as `PutObject`. `AWS_RESPONSE_CHECKSUM_VALIDATION=when_required` prevents automatic response-checksum negotiation and validation unless Terraform explicitly enables it. These settings do not disable checksums required by an API; they avoid optional checksum behavior that an S3-compatible implementation might not support.
>
> If a state upload has already failed, do not run another apply because it can create divergent state histories. Recover the state from `errored.tfstate` as instructed by Terraform, or manually repair the state and roll back to an earlier Terraform version. For more details, see the [Terraform issue](https://github.com/hashicorp/terraform/issues/37130) and [AWS checksum configuration](https://docs.aws.amazon.com/sdkref/latest/guide/feature-dataintegrity.html).

A `backend` defines where Terraform stores its [state](https://developer.hashicorp.com/terraform/language/state) data files.

The main terraform reference you can get at:
`https://developer.hashicorp.com/terraform/language/settings/backends/configuration`.

## Default Backend
If a configuration includes no `backend` block, Terraform defaults to using the `local` backend, which stores state as a plain file in the current working directory.

## Initialization
When you change a backend's configuration, you must run `terraform init` again to validate and configure the backend before you can perform any plans, applies, or state operations.

After you initialize, Terraform creates a `.terraform/` directory locally. This directory contains the most recent backend configuration, including any authentication parameters you provided to the Terraform CLI. Do not check this directory into Git, as it may contain sensitive credentials for your remote backend.

The local backend configuration is different and entirely separate from the `terraform.tfstate` file that contains state data about your real-world infrastructure. Terraform stores the `terraform.tfstate` file in your remote backend.

When you change backends, Terraform gives you the option to migrate your state to the new backend. This lets you adopt backends without losing any existing state.

## Using a T-Cloud Public (former OpenTelekomCloud) S3 Backend Block
To configure a backend, add a nested backend block within the top-level terraform block. The following example configures the remote backend for S3.
T-Cloud Public (former OpenTelekomCloud) S3 endpoints:
 - `https://obs.eu-de.otc.t-systems.com/`
 - `https://obs.eu-nl.otc.t-systems.com/`
 - `https://obs.eu-ch2.sc.otc.t-systems.com/`

Also, after latest terraform updates uso should also specify `secret_key` and `access_key` in `backend` block.

```hcl
terraform {
  required_version = ">= 1.6.3"

  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.36.0"
    }
  }
  backend "s3" {
    endpoints = {
      s3 = "https://obs.eu-de.otc.t-systems.com/"
    }
    key                         = "terraform_state/test"
    bucket                      = "tf-test-bucket"
    region                      = "eu-de"
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    skip_metadata_api_check     = true
    skip_s3_checksum            = true
    secret_key                  = "secret"
    access_key                  = "access"
  }
}
```

## Support for "S3 Compatible" Storage Providers

Support for S3 Compatible storage providers is offered by `HashiCorp` as “best effort”.
`HashiCorp` only tests the `s3` backend against `Amazon S3`, so cannot offer any guarantees when using an alternate provider.
