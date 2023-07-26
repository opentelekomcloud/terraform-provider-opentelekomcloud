---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF datamasking rule you can get at
`https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/dedicated_instance_management/index.html`.

# opentelekomcloud_waf_dedicated_instance_v1

Manages a WAF dedicated instance resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "my_subnet"
}

data "opentelekomcloud_networking_secgroup_v2" "default_secgroup" {
  name = "default"
}

resource "opentelekomcloud_waf_dedicated_instance_v1" "wafd_1" {
  name              = "wafd-instance-1"
  availability_zone = "eu-de-01"
  specification     = "waf.instance.professional"
  flavor            = "s2.large.2"
  architecture      = "x86"
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id

  security_group = [
    data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) Region where a dedicated engine is to be created. If omitted, the
  provider-level region will be used. Changing this setting will create a new instance.

* `name` - (Required, String) The name of WAF dedicated instance. Duplicate names are allowed, we suggest to keeping the
  name unique.

* `availability_zone` - (Required, ForceNew) AZ where the dedicated engine is to be created. Changing this will create a new instance.

* `specification` - (Required, ForceNew) Specifications of the dedicated engine version. Values are:
  + `waf.instance.professional` - The professional edition, throughput: 100 Mbit/s; QPS: 2,000 (Reference only).
  + `waf.instance.enterprise` - The enterprise edition, throughput: 500 Mbit/s; QPS: 10,000 (Reference only).

* `flavor` - (Required, ForceNew) ID of the specifications of the ECS hosting the dedicated engine.
  You can go to the management console and confirm supported specifications. Changing this will create a new instance.

* `vpc_id` - (Required, ForceNew) ID of the VPC where the dedicated engine is located. Changing this will create a new
  instance.

* `subnet_id` - (Required, ForceNew) ID of the VPC subnet where the dedicated engine is located.
  Subnet_id has the same value as network_id obtained by calling the OpenStack APIs. Changing this will create a
  new instance.

* `security_group` - (Required, ForceNew) ID of the security group where the dedicated engine is located.
  Changing this will create a new instance.

* `architecture` - (Optional, String, ForceNew) Dedicated engine CPU architecture. Default value is `x86`.
  Changing this will create a new instance.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the instance.

* `server_id` - The id of the instance server.

* `service_ip` - The ip of the instance service.

* `status` - Running status of the dedicated engine.
  The value can be:
  + `0` - Instance is creating.
  + `1` - Instance has created.
  + `2` - Instance is deleting.
  + `3` - Instance has deleted.
  + `4` - Instance create failed.
  + `5` - Instance is frozen.
  + `6` - Instance in abnormal state.
  + `7` - Instance in updating.
  + `8` - Instance update failed.

* `access_status` - The access status of the instance.
  + `0`: inaccessible
  + `1`: accessible.

* `billing_status` - Billing status of dedicated WAF engine. The value can be `0`, `1`, or `2`.
  + `0`: The billing is normal.
  + `1`: The billing account is frozen. Resources and data will be retained, but the cloud services cannot be used by the account.
  + `2`: The billing is terminated. Resources and data will be cleared.

* `upgradable` - The instance is to support upgrades. `false`: Cannot be upgraded, `true`: Can be upgraded.

* `created_at` - Timestamp when the dedicated WAF engine was created.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minute.
* `delete` - Default is 20 minute.

## Import

WAF dedicated instance can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_waf_dedicated_instance_v1.wafd <id>
```
