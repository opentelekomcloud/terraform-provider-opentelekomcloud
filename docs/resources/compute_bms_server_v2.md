---
subcategory: "Bare Metal Server (BMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_bms_server_v2"
sidebar_current: "docs-opentelekomcloud-resource-compute-bms-server-v2"
description: |-
  Manages a BMS Server resource within OpenTelekomCloud.
---

# opentelekomcloud_compute_bms_server_v2

Manages a BMS Server resource within OpenTelekomCloud.

## Example Usage

### Basic Instance

```hcl
variable "image_id" {}
variable "flavor_id" {}
variable "keypair_name" {}
variable "network_id" {}
variable "availability_zone" {}

resource "opentelekomcloud_compute_bms_server_v2" "basic" {
  name              = "basic"
  image_id          = var.image_id
  flavor_id         = var.flavor_id
  key_pair          = var.keypair_name
  security_groups   = ["default"]
  availability_zone = var.availability_zone
  metadata = {
    this = "that"
  }

  network {
    uuid = var.network_id
  }
}
```

### Instance Boot From Volume Image

```hcl
variable "flavor_id" {}
variable "keypair_name" {}
variable "network_id" {}
variable "availability_zone" {}

resource "opentelekomcloud_compute_bms_server_v2" "basic" {
  name              = "basic"
  flavor_id         = var.flavor_id
  key_pair          = var.keypair_name
  security_groups   = ["default"]
  availability_zone = var.availability_zone
  metadata = {
    this = "that"
  }

  network {
    uuid = var.network_id
  }

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    volume_type           = "SATA"
    volume_size           = 100
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    device_name           = "/dev/sda"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the BMS.

* `image_id` - (Optional; Required if `image_name` is empty.) Changing this creates a new bms server.

* `image_name` - (Optional; Required if `image_id` is empty.) The name of the
  desired image for the bms server. Changing this creates a new BMS server.

* `flavor_id` - (Optional; Required if `flavor_name` is empty) The flavor ID of
  the desired flavor for the BMS server. Changing this resizes the existing BMS server.

* `flavor_name` - (Optional; Required if `flavor_id` is empty) The name of the
  desired flavor for the BMS server. Changing this resizes the existing BMS server.

* `user_data` - (Optional) The user data to provide when launching the instance.
  Changing this creates a new BMS server.

* `security_groups` - (Optional) An array of one or more security group names
  to associate with the BMS server. Changing this results in adding/removing
  security groups from the existing BMS server.

* `availability_zone` - (Required) The availability zone in which to create
  the BMS server.

* `network` - (Optional) An array of one or more networks to attach to the
  BMS instance. Changing this creates a new BMS server.

* `metadata` - (Optional) Metadata key/value pairs to make available from
  within the instance. Changing this updates the existing BMS server metadata.

* `admin_pass` - (Optional) The administrative password to assign to the BMS server.
  Changing this changes the root password on the existing server.

* `key_pair` - (Optional) The name of a key pair to put on the BMS server. The key
  pair must already be created and associated with the tenant's account.
  Changing this creates a new BMS server.

* `stop_before_destroy` - (Optional) Whether to try stop instance gracefully
  before destroying it, thus giving chance for guest OS daemons to stop correctly.
  If instance doesn't stop within timeout, it will be destroyed anyway.

* `tags` - (Optional) Tags key/value pairs to associate with the instance.

The `network` block supports:

* `uuid` - (Required unless `port`  or `name` is provided) The network UUID to
  attach to the BMS server. Changing this creates a new BMS server.

* `name` - (Required unless `uuid` or `port` is provided) The human-readable
  name of the network. Changing this creates a new BMS server.

* `port` - (Required unless `uuid` or `name` is provided) The port UUID of a
  network to attach to the BMS server. Changing this creates a new server.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this
  network. Changing this creates a new BMS server.

* `fixed_ip_v6` - (Optional) Specifies a fixed IPv6 address to be used on this
  network. Changing this creates a new BMS server.

* `access_network` - (Optional) Specifies if this network should be used for
  provisioning access. Accepts true or false. Defaults to false.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the BMS server.

* `config_drive` - Whether to use the config_drive feature to configure the instance.

* `kernel_id` - The UUID of the kernel image when the AMI image is used.

* `user_id` - The ID of the user to which the BMS belongs.

* `host_status` - The nova-compute status: `UP`, `UNKNOWN`, `DOWN`, `MAINTENANCE` and `Null`.
