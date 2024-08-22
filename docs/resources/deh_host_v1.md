---
subcategory: "Dedicated Host (DEH)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_deh_host_v1"
sidebar_current: "docs-opentelekomcloud-resource-deh-host-v1"
description: |-
  Manages a DEH Host resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DEH host you can get at
[documentation portal](https://docs.otc.t-systems.com/dedicated-host/api-ref/api)

# opentelekomcloud_deh_host_v1

Allocates a Dedicated Host to a tenant and set minimum required parameters for this Dedicated Host.

## Example Usage

```hcl
resource "opentelekomcloud_deh_host_v1" "deh_host" {
  name              = "high_performance_deh"
  auto_placement    = "on"
  availability_zone = "eu-de-02"
  host_type         = "h1"

  tags = {
    muh = "kuh"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the Dedicated Host.

* `auto_placement` - (Optional, String) Allows a instance to be automatically placed onto the available Dedicated Hosts. The default value is `on`.

* `availability_zone` - (Required, String, ForceNew) The Availability Zone to which the Dedicated Host belongs.

* `host_type` - (Required, String, ForceNew) The Dedicated Host type.
  Expected values are:
  * in `eu-de`: `c4`, `s2-medium`, `m3`, `c3` and `s2`.
  * in `eu-nl`: `i3_pro`, `i3`, `m4`, `s3` and `c4`.
  * in `eu-ch2`: `s3`.

* `tags` - (Optional, Map) Tags key/value pairs to associate with Host.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Specifies the Dedicated Host status.

* `available_vcpus` - The number of available vCPUs for the Dedicated Host.

* `available_memory` - The size of available memory for the Dedicated Host.

* `instance_total` - The number of the placed VMs.

* `instance_uuids` - The VMs started on the Dedicated Host.

* `host_type_name` -  The name of the Dedicated Host type.

* `vcpus` - The number of host vCPUs.

* `cores` -  The number of host physical cores.

* `sockets` -  The number of host physical sockets.

* `memory` - The size of host physical memory (MB).

* `available_instance_capacities` - The VM flavors placed on the Dedicated Host.

## Import

DeH can be imported using the `dedicated_host_id`, e.g.

```sh
terraform import opentelekomcloud_deh_host_v1.deh_host 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
