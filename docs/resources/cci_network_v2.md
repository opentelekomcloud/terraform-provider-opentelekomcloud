---
subcategory: "Cloud Container Instance (CCI)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cci_network_v2"
sidebar_current: "docs-opentelekomcloud-resource-cci-network-v2"
description: |-
  Manages a CCI Network resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CCI network you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-container-instance/api-ref/proprietary_apis/index.html)

# opentelekomcloud_cci_network_v2

Manages a CCI Network resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "project_id" {}
variable "domain_id" {}
variable "subnet_id" {}
variable "security_group_id" {}

resource "opentelekomcloud_cci_network_v2" "test" {
  namespace = var.namespace
  name      = var.name

  annotations = {
    "yangtse.io/project-id"                 = var.project_id
    "yangtse.io/domain-id"                  = var.domain_id
    "yangtse.io/warm-pool-size"             = "10"
    "yangtse.io/warm-pool-recycle-interval" = "2"
  }

  subnets {
    subnet_id = var.subnet_id
  }

  security_group_ids = [var.security_group_id]
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, String, ForceNew) Specifies the namespace of the CCI network.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the CCI network.
  Changing this creates a new resource.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI network.
  Annotations is an unstructured key value map that may be set by external tools to store and retrieve arbitrary
  metadata.

* `ip_families` - (Optional, List, ForceNew) Specifies the IP families of the CCI network.
  The value can be `IPv4` or `IPv6`. When IPv6 is enabled, the value can be `["IPv4", "IPv6"]`.
  Changing this creates a new resource.

* `security_group_ids` - (Optional, List) Specifies the security group IDs of the CCI network.

* `subnets` - (Optional, List, ForceNew) Specifies the subnets of the CCI network.
  Changing this creates a new resource.
  The [subnets](#block_subnets) structure is documented below.

<a name="block_subnets"></a>
The `subnets` block supports:

* `subnet_id` - (Required, String, ForceNew) Specifies the subnet ID (the network ID of the subnet) of the
  CCI network. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<namespace>/<name>`.

* `region` - The region of the CCI network.

* `api_version` - The API version of the CCI network.

* `kind` - The kind of the CCI network.

* `creation_timestamp` - The creation timestamp of the CCI network. The value is in RFC3339 format and is in UTC.

* `finalizers` - The finalizers of the CCI network. Must be empty before the object is deleted from the registry.

* `resource_version` - The resource version of the CCI network.

* `uid` - The uid of the CCI network.

* `status` - The status of the CCI network.
  The [status](#attrblock_status) structure is documented below.

<a name="attrblock_status"></a>
The `status` block supports:

* `status` - The status of the CCI network. The value can be `Ready`, `Failed`, or `IPInsufficient`.

* `conditions` - The conditions of the CCI network.
  The [conditions](#attrblock_status_conditions) structure is documented below.

* `subnet_attrs` - The subnet attributes of the CCI network.
  The [subnet_attrs](#attrblock_status_subnet_attrs) structure is documented below.

<a name="attrblock_status_conditions"></a>
The `conditions` block supports:

* `type` - The type of the CCI network condition.

* `status` - The status of the CCI network condition. The value can be `True`, `False`, or `Unknown`.

* `last_transition_time` - The last time the condition transitioned from one status to another.

* `reason` - The reason for the condition's last transition.

* `message` - A human-readable message indicating details about the transition.

<a name="attrblock_status_subnet_attrs"></a>
The `subnet_attrs` block supports:

* `network_id` - The network ID of the subnet.

* `subnet_v4_id` - The subnet IPv4 ID.

* `subnet_v6_id` - The subnet IPv6 ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The CCI Network can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import opentelekomcloud_cci_network_v2.test <namespace>/<name>
```
