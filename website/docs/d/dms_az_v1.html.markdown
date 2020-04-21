---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_az_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dms-az-v1"
description: |-
  Get information on an OpenTelekomCloud dms az.
---

# opentelekomcloud\_dms\_az_v1

Use this data source to get the ID of an available OpenTelekomCloud dms az.

## Example Usage

```hcl
data "opentelekomcloud_dms_az_v1" "az1" {
  name = "eu-de-01"
  port = "8002"
}
```

## Argument Reference

* `name` - (Required) Indicates the name of an AZ.

* `port` - (Optional) Indicates the port number of an AZ.

* `code` - (Optional) Indicates the code of an AZ.

## Attributes Reference

`id` is set to the ID of the found az. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `port` - See Argument Reference above.
* `code` - See Argument Reference above.
