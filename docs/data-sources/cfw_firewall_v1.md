---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_firewall_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-firewall-v1"
description: |-
  Get details about a CFW Firewall Instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW firewall instance you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/firewall_management/index.html)

# opentelekomcloud_cfw_firewall_v1

Get details about a CFW Firewall Instance resource within OpenTelekomCloud.

## Example Usage: Creating a basic CFW firewall instance
```hcl
variable firewall_id {}

resource "opentelekomcloud_cfw_firewall_v1" "firewall_1" {
  id           = var.firewall_id
  service_type = "0"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required, String) Specifies the Firewall instance ID.

* `service_type` - (Optional, String) Specifies the Firewall protection type. Currently, its value can only be `0` (Internet protection).

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `name` - Specifies the CFW firewall instance name. The CFW firewall instance name of the same
  type is unique in the same tenant.
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
