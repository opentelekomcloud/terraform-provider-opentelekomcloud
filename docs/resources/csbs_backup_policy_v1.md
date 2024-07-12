---
subcategory: "Cloud Server Backup Service (CSBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_csbs_backup_policy_v1"
sidebar_current: "docs-opentelekomcloud-resource-csbs-backup-policy-v1"
description: |-
  Manages a CSBS Backup Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CSBS backup policy you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-server-backup-service/api-ref/api_description/backup_policy_management)

# opentelekomcloud_csbs_backup_policy_v1

Provides an OpenTelekomCloud Backup Policy of Resources.

## Example Usage

### Basic example

```hcl
variable "name" {}
variable "id" {}
variable "resource_name" {}

resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = var.name

  resource {
    id   = var.id
    type = "OS::Nova::Server"
    name = var.resource_name
  }

  scheduled_operation {
    enabled         = true
    operation_type  = "backup"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
  }
}
```

### Basic example with configured the week and month backups

```
variable "name" { }
variable "id" { }
variable "resource_name" { }
var "scheduled_operation_name" { }

resource "opentelekomcloud_csbs_backup_policy_v1" "backup_policy_v1" {
  name = var.name

  resource {
    id   = var.id
    type = "OS::Nova::Server"
    name = var.resource_name
  }
  scheduled_operation {
    name            = var.scheduled_operation_name
    enabled         = true
    operation_type  = "backup"
    max_backups     = "6"
    trigger_pattern = "BEGIN:VCALENDAR\r\nBEGIN:VEVENT\r\nRRULE:FREQ=WEEKLY;BYDAY=TH;BYHOUR=12;BYMINUTE=27\r\nEND:VEVENT\r\nEND:VCALENDAR\r\n"
    week_backups    = "4"
    month_backups   = "2"
    timezone        = "UTC+03:00"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of backup policy. The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional) Backup policy description. The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

* `provider_id` - (Required) Specifies backup provider ID. Default value is **fc4d5750-22e7-4798-8a46-f48f62c4c1da**

* `common` - (Optional) General backup policy parameters, which are blank by default.

The `scheduled_operation` block supports the following arguments:

* `name` - (Optional) Specifies Scheduling period name.The value consists of 1 to 255 characters and can contain only letters, digits, underscores (_), and hyphens (-).

* `description` - (Optional) Specifies Scheduling period description.The value consists of 0 to 255 characters and must not contain a greater-than sign (>) or less-than sign (<).

* `enabled` - (Optional) Specifies whether the scheduling period is enabled. Default value is **true**

* `max_backups` - (Optional) Specifies maximum number of backups that can be automatically created for a backup object.

* `retention_duration_days` - (Optional) Specifies duration of retaining a backup, in days.

-> **Note:** If `day_backups`, `week_backups`, `month_backups` or `year_backups` is configured
  `timezone` is mandatory.

* `day_backups` - (Optional) Specifies the maximum number of retained daily backups.
  The latest backup of each day is saved in the long term. This parameter can be effective
  together with the maximum number of retained backups specified by `max_backups`.

* `week_backups` - (Optional) Specifies the maximum number of retained weekly backups.
  The latest backup of each week is saved in the long term. This parameter can be effective
  together with the maximum number of retained backups specified by `max_backups`.

* `month_backups` - (Optional) Specifies the maximum number of retained monthly backups.
  The latest backup of each month is saved in the long term. This parameter can be effective
  together with the maximum number of retained backups specified by `max_backups`.

* `year_backups` - (Optional) Specifies the maximum number of retained yearly backups.
  The latest backup of each year is saved in the long term. This parameter can be effective
  together with the maximum number of retained backups specified by `max_backups`.

* `timezone` - (Optional) Time zone where the user is located, for example, `UTC+08:00`.

* `permanent` - (Optional) Specifies whether backups are permanently retained.

* `trigger_pattern` - (Required) Specifies Scheduling policy of the scheduler.

* `operation_type` - (Required) Specifies Operation type, which can be backup.

The `resource` block supports the following arguments:

* `id` - (Required) Specifies the ID of the object to be backed up.

* `type` - (Required) Entity object type of the backup object. If the type is VMs, the value is **OS::Nova::Server**.

* `name` - (Required) Specifies backup object name.

The `tags` block supports the following arguments:

* `key` - (Required) Tag key. It cannot be an empty string.

* `value` - (Required) Tag value. It can be an empty string.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - Status of Backup Policy.

* `id` - Backup Policy ID.

* `created_at` - Backup creation time.

* `scheduled_operation` - Backup plan information

  * `id` -  Specifies Scheduling period ID.

  * `trigger_id` - Specifies Scheduler ID.

  * `trigger_name` - Specifies Scheduler name.

  * `trigger_type` - Specifies Scheduler type.


## Import

Backup Policy can be imported using `id`, e.g.

```sh
terraform import opentelekomcloud_csbs_backup_policy_v1.backup_policy_v1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```
