---
subcategory: "Enterprise Router (ER)"
---

Up-to-date reference of API arguments for Enterprise Router you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/enterprise-router/api-ref/apis/enterprise_routers/index.html#enterpriserouterinstance).

# opentelekomcloud_er_instance_v3

Manages an ER instance resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "router_name" {}
variable "availability_zones" {
  type = list(string)
}

resource "opentelekomcloud_er_instance_v3" "test" {
  availability_zones = var.availability_zones

  name = var.router_name
  asn  = 64512
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The router name.
  The name can contain 1 to 64 characters, only letters, digits, underscore (_) and hyphens (-) are allowed.

* `availability_zones` - (Required, List) The availability zone list where the ER instance is located.
  The maximum number of availability zone is two. Select two AZs to configure active-active deployment for high
  availability which will ensure reliability and disaster recovery.

* `asn` - (Required, Int, ForceNew) The BGP AS number of the ER instance.
  The valid value is range from `64,512` to `65534` or range from `4,200,000,000` to `4,294,967,294`.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) The description of the ER instance.
  The description contain a maximum of 255 characters, and the angle brackets (< and >) are not allowed.

* `enable_default_propagation` - (Optional, Bool) Whether to enable the propagation of the default route table.
  The default value is **false**.

* `enable_default_association` - (Optional, Bool) Whether to enable the association of the default route table.
  The default value is **false**.

* `auto_accept_shared_attachments` - (Optional, Bool) Whether to automatically accept the creation of shared
attachment.
  The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - Current status of the router.

* `created_at` - The creation time.

* `updated_at` - The latest update time.

* `region` - Specifies the region of the ER instance.

* `default_propagation_route_table_id` - The ID of the default propagation route table.

* `default_association_route_table_id` - The ID of the default association route table.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 5 minutes.

## Import

The router instance can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_er_instance_v3.test 0ce123456a00f2591fabc00385ff1234
```
