---
subcategory: "Elastic Cloud Server (ECS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_flavor_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-flavor-v2"
description: |-
  Get the ID of an available ECS flavor from OpenTelekomCloud
---

Up-to-date reference of API arguments for ECS flavor you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/native_openstack_nova_apis/flavor_management/querying_ecs_flavors.html#en-us-topic-0065817705)

# opentelekomcloud_compute_flavor_v2

Use this data source to get the ID of an available OpenTelekomCloud flavor.

## Example Usage

```hcl
data "opentelekomcloud_compute_flavor_v2" "medium-s2" {
  vcpus         = 1
  ram           = 4096
  resource_type = "IOoptimizedS2"
}
```

## Argument Reference

-> If multiple flavors have found only the first will be returned.

* `region` - (Optional) The region in which to obtain the V2 Compute client.
  If omitted, the `region` argument of the provider is used.

* `flavor_id` - (Optional) The ID of the flavor. Conflicts with the `name`,
  `min_ram` and `min_disk`

* `name` - (Optional) The name of the flavor. Conflicts with the `flavor_id`.

* `min_ram` - (Optional) The minimum amount of RAM (in megabytes). Conflicts
  with the `flavor_id`.

* `ram` - (Optional) The exact amount of RAM (in megabytes).

* `min_disk` - (Optional) The minimum amount of disk (in gigabytes). Conflicts
  with the `flavor_id`.

* `availability_zone` - (Optional) Whether flavor should be in `normal` state.

* `resource_type` - (Optional) Flavor resource type.

* `disk` - (Optional) The exact amount of disk (in gigabytes).

* `vcpus` - (Optional) The amount of VCPUs.

* `swap` - (Optional) The amount of swap (in gigabytes).

* `rx_tx_factor` - (Optional) The `rx_tx_factor` of the flavor.

## Attributes Reference

`id` is set to the ID of the found flavor. In addition, the following attributes
are exported:

* `extra_specs` - Key/Value pairs of metadata for the flavor.
