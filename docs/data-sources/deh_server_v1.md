---
subcategory: "Dedicated Host (DEH)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_deh_server_v1"
sidebar_current: "docs-opentelekomcloud-datasource-deh-server-v1"
description: |-
Get details about the server on a specified Dedicated Host from OpenTelekomCloud
---

Up-to-date reference of API arguments for DEH ECSs you can get at
[documentation portal](https://docs.otc.t-systems.com/dedicated-host/api-ref/api/querying_ecss_on_a_deh.html)

# opentelekomcloud_deh_server_v1

Use this data source to get details about the server on a specified Dedicated Host.

## Example Usage

```hcl
variable "deh_id" {}
variable "server_id" {}

data "opentelekomcloud_deh_server_v1" "deh_server" {
  id        = var.deh_id
  server_id = var.server_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the server on specified dedicated host.

* `dedicated_host_id` - (Optional) The Dedicated Host ID.

* `server_id` - (Optional) The Server ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `user_id` - The ID of the user to which the server belongs.

* `name` - The server name.

* `flavor` - The ID of server specifications.

* `metadata` - The metadata of the server.

* `status` - The status of the server.

* `tenant_id` - The ID of the tenant to which the server belongs.

* `addresses` - The network addresses of the server.
