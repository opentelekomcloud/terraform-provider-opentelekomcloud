---
subcategory: "Elastic Cloud Server (ECS)"
---

# opentelekomcloud_compute_instance_v2

Manages a V2 VM instance resource within OpenTelekomCloud.

## Example Usage

### Basic Instance

```hcl
variable image_id {}

resource "opentelekomcloud_compute_instance_v2" "basic" {
  name            = "basic"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }

  metadata = {
    this = "that"
  }

  tags = {
    muh = "kuh"
  }
}
```

### Instance With Attached Volume

```hcl
variable image_id {}

resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 4
}

resource "opentelekomcloud_compute_instance_v2" "myinstance" {
  name            = "myinstance"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_network"
  }
}

resource "opentelekomcloud_compute_volume_attach_v2" "attached" {
  instance_id = opentelekomcloud_compute_instance_v2.myinstance.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.myvol.id
}
```

### Boot From Volume

```hcl
variable image_id {}

resource "opentelekomcloud_compute_instance_v2" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    volume_size           = 5
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    volume_type           = "SSD"
  }

  network {
    name = "my_network"
  }
}
```

### Boot From an Existing Volume

```hcl
variable image_id {}

resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = var.image_id
}

resource "opentelekomcloud_compute_instance_v2" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = opentelekomcloud_blockstorage_volume_v2.myvol.id
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    name = "my_network"
  }
}
```

### Boot Instance, Create Volume, and Attach Volume as a Block Device

```hcl
variable image_id {}
variable data_image_id {}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = var.data_image_id
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 1
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Boot Instance and Attach Existing Volume as a Block Device

```hcl
variable image_id {}
variable data_image_id {}

resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = var.data_image_id
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    uuid                  = opentelekomcloud_blockstorage_volume_v2.volume_1.id
    source_type           = "volume"
    destination_type      = "volume"
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Instance With Multiple Networks

```hcl
variable image_id {}

resource "opentelekomcloud_networking_floatingip_v2" "myip" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_compute_instance_v2" "multi-net" {
  name            = "multi-net"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    name = "my_first_network"
  }

  network {
    name = "my_second_network"
  }
}

resource "opentelekomcloud_compute_floatingip_associate_v2" "myip" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.myip.address
  instance_id = opentelekomcloud_compute_instance_v2.multi-net.id
  fixed_ip    = opentelekomcloud_compute_instance_v2.multi-net.network.1.fixed_ip_v4
}
```

### Instance with Multiple Ephemeral Disks

```hcl
variable image_id {}
variable data_image_id {}

resource "opentelekomcloud_compute_instance_v2" "multi-eph" {
  name            = "multi_eph"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "image"
    uuid                  = var.data_image_id
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instance with User Data (cloud-init)

```hcl
variable image_id {}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "basic"
  image_id        = var.image_id
  flavor_id       = "s2.large.4"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  user_data       = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    name = "my_network"
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource.

* `image_id` - (Optional; Required if `image_name` is empty and not booting from a volume. Do not specify if booting
  from a volume.) The image ID of the desired image for the server. Changing this creates a new server.

* `image_name` - (Optional; Required if `image_id` is empty and not booting from a volume. Do not specify if booting
  from a volume.) The name of the desired image for the server. Changing this creates a new server.

* `flavor_id` - (Optional; Required if `flavor_name` is empty) The flavor ID of the desired flavor for the server.
  Changing this resizes the existing server.

* `flavor_name` - (Optional; Required if `flavor_id` is empty) The name of the desired flavor for the server. Changing
  this resizes the existing server.

* `user_data` - (Optional) The user data to provide when launching the instance. Changing this creates a new server.

* `ssh_private_key_path` - (Optional) The path to the private key to use for SSH access. Required only if you want to
  get the password from the windows instance.

* `security_groups` - (Optional) An array of one or more security group names to associate with the server. Changing
  this results in adding/removing security groups from the existing server.

~> **Warning** Names should be used and not IDs. Security group names should be **unique**, otherwise it will return an
error.

-> When attaching the instance to networks using Ports, place the security groups on the Port and not the instance.

* `availability_zone` - (Optional) The availability zone in which to create the server. Changing this creates a new
  server.

* `network` - (Optional) An array of one or more networks to attach to the instance. Required when there are multiple
  networks defined for the tenant. The network object structure is documented below. Changing this creates a new server.

* `metadata` - (Optional) Metadata key/value pairs to make available from within the instance. Changing this updates the
  existing server metadata.

* `config_drive` - (Optional) Whether to use the config_drive feature to configure the instance. Changing this creates a
  new server.

* `admin_pass` - (Optional) The administrative password to assign to the server. Changing this changes the root password
  on the existing server.

* `key_pair` - (Optional) The name of a key pair to put on the server. The key pair must already be created and
  associated with the tenant's account. Changing this creates a new server.

* `block_device` - (Optional) Configuration of block devices. The block_device structure is documented below. Changing
  this creates a new server. You can specify multiple block devices which will create an instance with multiple disks.
  This configuration is very flexible, so please see the
  following [reference](http://docs.openstack.org/developer/nova/block_device_mapping.html)
  for more information.

* `scheduler_hints` - (Optional) Provide the Nova scheduler with hints on how the instance should be launched. The
  available hints are described below.

* `tags` -  (Optional) Tags key/value pairs to associate with the instance.

* `stop_before_destroy` - (Optional) Whether to try stop instance gracefully before destroying it, thus giving chance
  for guest OS daemons to stop correctly. If instance doesn't stop within a timeout, it will be destroyed anyway.

* `force_delete` - (Optional) Whether to force the OpenTelekomCloud instance to be forcefully deleted. This is useful
  for environments that have reclaim/soft deletion enabled.

* `auto_recovery` - (Optional) Configures or deletes automatic recovery of an instance. Defaults to true.

* `power_state` - (Optional) Provide the VM state. Only `active` and `shutoff` are supported values.

  ->
  If the initial `power_state` is the `shutoff` the VM will be stopped immediately after build, and the provisioners
  like remote-exec or files are not supported.

The `network` block supports:

* `uuid` - (Required unless `port`  or `name` is provided) The network UUID to attach to the server. Changing this
  creates a new server.

* `name` - (Required unless `uuid` or `port` is provided) The human-readable name of the network. Changing this creates
  a new server.

* `port` - (Required unless `uuid` or `name` is provided) The port UUID of a network to attach to the server. Changing
  this creates a new server.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this network. Changing this creates a new
  server.

* `fixed_ip_v6` - (Optional) Specifies a fixed IPv6 address to be used on this network. Changing this creates a new
  server.

* `access_network` - (Optional) Specifies if this network should be used for provisioning access. Accepts true or false.
  Defaults to false.

The `block_device` block supports:

* `uuid` - (Required unless `source_type` is set to `"blank"` ) The UUID of the image, volume, or snapshot. Changing
  this creates a new server.

* `source_type` - (Required) The source type of the device. Must be one of
  "blank", "image", "volume", or "snapshot". Changing this creates a new server.

* `volume_size` - The size of the volume to create (in gigabytes). Required in the following combinations: source=image
  and destination=volume, and source=blank and destination=volume. Changing this creates a new server.

* `volume_type` - (Optional) Currently, the value can be `SSD` (ultra-I/O disk type),
  `SAS` (high I/O disk type), or `SATA` (common I/O disk type)
  [OTC-API](https://docs.otc.t-systems.com/en-us/api/ecs/en-us_topic_0065817708.html)

* `boot_index` - (Optional) The boot index of the volume. It defaults to 0. Changing this creates a new server.

* `destination_type` - (Optional) The type that gets created. Currently only support "volume". Changing this creates a
  new server.

* `delete_on_termination` - (Optional) Delete the volume / block device upon termination of the instance. Defaults to
  false. Changing this creates a new server.

The `scheduler_hints` block supports:

* `group` - (Optional) A UUID of a Server Group. The instance will be placed into that group.

* `different_host` - (Optional) A list of instance UUIDs. The instance will be scheduled on a different host than all
  other instances.

* `same_host` - (Optional) A list of instance UUIDs. The instance will be scheduled on the same host of those specified.

* `query` - (Optional) A conditional query that a compute node must pass in order to host an instance.

* `target_cell` - (Optional) The name of a cell to host the instance.

* `build_near_host_ip` - (Optional) An IP Address in CIDR form. The instance will be placed on a compute node that is in
  the same subnet.

* `tenancy` - (Optional) The tenancy specifies whether the ECS is to be created on a Dedicated Host
  (DeH) or in a shared pool.

* `deh_id` - (Optional) The ID of DeH. This parameter takes effect only when the value of tenancy is dedicated.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `access_ip_v4` - The first detected Fixed IPv4 address _or_ the Floating IP.

* `access_ip_v6` - The first detected Fixed IPv6 address.

* `metadata` - See Argument Reference above.

* `security_groups` - See Argument Reference above.

* `flavor_id` - See Argument Reference above.

* `flavor_name` - See Argument Reference above.

* `network/uuid` - See Argument Reference above.

* `network/name` - See Argument Reference above.

* `network/port` - See Argument Reference above.

* `network/fixed_ip_v4` - The Fixed IPv4 address of the Instance on that network.

* `network/fixed_ip_v6` - The Fixed IPv6 address of the Instance on that network.

* `network/mac` - The MAC address of the NIC on that network.

* `volume_attached/id` - The volume id on that attachment.

* `all_metadata` - Contains all instance metadata, even metadata not set by Terraform.

* `auto_recovery` - See Argument Reference above.

## Notes

### Multiple Ephemeral Disks

It's possible to specify multiple `block_device` entries to create an instance with multiple ephemeral (local) disks. In
order to create multiple ephemeral disks, the sum of the total amount of ephemeral space must be less than or equal to
what the chosen flavor supports.

The following example shows how to create an instance with multiple ephemeral disks:

```hcl
resource "opentelekomcloud_compute_instance_v2" "foo" {
  name = "terraform-test"

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "image"
    uuid                  = var.image_id
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "volume"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instances and Ports

Neutron Ports are a great feature and provide a lot of functionality. However, there are some notes to be aware of when
mixing Instances and Ports:

* When attaching an Instance to one or more networks using Ports, place the security groups on the Port and not the
  Instance. If you place the security groups on the Instance, the security groups will not be applied upon creation, but
  they will be applied upon a refresh. This is a known OpenTelekomCloud bug.

* Network IP information is not available within an instance for networks that are attached with Ports. This is mostly
  due to the flexibility Neutron Ports provide when it comes to IP addresses. For example, a Neutron Port can have
  multiple Fixed IP addresses associated with it. It's not possible to know which single IP address the user would want
  returned to the Instance's state information. Therefore, in order for a Provisioner to connect to an Instance via it's
  network Port, customize the `connection` information:

```hcl
resource "opentelekomcloud_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"

  network_id = "0a1d0a27-cffa-4de3-92c5-9d3fd3f2e74d"

  security_group_ids = [
    "2f02d20a-8dca-49b7-b26f-b6ce9fddaf4f",
    "ca1e5ed7-dae8-4605-987b-fadaeeb30461",
  ]
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network {
    port = opentelekomcloud_networking_port_v2.port_1.id
  }

  connection {
    user        = "root"
    host        = opentelekomcloud_networking_port_v2.port_1.fixed_ip.0.ip_address
    private_key = "~/path/to/key"
  }

  provisioner "remote-exec" {
    inline = [
      "echo terraform executed > /tmp/foo",
    ]
  }
}
```

## Importing instances

Importing instances can be tricky, since the nova api does not offer all information provided at creation time for later
retrieval. Network interface attachment order, and number and sizes of ephemeral disks are examples of this.

### Importing basic instance

Assume you want to import an instance with one ephemeral root disk, and one network interface.

Your configuration would look like the following:

```hcl
resource "opentelekomcloud_compute_instance_v2" "basic_instance" {
  name      = "basic"
  flavor_id = var.flavor_id
  key_pair  = var.key_pair
  image_id  = var.image_id

  network {
    name = var.network_name
  }
}

```

Then you execute

```
terraform import opentelekomcloud_compute_instance_v2.basic_instance <instance_id>
```

### Importing instance with multiple network interfaces.

Nova returns the network interfaces grouped by network, thus not in creation order. That means that if you have multiple
network interfaces you must take care of the order of networks in your configuration.

As example, we want to import an instance with one ephemeral root disk, and 3 network interfaces.

Examples

```hcl
resource "opentelekomcloud_compute_instance_v2" "boot-from-volume" {
  name      = "boot-from-volume"
  flavor_id = var.flavor_id
  key_pair  = var.key_pair
  image_id  = var.image_id

  network {
    name = var.network_1_name
  }
  network {
    name = var.network_2_name
  }
  network {
    name        = var.network_1_name
    fixed_ip_v4 = var.fixed_ip_v4
  }

}
```

In the above configuration the networks are out of order compared to what nova and thus the import code returns, which
means the plan will not be empty after import.

So either with care check the plan and modify configuration, or read the network order in the state file after import
and modify your configuration accordingly.

* A note on ports. If you have created a neutron port independent of an instance, then the import code has no way to
  detect that the port is created idenpendently, and therefore on deletion of imported instances you might have port
  resources in your project, which you expected to be created by the instance and thus to also be deleted with the
  instance.

### Importing instances with multiple block storage volumes.

We have an instance with two block storage volumes, one bootable and one non-bootable. Note that we only configure the
bootable device as block_device. The other volumes can be specified as `opentelekomcloud_blockstorage_volume_v2`

```hcl
resource "opentelekomcloud_compute_instance_v2" "instance_2" {
  name      = "instance_2"
  image_id  = var.image_id
  flavor_id = var.flavor_id
  key_pair  = var.key_pair

  block_device {
    uuid                  = var.image_id
    source_type           = "image"
    destination_type      = "volume"
    boot_index            = 0
    delete_on_termination = true
  }

  network {
    name = var.network_name
  }
}
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  size = 1
  name = var.volume_name
}
resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
  instance_id = opentelekomcloud_compute_instance_v2.instance_2.id
}
```

To import the instance outlined in the above configuration do the following:

```
terraform import opentelekomcloud_compute_instance_v2.instance_2 <instance_id>
import opentelekomcloud_blockstorage_volume_v2.volume_1 <volume_id>
terraform import opentelekomcloud_compute_volume_attach_v2.va_1
<instance_id>/<volume_id>
```

* A note on block storage volumes, the importer does not read delete_on_termination flag, and always assumes true. If
  you import an instance created with delete_on_termination false, you end up with "orphaned" volumes after destruction
  of instances.
