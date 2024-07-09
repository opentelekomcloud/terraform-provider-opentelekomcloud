---
subcategory: "Cloud Server Backup Service (CSBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_csbs_backup_policy_v1"
sidebar_current: "docs-opentelekomcloud-datasource-csbs-backup-policy-v1"
description: |-
Get details about backup Policy resources from OpenTelekomCloud
---

Up-to-date reference of API arguments for CSBS backup policy you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-server-backup-service/api-ref/api_description/backup_policy_management/querying_the_backup_policy_list.html#en-us-topic-0059304227)

# opentelekomcloud_csbs_backup_policy_v1

Use this data source to get details about backup Policy resources from OpenTelekomCloud.

## Example Usage
```hcl
variable "policy_id" {}

data "opentelekomcloud_csbs_backup_policy_v1" "csbs_policy" {
  id = var.policy_id
}
```

## Argument Reference
The following arguments are supported:

* `id` - (Optional) Specifies the ID of backup policy.

* `name` - (Optional) Specifies the backup policy name.

* `status` - (Optional) Specifies the backup policy status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `description` - Specifies the backup policy description.

* `provider_id` - Provides the Backup provider ID.

* `parameters` - Specifies the parameters of a backup policy.

The `scheduled_operation` block supports the following arguments:

* `name` - Specifies Scheduling period name.

* `description` - Specifies Scheduling period description.

* `enabled` - Specifies whether the scheduling period is enabled.

* `max_backups` - Specifies maximum number of backups that can be automatically created for a backup object.

* `retention_duration_days` - Specifies duration of retaining a backup, in days.

* `permanent` - Specifies whether backups are permanently retained.

* `trigger_pattern` - Specifies Scheduling policy of the scheduler.

* `operation_type` - Specifies Operation type, which can be backup.

* `id` - Specifies Scheduling period ID.

* `trigger_id` - Specifies Scheduler ID.

* `trigger_name` - Specifies Scheduler name.

* `trigger_type` - Specifies Scheduler type.

The `resource` block supports the following arguments:

* `id` - Specifies the ID of the object to be backed up.

* `type` - Entity object type of the backup object.

* `name` - Specifies backup object name.

The `tags` block supports the following arguments:

* `key` - Tag key. It cannot be an empty string.

* `value` - Tag value. It can be an empty string.
