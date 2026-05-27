---
subcategory: "Elastic Volume Service (EVS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_evs_snapshot_v2"
sidebar_current: "docs-opentelekomcloud-resource-evs-snapshot-v2"
description: |-
Manages an EVS snapshot resource within OpenTelekomCloud.
---

# opentelekomcloud_evs_snapshot_v2

Provides an EVS snapshot resource.

## Example Usage

```hcl
data "opentelekomcloud_compute_availability_zones_v2" "test" {}

resource "opentelekomcloud_evs_volume_v3" "test" {
  name              = "volume-001"
  description       = "Created by acc test"
  availability_zone = data.opentelekomcloud_compute_availability_zones_v2.test.names[0]
  volume_type       = "SAS"
  size              = 12
}

resource "opentelekomcloud_evs_snapshot_v2" "test" {
  name        = "snapshot-001"
  description = "Daily backup"
  volume_id   = opentelekomcloud_evs_volume_v3.test.id

  metadata = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - (Required, String, ForceNew) The ID of the snapshot's source disk. Changing this parameter creates a new
  snapshot.

* `name` - (Required, String) The name of the snapshot. The value can contain a maximum of 255 bytes.

* `metadata` - (Optional, Map, ForceNew) Specifies the user-defined metadata key-value pairs. Changing this parameter
  creates a new snapshot.

* `description` - (Optional, String) The description of the snapshot. The value can contain a maximum of 255 bytes.

* `force` - (Optional, Boolean) Specifies whether to forcibly create a snapshot. Defaults to `false`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the snapshot.

* `status` - The status of the snapshot.

* `region` - The region in which the EVS snapshot resource was created.

* `size` - The size of the snapshot in GB.

* `created_at` - The time when the snapshot was created.

* `updated_at` - The time when the snapshot was updated.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 3 minutes.

## Import

EVS snapshot can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_evs_snapshot_v2.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `force`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the snapshot. Also, you can
ignore changes as below.

```hcl
resource "opentelekomcloud_evs_snapshot_v2" "test" {
  # ...

  lifecycle {
    ignore_changes = [
      force,
    ]
  }
}
```
