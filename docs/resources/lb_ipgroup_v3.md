---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_ipgroup_v3

Manages a Dedicated Load Balancer IP address group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_ipgroup_v3" "group_1" {
  name        = "group_1"
  description = "group description"

  ip_list {
    ip          = "192.168.50.10"
    description = "one"
  }
  ip_list {
    ip          = "192.168.100.10"
    description = "two"
  }
  ip_list {
    ip          = "192.168.150.10"
    description = "three"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the IP address group name.

* `description` - (Optional) Provides supplementary information about the IP address group.

* `project_id` - (Optional) Specifies the project ID of the IP address group.

* `ip_list` - (Required) Specifies the IP addresses or CIDR blocks in the IP address group. [] indicates any IP address.
  * `ip` - (Required) Specifies the IP addresses in the IP address group.
    IPv6 is unsupported. The value cannot be an IPv6 address.
  * `description` - (Optional) Provides remarks about the IP address group.

## Attributes Reference

In addition, the following attributes are exported:

* `listeners` - Lists the IDs of listeners with which the IP address group is associated.

* `updated_at` - Indicates the update time.

* `created_at` - Indicates the creation time.

## Import

Ip groups can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_ipgroup_v3.group_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
