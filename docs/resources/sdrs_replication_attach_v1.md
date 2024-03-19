---
subcategory: "Storage Disaster Recovery Service (SDRS)"
---

Up-to-date reference of API arguments for SDRS replication pair attachment you can get at
`https://docs.otc.t-systems.com/storage-disaster-recovery-service/api-ref/sdrs_apis/protected_instance/index.html`.

# opentelekomcloud_sdrs_replication_attach_v1

Manages a SDRS replication pair attachment resource within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_sdrs_domain_v1" "domain_1" {}

resource "opentelekomcloud_sdrs_protectiongroup_v1" "group_1" {
  name                     = "group_1"
  source_availability_zone = "eu-de-02"
  target_availability_zone = "eu-de-01"
  domain_id                = data.opentelekomcloud_sdrs_domain_v1.domain_1.id
  source_vpc_id            = var.vpc_id
  dr_type                  = "migration"
}

resource "opentelekomcloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = var.image_id
  flavor   = "s3.medium.1"
  vpc_id   = var.vpc_id

  nics {
    network_id = var.network_id
  }

  availability_zone = "eu-de-02"
}

resource "opentelekomcloud_sdrs_protected_instance_v1" "instance_1" {
  name                 = "instance_create"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  server_id            = opentelekomcloud_ecs_instance_v1.instance_1.id
  delete_target_server = true
}

resource "opentelekomcloud_evs_volume_v3" "volume_1" {
  name              = "volume_1"
  description       = "first test volume"
  availability_zone = "eu-de-02"
  volume_type       = "SATA"
  size              = 12
}

resource "opentelekomcloud_sdrs_replication_pair_v1" "pair_1" {
  name                 = "replication_1"
  group_id             = opentelekomcloud_sdrs_protectiongroup_v1.group_1.id
  volume_id            = opentelekomcloud_evs_volume_v3.volume_1.id
  delete_target_volume = true
}

resource "opentelekomcloud_sdrs_replication_attach_v1" "attach_1" {
  instance_id    = opentelekomcloud_sdrs_protected_instance_v1.instance_1.id
  replication_id = opentelekomcloud_sdrs_replication_pair_v1.pair_1.id
  device         = "/dev/vdb"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a protected instance.

* `replication_id` - (Required, String, ForceNew) Specifies the ID of a replication pair.

* `device` - (Required, String, ForceNew) Specifies the disk device name of a replication pair. There are several
  restrictions on this field as follows：

    + The new disk device name cannot be the same as an existing one.

    + Set the parameter value to /dev/sda for the system disks of protected instances created using Xen servers and to
      /dev/sdx for data disks, where x is a letter in alphabetical order. For example, if there are two data disks, set the
      device names of the two data disks to /dev/sdb and /dev/sdc, respectively. If you set a device name starting with
      /dev/vd, the system uses /dev/sd by default.

    + Set the parameter value to /dev/vda for the system disks of protected instances created using KVM servers and
      to /dev/vdx for data disks, where x is a letter in alphabetical order. For example, if there are two data disks,
      set the device names of the two data disks to /dev/vdb and /dev/vdc, respectively. If you set a device name starting
      with /dev/sd, the system uses /dev/vd by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the SDRS protected instance.

* `region` - The attachment region.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The SDRS replication attach can be imported using the `protected_instance_id` and `replication_id`, separated
by a slash , e.g.

```bash
$ terraform import opentelekomcloud_sdrs_replication_attach_v1.test <protected_instance_id>/<replication_id>
```
