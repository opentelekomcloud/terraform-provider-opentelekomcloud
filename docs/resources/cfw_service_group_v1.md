---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_service_group_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-service-group-v1"
description: |-
  Manages a CFW Service group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW service group you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/service_group_management/index.html)

# opentelekomcloud_cfw_service_group_v1

Manages a CFW Service Group resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable object_id {}

resource "opentelekomcloud_cfw_service_group_v1" "group_1" {
  object_id = var.object_id
  name      = "test-acc-tf-service-group"
}
```

## Argument Reference

The following arguments are supported:

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `name` - (Required, String) Specifies the CFW Service group name. The CFW service group name of the same type is unique in the same firewall.

* `description` - (Optional, String) Specifies the description of the service group.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the service group ID.
* `service_set_type` - Indicates the Service group type: `0` (user-defined service group), `1` (common web service), `2` (common remote login and ping), or `3` (common database).

## Import

CFW Service Group V1 resource can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_cfw_service_group_v1.group_1 b4cd6aeb0b7445d3bf271457c6941544in09
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_cfw_service_group_v1" "group_1" {
  # ...

  lifecycle {
    ignore_changes = [
      object_id,
    ]
  }
}
```
