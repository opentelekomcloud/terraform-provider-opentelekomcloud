---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_firewall_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-firewall-v1"
description: |-
  Manages a CFW Firewall Instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW firewall instance you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/firewall_management/index.html)

# opentelekomcloud_cfw_firewall_v1

Manages a CFW Firewall Instance resource within OpenTelekomCloud.

## Example Usage: Creating a basic CFW firewall instance
```hcl
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  name = "test-acc-tf-firewall"
  flavor {
    version = "standard"
  }
  charge_info {
    charge_mode = "postPaid"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the CFW firewall instance name. The CFW firewall instance name of the same
  type is unique in the same tenant.

* `service_type` - (Optional, String, ForceNew) Specifies the Firewall protection type. Currently, its value can only be `0` (Internet protection).

* `flavor` - (Required, List, ForceNew) Specifies the Firewall specifications. The [flavor](#flavor) structure is documented below.

* `charge_info` - (Required, List, ForceNew) Specifies the billing type, which can be yearly/monthly or pay-per-use (default setting).
  The [charge_info](#charge_info) structure is documented below.

<a name="flavor"></a>
The `flavor` block supports:

* `version` - (Optional, String, ForceNew) Specifies the Firewall edition. Only the professional edition `standard` is supported.

<a name="charge_info"></a>
The `charge_info` block supports:

* `charge_mode` - (Optional, String, ForceNew) Specifies the Billing mode. The value can only be `postPaid` (case-sensitive), indicating pay-per-use billing.


## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the Firewall instance ID.
* `flavor` - Indicates the Firewall specifications. The [flavor](#flavor_res) structure is documented below.
* `enterprise_project_id` - Indicates the Enterprise project ID, which is the ID of a project planned based on organizations.
* `ha_type` - Indicates the Cluster type: 0 (active/standby), 1 (cluster). In active/standby mode, there are four nodes. Two active nodes form a cluster, and the other two are the standby of the active nodes. In cluster mode, only two nodes are started to form a cluster..
* `charge_mode` - Indicates the billing mode: 0 (yearly/monthly), 1 (pay-per-use).
* `engine_type` - Indicates the engine type. Its value can only be 1 (Hillstone engine).
* `protect_objects` - Indicates the protected object list. The [protect_objects](#protect_objects) structure is documented below.
* `status` - Indicates the firewall status: -1 (waiting for payment), 0 (creating), 1 (deleting), 2 (running), 3 (upgrading), 4 (deleted), 5 (frozen), 6 (creation failed), 7 (deletion failed), 8 (freezing failed), or 9 (being stored), 10 (storage failed), or 11 (upgrade failed).
* `is_old_firewall_instance` - Indicates whether an engine is old: true (yes), false (no)..
* `is_available_obs` - Indicates whether OBS is supported: true (yes), false (no).
* `is_support_threat_tags` - Indicates whether threat intelligence tags are supported: true (yes), false (no).
* `support_ipv6` - Indicates whether IPv6 is supported: true (yes), false (no).
* `feature_toggle` - Provides a map of features indicating whether a feature is enabled: true (yes), false (no).
* `resources` - Indicates the firewall resource list. The [resources](#resources) structure is documented below.
* `resource_id` - Indicates the Firewall resource ID, which is the same as `id`.
* `support_url_filtering` - Indicates whether website filtering is supported: true (yes), false (no).

<a name="flavor"></a>
The `flavor` block supports:

* `version` - See Argument Reference above.
* `version_code` - Indicates the firewall version. Its value can only be 1 (professional edition).
* `eip_count` - Indicates the number of EIPs.
* `vpc_count` - Indicates the number of VPCs.
* `bandwidth` - Indicates the bandwidth, in Mbits/s.
* `log_storage` - Indicates the log storage, in bytes.
* `default_bandwidth` - Indicates the default firewall bandwidth, in Mbits/s.
* `default_eip_count` - Indicates the default number of EIPs.
* `default_log_storage` - Indicates the default log storage, in bytes.
* `default_vpc_count` - Indicates the default number of VPCs.

<a name="protect_objects"></a>
The `protect_objects` block supports:

* `object_id` - Indicates the protected object ID. It is used to distinguish Internet border protection from VPC border protection after a CFW instance is created.
* `object_name` - Indicates the protected object name.
* `type` - Indicates the project type: 0 (north-south), 1 (east-west).


<a name="resources"></a>
The `resources` block supports:

* `resource_id` - Indicates the resource ID. It can be the firewall ID, bandwidth ID, EIP ID, VPC ID, or the ID returned after CBC callback.
* `cloud_service_type` - Indicates the Service type, which is used by CBC.
* `resource_type` - Indicates the resource type.
* `resource_spec_code` - Indicates the inventory unit code.
* `resource_size` - Indicates the resource quantity.
* `resource_size_measure_id` - Indicates the resource unit.

## Import

CFW Firewall V1 Instance can be imported using the CFW firewall instance ID, `id` and service type `service_type`, e.g.

```shell
terraform import opentelekomcloud_cfw_firewall_v1.firewall_1 b4cd6aeb0b7445d3bf271457c6941544in09/service_type
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  # ...

  lifecycle {
    ignore_changes = [
      flavor.version,
      charge_info,
    ]
  }
}
```
