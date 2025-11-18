---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_instance_minor_version_upgrade_v3"
sidebar_current: "docs-opentelekomcloud-resource-rds-instance-minor-version-upgrade-v3"
description: |-
  Manages an RDS instance minor version upgrade resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RDS instance minor version upgrade you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/db_instance_management/upgrading_a_minor_version.html#rds-05-0024)

# opentelekomcloud_rds_instance_minor_version_upgrade_v3

Manages an RDSv3 instance minor version upgrade resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_rds_instance_minor_version_upgrade_v3" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew, String) Specifies the ID of the RDS instance to upgrade.

* `delay` - (Optional, ForceNew, Bool) Specifies whether the upgrade is delayed to the maintenance window.
  Defaults to `false`.
    + **true**: The upgrade is delayed and performed within the maintenance window.
    + **false**: The upgrade is performed immediately.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the instance ID.
