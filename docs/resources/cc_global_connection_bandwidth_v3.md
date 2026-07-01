---
subcategory: "Cloud Connect (CC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cc_global_connection_bandwidth_v3"
sidebar_current: "docs-opentelekomcloud-resource-cc-global-connection-bandwidth-v3"
description: |-
  Manages a Cloud Connect Global Connection Bandwidth resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for Cloud Connect Global Connection Bandwidth you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-connect/api-ref/api/global_connection_bandwidths/index.html)

# opentelekomcloud_cc_global_connection_bandwidth_v3

Manages a Cloud Connect (CC) global connection bandwidth resource within OpenTelekomCloud.

A global connection bandwidth provides the network capacity that connects instances (such as cloud connections,
global EIPs, and central networks) across regions. It can be dedicated to a single instance or shared by
multiple instances.

-> The global connection bandwidth APIs are account-level (domain-scoped). Make sure the credentials used by the
provider have permission to manage Cloud Connect resources.

## Example Usage

```hcl
resource "opentelekomcloud_cc_global_connection_bandwidth_v3" "test" {
  name        = "gcb-demo"
  description = "managed by terraform"
  type        = "Region"
  bordercross = false
  charge_mode = "bwd"
  size        = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the global connection bandwidth.

* `bordercross` - (Required, Bool, ForceNew) Whether the bandwidth is used for cross-border connections.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) The bandwidth type. The value can be **Area**, **TrsArea**,
  **SubArea** or **Region**. Changing this parameter will create a new resource.

* `charge_mode` - (Required, String) The billing option. The value can be **bwd** (billing by bandwidth),
  **95** (standard 95th percentile billing) or **95avr** (enhanced 95th percentile billing).

* `size` - (Required, Int) The bandwidth capacity, in Mbit/s. The value ranges from **2** to **300**.

* `description` - (Optional, String) The description of the global connection bandwidth.
  The angle brackets (`<` and `>`) are not allowed.

* `enterprise_project_id` - (Optional, String, ForceNew) The ID of the enterprise project that the global
  connection bandwidth belongs to. Changing this parameter will create a new resource.

* `sla_level` - (Optional, String) The service class of the bandwidth. The value can be **Pt** (Platinum),
  **Au** (Gold) or **Ag** (Silver).

* `local_area` - (Optional, String, ForceNew) The local access point. The value can contain 1 to 64
  characters, only letters, digits, underscores (_), hyphens (-) and periods (.) are allowed.
  Changing this parameter will create a new resource.

* `remote_area` - (Optional, String, ForceNew) The remote access point. The value can contain 1 to 64
  characters, only letters, digits, underscores (_), hyphens (-) and periods (.) are allowed.
  Changing this parameter will create a new resource.

* `spec_code_id` - (Optional, String) The UUID of the line specification code.

* `binding_service` - (Optional, String) The service binding type. The value can be **CC**, **GEIP**,
  **GCN**, **GSN** or **ALL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format. Equal to the global connection bandwidth ID.

* `domain_id` - The ID of the account that the global connection bandwidth belongs to.

* `local_site_code` - The local access point code.

* `remote_site_code` - The remote access point code.

* `admin_state` - The status of the global connection bandwidth. The value can be **NORMAL** or **FREEZED**.

* `frozen` - Whether the global connection bandwidth is frozen.

* `enable_share` - Whether the bandwidth can be shared by multiple instances.

* `eps_id` - The ID of the enterprise project that the global connection bandwidth belongs to.

* `created_at` - The creation time of the global connection bandwidth.

* `updated_at` - The latest update time of the global connection bandwidth.

* `instances` - The instances associated with the global connection bandwidth.
  The [instances](#gcb_instances) structure is documented below.

* `directional_connections` - The directional connections of the global connection bandwidth.
  The [directional_connections](#gcb_directional_connections) structure is documented below.

* `region` - The region in which the global connection bandwidth is managed.

<a name="gcb_instances"></a>
The `instances` block supports:

* `id` - The instance ID.

* `type` - The instance type.

* `region_id` - The region ID of the instance.

* `project_id` - The project ID of the instance.

<a name="gcb_directional_connections"></a>
The `directional_connections` block supports:

* `id` - The connection ID.

* `name` - The connection name.

* `local_site_code` - The local access point code.

* `remote_site_code` - The remote access point code.

## Import

The global connection bandwidth can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_cc_global_connection_bandwidth_v3.test 0ce123456a00f2591fabc00385ff1234
```
