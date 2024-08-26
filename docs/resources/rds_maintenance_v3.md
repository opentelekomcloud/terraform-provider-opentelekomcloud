---
subcategory: "Relational Database Service (RDS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rds_maintenance_v3"
sidebar_current: "docs-opentelekomcloud-resource-rds-maintenance-v3"
description: |-
  Manages an RDS Maintenance windows resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RDS parameter group rule you can get at
[documentation portal](https://docs.otc.t-systems.com/relational-database-service/api-ref/api_v3_recommended/db_instance_management/configuring_the_maintenance_window.html)

# opentelekomcloud_rds_maintenance_v3

Manages a RDSv3 maintenance windows resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_rds_maintenance_v3" "test" {
  instance_id = var.instance_id
  start_time  = "12:00"
  end_time    = "16:00"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew, String) The ID of the RDS instance to which the maintenance window belongs.

-> **NOTE:** The interval between the `start_time` and `end_time` must be four hours.

* `start_time` - (Required, ForceNew, String) Specifies the start time.
  The value must be a valid value in the "HH:MM" format. The current time is in the UTC format.

* `end_time` - (Required, ForceNew, String) Specifies the end time.
  The value must be a valid value in the "HH:MM" format. The current time is in the UTC format.
