---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_environments_v2"
sidebar_current: "docs-opentelekomcloud-datasource-apigw-environments-v2"
description: |-
  Get the all APIGW gateway environments from OpenTelekomCloud
---

# opentelekomcloud_apigw_environments_v2

Use this data source to query the environment list under the APIGW instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "environment_name" {}

data "opentelekomcloud_apigw_environments_v2" "test" {
  instance_id = var.instance_id
  name        = var.environment_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies an ID of the APIGW dedicated instance to which the API
  environment belongs.

* `name` - (Optional, String) Specifies the name of the API environment. The API environment name consists of 3 to 64
  characters, starting with a letter. Only letters, digits and underscores (_) are allowed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `environments` - List of APIGW environment details. The structure is documented below.

* `region` - The region in which queried the data source.

The `environments` block supports:

* `id` - ID of the APIGW environment.

* `name` - The environment name.

* `description` - The description about the API environment.

* `created_at` - Time when the APIGW environment was created, in RFC-3339 format.
