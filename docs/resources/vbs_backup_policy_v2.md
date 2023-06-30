---
subcategory: "Volume Backup Service (VBS)"
---

Up-to-date reference of API arguments for VBS backup policy you can get at
`https://docs.otc.t-systems.com/volume-backup-service/api-ref/api_description/backup_policies`.

# opentelekomcloud_vbs_backup_policy_v2

Provides an VBS Backup Policy resource within OpenTelekomCloud.

## Example Usage

### Basic Backup Policy

```hcl
resource "opentelekomcloud_vbs_backup_policy_v2" "vbs_policy1" {
  name   = "policy_001"
  status = "ON"

  start_time          = "12:00"
  retain_first_backup = "N"
  rentention_num      = 7
  frequency           = 1

  tags = {
    key   = "k1"
    value = "v1"
  }
}
```

### Backup Policy with EVS Disks

```hcl
variable "volume_id" {}

resource "opentelekomcloud_vbs_backup_policy_v2" "vbs_policy2" {
  name   = "policy_002"
  status = "ON"

  start_time          = "12:00"
  retain_first_backup = "N"
  rentention_num      = 5

  frequency = 3
  resources = [var.volume_id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the policy name. The value is a string of 1 to 64 characters that
  can contain letters, digits, underscores (_), and hyphens (-). It cannot start with `default`.

* `start_time` - (Required) Specifies the start time(UTC) of the backup job. The value is in the
  HH:mm format. You need to set the execution time on a full hour. You can set multiple execution
  times, and use commas (,) to separate one time from another.

* `status` - (Optional) Specifies the backup policy status. Possible values are ON or OFF. Defaults to ON.

* `retain_first_backup` - (Required) Specifies whether to retain the first backup in the current month.
  Possible values are Y or N.

* `rentention_num` - (Optional) Specifies number of retained backups. Minimum value is 2.
  Either this field or `rentention_day` must be specified.

* `rentention_day` - (Optional) Specifies days of retained backups. Minimum value is 2.
  Either this field or `rentention_num` must be specified.

* `frequency` - (Optional) Specifies the backup interval. The value is in the range of 1 to 14 days.
  Either this field or `week_frequency` must be specified.

* `week_frequency` - (Optional) Specifies on which days of each week backup jobs are executed.
  The value can be one or more of the following: "SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT".
  Either this field or `frequency` must be specified.

* `resources` - (Optional) Specifies one or more volumes associated with the backup policy.
  Any previously associated backup policy will no longer apply.

* `tags` - (Optional) Represents the list of tags to be configured for the backup policy.

  * `key` - (Required) Specifies the tag key. A tag key consists of up to 36 characters, chosen from letters, digits, hyphens (-), and underscores (_).

  * `value` - (Required) Specifies the tag value. A tag value consists of 0 to 43 characters, chosen from letters, digits, hyphens (-), and underscores (_).


## Attributes Reference

All of the argument attributes are also exported as result attributes:

* `id` - Specifies a backup policy ID.

* `policy_resource_count` - Specifies the number of volumes associated with the backup policy.

## Import

Backup Policy can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vbs_backup_policy_v2.vbs 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
