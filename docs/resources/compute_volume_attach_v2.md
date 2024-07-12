---
subcategory: "Elastic Cloud Server (ECS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_volume_attach_v2"
sidebar_current: "docs-opentelekomcloud-resource-compute-volume-attach-v2"
description: |-
  Manages an ECS Disk Management resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for ECS disk management you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/openstack_nova_apis/disk_management)

# opentelekomcloud_compute_volume_attach_v2

Attaches a Block Storage Volume to an Instance using the OpenTelekomCloud
Compute (Nova) v2 API.

## Example Usage

### Basic attachment of a single volume to a single instance

```hcl
resource "opentelekomcloud_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "opentelekomcloud_compute_volume_attach_v2" "va_1" {
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volume_1.id
}
```

### Attaching multiple volumes to a single instance

```hcl
resource "opentelekomcloud_blockstorage_volume_v2" "volumes" {
  count = 2
  name  = format("vol-%02d", count.index + 1)
  size  = 1
}

resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  security_groups = ["default"]
}

resource "opentelekomcloud_compute_volume_attach_v2" "attachments" {
  count       = 2
  instance_id = opentelekomcloud_compute_instance_v2.instance_1.id
  volume_id   = opentelekomcloud_blockstorage_volume_v2.volumes[count.index].id
}

output "volume devices" {
  value = opentelekomcloud_compute_volume_attach_v2.attachments.*.device
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The ID of the Instance to attach the Volume to.

* `volume_id` - (Required) The ID of the Volume to attach to an Instance.

* `device` - (Optional) The device of the volume attachment (ex: `/dev/vdc`).

-> **Note:** Being able to specify a device is dependent upon the hypervisor in
  use. There is a chance that the device specified in Terraform will not be
  the same device the hypervisor chose. If this happens, Terraform will wish
  to update the device upon subsequent applying which will cause the volume
  to be detached and reattached indefinitely. Please use with caution.

## Attributes Reference

The following attributes are exported:

* `instance_id` - See Argument Reference above.

* `volume_id` - See Argument Reference above.

* `device` - See Argument Reference above.
-> **Note:** The correctness of this information is dependent upon the hypervisor in use.
  In some cases, this should not be used as an authoritative piece of information.

## Import

Volume Attachments can be imported using the Instance ID and Volume ID
separated by a slash, e.g.

```
$ terraform import opentelekomcloud_compute_volume_attach_v2.va_1 89c60255-9bd6-460c-822a-e2b959ede9d2/45670584-225f-46c3-b33e-6707b589b666
```
