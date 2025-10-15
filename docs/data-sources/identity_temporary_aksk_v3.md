---
subcategory: "Identity and Access Management (IAM)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_identity_temporary_aksk_v3"
sidebar_current: "docs-opentelekomcloud-datasource-identity-temporary-aksk-v3"
description: |-
  Get temporary AK/SK credentials for accessing OpenTelekomCloud services
---

# opentelekomcloud_identity_temporary_aksk_v3

Use this data source to generate temporary access credentials (AK/SK and security token)
for accessing OpenTelekomCloud services. These credentials are automatically generated
and have a limited lifetime, making them suitable for temporary access scenarios.

## Example Usage

### Basic Usage
```hcl
data "opentelekomcloud_identity_temporary_aksk_v3" "temp_creds" {}

output "access_key" {
  value     = data.opentelekomcloud_identity_temporary_aksk_v3.temp_creds.access
  sensitive = true
}
```

### With Custom Duration
```hcl
data "opentelekomcloud_identity_temporary_aksk_v3" "temp_creds" {
  duration_seconds = 3600 # 1 hour
}
```

### With Agency
```hcl
data "opentelekomcloud_identity_temporary_aksk_v3" "temp_creds" {
  duration_seconds = 7200 # 2 hours
  agency_name      = "my_agency"
}
```

## Argument Reference

The following arguments are supported:

* `duration_seconds` - (Optional) Validity period (in seconds) of the temporary credentials.
  The value ranges from 900 (15 minutes) to 86400 (24 hours). The default value is 900 seconds (15 minutes).

* `agency_name` - (Optional) Name of the agency created by a delegating party.
  When specified, the temporary credentials will be generated with the permissions of the agency.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The access key, which serves as the unique identifier for the temporary credentials.

* `access` - The temporary Access Key (AK) for authentication.

* `secret` - The temporary Secret Key (SK) for authentication.

* `security_token` - The security token used for subsequent API calls. This token must be included
  in requests when using the temporary AK/SK.

* `expires_at` - The expiration time of the temporary credentials in ISO 8601 format.

## Notes

* Temporary credentials are generated each time the data source is read, so they will be regenerated
  on every `terraform refresh` or `terraform apply`.

* The temporary credentials consist of three parts that must be used together: AK, SK, and security token.
