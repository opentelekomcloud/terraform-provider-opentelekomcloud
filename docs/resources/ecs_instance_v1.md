---
subcategory: "Elastic Cloud Server (ECS)"
---

Up-to-date reference of API arguments for ECS you can get at
`https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/apis_recommended/lifecycle_management`.

# opentelekomcloud_ecs_instance_v1

Manages a V1 ECS instance resource within OpenTelekomCloud.

## Example Usage

### Basic Instance

```hcl
resource "opentelekomcloud_ecs_instance_v1" "basic" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  availability_zone = "eu-de-01"
  key_name          = "KeyPair-test"

  tags = {
    muh = "kuh"
  }
}
```

### Basic Instance with security group

```hcl
resource "opentelekomcloud_compute_secgroup_v2" "secgroup_1" {
  name = "ecs_secgroup"
}

resource "opentelekomcloud_ecs_instance_v1" "basic" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  security_groups   = [opentelekomcloud_compute_secgroup_v2.secgroup_1.id]
  availability_zone = "eu-de-01"
  key_name          = "KeyPair-test"
}
```

### Instance with Data Disks

```hcl
resource "opentelekomcloud_ecs_instance_v1" "basic" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  system_disk_type = "SAS"
  system_disk_size = 40

  data_disks {
    type = "SATA"
    size = "10"
  }
  data_disks {
    type = "SAS"
    size = "20"
  }

  delete_disks_on_termination = true
  availability_zone           = "eu-de-01"
  key_name                    = "KeyPair-test"
}
```

### Instance With Attached Volume

```hcl
resource "opentelekomcloud_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 1
}

resource "opentelekomcloud_ecs_instance_v1" "basic" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  availability_zone = "eu-de-01"
  key_name          = "KeyPair-test"
}

resource "opentelekomcloud_compute_volume_attach_v2" "attached" {
  instance_id = opentelekomcloud_ecs_instance_v1.basic.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.myvol.id
}
```

### Instance With Multiple Networks

```hcl
resource "opentelekomcloud_networking_floatingip_v2" "this" {
  pool = "admin_external_net"
}

resource "opentelekomcloud_ecs_instance_v1" "this" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  nics {
    network_id = "2c0a74a9-4395-4e62-a17b-e3e86fbf66b7"
  }

  availability_zone = "eu-de-01"
  key_name          = "KeyPair-test"
}

resource "opentelekomcloud_networking_floatingip_associate_v2" "this" {
  floating_ip = opentelekomcloud_networking_floatingip_v2.this.address
  port_id     = opentelekomcloud_ecs_instance_v1.this.nics.0.port_id
}
```

### Instance with User Data (cloud-init)

```hcl
resource "opentelekomcloud_ecs_instance_v1" "basic" {
  name     = "server_1"
  image_id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor   = "s2.large.2"
  vpc_id   = "8eed4fc7-e5e5-44a2-b5f2-23b3e5d46235"

  nics {
    network_id = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }

  user_data = "#cloud-config\nhostname: server_1.example.com\nfqdn: server_1.example.com"
  key_name  = "KeyPair-test"
}
```

-> `user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

### Instance with encrypted disks

```hcl
resource opentelekomcloud_ecs_instance_v1 ecs {
  name              = var.host_name
  flavor            = var.flavor_name
  availability_zone = var.az
  security_groups   = [data.opentelekomcloud_networking_secgroup_v2.default.id]
  vpc_id            = var.vpc_id
  image_id          = var.image_id
  auto_recovery     = true

  nics {
    network_id = var.vpc_subnetwork_id
    ip_address = var.private_ip
  }

  system_disk_type            = var.disk_type
  system_disk_size            = var.disk_size
  system_disk_kms_id          = var.key_disk_encryption
  delete_disks_on_termination = true

  data_disks {
    type   = "SSD"
    size   = 40
    kms_id = var.key_disk_encryption
  }
}
```

~>
  Encrypted disks requires EVS to be authorized to use KMS keys. The easiest way is to create an encrypted
  instance  via the console - this should be done only once per project. Another way is to use an agency,
  same way it's [done for CCE](cce_cluster_v3.md#creating-agency).

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the instance.

* `image_id` - (Required) The ID of the desired image for the server. Changing this creates a new server.

* `flavor` - (Required) The name of the desired flavor for the server.

* `user_data` - (Optional) The user data to provide when launching the instance.
  Changing this creates a new server.

* `password` - (Optional) The administrative password to assign to the server.
  Changing this creates a new server.

* `key_name` - (Optional) The name of a key pair to put on the server. The key
  pair must already be created and associated with the tenant's account.
  Changing this creates a new server.

* `vpc_id` - (Required) The ID of the desired VPC for the server. Changing this creates a new server.

* `nics` - (Required) An array of one or more networks to attach to the
  instance. The nics object structure is documented below. Changing this
  creates a new server.

* `system_disk_type` - (Optional) The system disk type of the server. For HANA, HL1, and HL2 ECSs use `co-p1` and `uh-l1` disks.
  Changing this creates a new server. Options are limited depending on AZ. Available options are:
  * `SATA`: common I/O disk type. Available for all AZs.
  * `SAS`: high I/O disk type. Available for all AZs.
  * `SSD`: ultra-high I/O disk type. Available for all AZs.
  * `co-p1`: high I/O(performance-optimized) disk type.
  * `uh-l1`: ultra-high I/O(latency-optimized) disk type.

* `system_disk_size` - (Optional) The system disk size in GB, The value range is 1 to 1024.
  Changing this creates a new server.

* `data_disks` - (Optional) An array of one or more data disks to attach to the
  instance. The `data_disks` object structure is documented below. Changing this
  creates a new server.

* `system_disk_kms_id` - (Optional) The Encryption KMS ID of the system disk. Changing this
  creates a new server.

* `security_groups` - (Optional) An array of one or more security group IDs
  to associate with the server. If this parameter is left blank, the `default`
  security group is bound to the ECS by default.

* `availability_zone` - (Required) The availability zone in which to create the server.
  Changing this creates a new server.

* `auto_recovery` - (Optional) Whether configure automatic recovery of an instance.

* `delete_disks_on_termination` - (Optional) Delete the data disks upon termination of the instance.
  Defaults to false. Changing this creates a new server.

* `tags` - (Optional) Tags key/value pairs to associate with the instance.

The `nics` block supports:

* `network_id` - (Required) The network UUID to attach to the server. Changing this creates a new server.

* `ip_address` - (Optional) Specifies a fixed IPv4 address to be used on this
  network. Changing this creates a new server.

The `data_disks` block supports:

* `type` - (Required) The data disk type of the server. For HANA, HL1, and HL2 ECSs use `co-p1` and `uh-l1` disks.
  Changing this creates a new server. Options are limited depending on AZ. Available options are:
  * `SATA`: common I/O disk type. Available for all AZs.
  * `SAS`: high I/O disk type. Available for all AZs.
  * `SSD`: ultra-high I/O disk type. Available for all AZs.
  * `co-p1`: high I/O(performance-optimized) disk type.
  * `uh-l1`: ultra-high I/O(latency-optimized) disk type.

* `size` - (Required) The size of the data disk in GB. The value range is 10 to 32768.
  Changing this creates a new server.

* `kms_id` - (Optional) The Encryption KMS ID of the data disk. Changing this
  creates a new server.

* `snapshot_id` - (Optional) Specifies the snapshot ID or ID of the original data disk contained in the full-ECS image.
  Changing this creates a new server.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `system_disk_id` - The ID of the system disk.

* `nics/mac_address` - The MAC address of the NIC on that network.

* `nics/type` - The type of the address of the NIC on that network.

* `nics/port_id` - The port ID of the NIC on that network.

* `volumes_attached/id` - The ID of the data disk.

## Import

Instances can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_ecs_instance_v1.instance_1 d90ce693-5ccf-4136-a0ed-152ce412b6b9
```
