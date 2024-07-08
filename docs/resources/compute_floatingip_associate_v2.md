---
subcategory: "Elastic Cloud Server (ECS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_floatingip_associate_v2"
sidebar_current: "docs-opentelekomcloud-resource-compute-floatingip-associate-v2"
description: |-
Manages an EIP Associate resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for EIP you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-ip/api-ref/native_openstack_neutron_apis_v2.0/floating_ip_address)

# opentelekomcloud_compute_floatingip_associate_v2

Associate a floating IP to an instance. This can be used instead of the
`floating_ip` options in `opentelekomcloud_compute_instance_v2`.

~>
Floating IP compute APIs are marked as discarded in [help center](https://docs.otc.t-systems.com/en-us/api/ecs/en-us_topic_0065817682.html).
Please use [`resource/opentelekomcloud_networking_floatingip_associate_v2`](networking_floatingip_associate_v2.md).

## Example Usage

### Automatically detect the correct network

```hcl
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
}
```

### Explicitly set the network to attach to

```hcl
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = 3
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }

  network {
    name = "default"
  }
}

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.fip_1.address
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  fixed_ip    = opentelekomcloud_compute_instance_v2.instance_1.network.1.fixed_ip_v4
}
```

## Argument Reference

The following arguments are supported:

* `floating_ip` - (Required) The floating IP to associate.

* `instance_id` - (Required) The instance to associte the floating IP with.

* `fixed_ip` - (Optional) The specific IP address to direct traffic to.

## Attributes Reference

The following attributes are exported:

* `floating_ip` - See Argument Reference above.

* `instance_id` - See Argument Reference above.

* `fixed_ip` - See Argument Reference above.

## Import

This resource can be imported by specifying all three arguments, separated
by a forward slash:

```sh
terraform import opentelekomcloud_compute_floatingip_associate_v2.fip_1 <floating_ip>/<instance_id>/<fixed_ip>
```
