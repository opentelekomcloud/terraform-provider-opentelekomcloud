---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_networks_v2"
sidebar_current: "docs-opentelekomcloud-data-source-cci-networks-v2"
description: |-
  Get the list of CCI v2 networks within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI network you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_networks_v2

Use this data source to get the list of CCI v2 networks under a namespace within OpenTelekomCloud.

## Example Usage

```hcl
variable "namespace" {}

data "opentelekomcloud_cci_networks_v2" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String) Specifies the namespace to which the networks belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The region in which the networks are queried.

* `networks` - The list of networks. The [networks](#networks) structure is documented below.

<a name="networks"></a>
The `networks` block supports:

* `name` - The name of the network.

* `namespace` - The namespace to which the network belongs.

* `annotations` - The annotations of the network.

* `labels` - The labels of the network.

* `creation_timestamp` - The creation timestamp of the network.

* `resource_version` - The resource version of the network.

* `finalizers` - The finalizers of the network.

* `uid` - The uid of the network.

* `ip_families` - The IP families of the network.

* `security_group_ids` - The security group IDs bound to the network.

* `subnets` - The subnets of the network. The [subnets](#subnets) structure is documented below.

* `status` - The status of the network. The [status](#status) structure is documented below.

<a name="subnets"></a>
The `subnets` block supports:

* `subnet_id` - The ID of the subnet.

<a name="status"></a>
The `status` block supports:

* `status` - The current state of the network.

* `conditions` - The list of network conditions. The [conditions](#conditions) structure is documented below.

* `subnet_attrs` - The list of subnet attributes. The [subnet_attrs](#subnet_attrs) structure is documented below.

<a name="conditions"></a>
The `conditions` block supports:

* `type` - The type of the condition.

* `status` - The status of the condition.

* `last_transition_time` - The last transition time of the condition.

* `reason` - The reason for the condition's last transition.

* `message` - The human readable message indicating details about the transition.

<a name="subnet_attrs"></a>
The `subnet_attrs` block supports:

* `network_id` - The ID of the underlying neutron network.

* `subnet_v4_id` - The IPv4 subnet ID.

* `subnet_v6_id` - The IPv6 subnet ID.
