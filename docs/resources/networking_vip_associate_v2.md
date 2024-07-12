---
subcategory: "Virtual Private Cloud (VPC)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_networking_vip_associate_v2"
sidebar_current: "docs-opentelekomcloud-resource-networking-vip-associate-v2"
description: |-
  Manages a VPC VIP Associate resource within OpenTelekomCloud.
---

# opentelekomcloud_networking_vip_associate_v2

Manages a V2 vip associate resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "opentelekomcloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = opentelekomcloud_networking_network_v2.network_1.id
}

resource "opentelekomcloud_networking_router_interface_v2" "router_interface_1" {
  router_id = opentelekomcloud_networking_router_v2.router_1.id
  subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_router_v2" "router_1" {
  name             = "router_1"
  external_gateway = "0a2228f2-7f8a-45f1-8e09-9039e1d09975"
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = opentelekomcloud_networking_network_v2.network_1.id
  fixed_ip {
    subnet_id = opentelekomcloud_networking_subnet_v2.subnet_1.id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_2" {
  name            = "instance_2"
  security_groups = ["default"]

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
}

resource "opentelekomcloud_networking_vip_v2" "vip_1" {
  network_id = opentelekomcloud_networking_network_v2.network_1.id
  subnet_id  = opentelekomcloud_networking_subnet_v2.subnet_1.id
}

resource "opentelekomcloud_networking_vip_associate_v2" "vip_associate_1" {
  vip_id = opentelekomcloud_networking_vip_v2.vip_1.id
  port_ids = [
    opentelekomcloud_networking_port_v2.port_1.id,
    opentelekomcloud_networking_port_v2.port_2.id
  ]
}
```

## Example VIP-EIP association

```hcl
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "my-shared-subnet"
}

data "opentelekomcloud_images_image_v2" "latest_image" {
  name        = "Standard_Debian_11_latest"
  most_recent = true
}

resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  fixed_ip {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_vip_ass_1"
  security_groups = ["default"]
  image_id        = data.opentelekomcloud_images_image_v2.latest_image.id
  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }
}

resource "opentelekomcloud_networking_port_v2" "port_2" {
  name           = "port_2"
  admin_state_up = "true"
  network_id     = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  fixed_ip {
    subnet_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
  }
}

resource "opentelekomcloud_compute_instance_v2" "instance_2" {
  name            = "instance_vip_ass_2"
  security_groups = ["default"]
  image_id        = data.opentelekomcloud_images_image_v2.latest_image.id
  network {
    port = opentelekomcloud_networking_port_v2.port_2.id
  }
}

resource "opentelekomcloud_vpc_eip_v1" "vip_eip_1" {
  publicip {
    type = "5_bgp"
    name = "eip-vip"
  }
  bandwidth {
    name        = "eip-bandwidth-vip"
    size        = 10
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_networking_vip_v2" "vip_1" {
  network_id = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  subnet_id  = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.subnet_id
}

resource "opentelekomcloud_networking_vip_associate_v2" "vip_associate_1" {
  vip_id = opentelekomcloud_networking_vip_v2.vip_1.id
  port_ids = [
    opentelekomcloud_networking_port_v2.port_1.id,
    opentelekomcloud_networking_port_v2.port_2.id,
  ]
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "vip_eip_associate_1" {
  floating_ip = opentelekomcloud_vpc_eip_v1.vip_eip_1.publicip.0.ip_address
  port_id     = opentelekomcloud_networking_vip_v2.vip_1.id
}
```

## Argument Reference

The following arguments are supported:

* `vip_id` - (Required) The ID of vip to attach the port to.
  Changing this creates a new vip associate.

* `port_ids` - (Required) An array of one or more IDs of the ports to attach the vip to.
  Changing this creates a new vip associate.

## Attributes Reference

The following attributes are exported:

* `vip_id` - See Argument Reference above.

* `port_ids` - See Argument Reference above.

* `vip_subnet_id` - The ID of the subnet this vip connects to.

* `vip_ip_address` - The IP address in the subnet for this vip.
