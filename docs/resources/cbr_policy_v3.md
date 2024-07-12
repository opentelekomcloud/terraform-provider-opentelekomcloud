---
subcategory: "Cloud Backup and Recovery (CBR)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cbr_policy_v3"
sidebar_current: "docs-opentelekomcloud-resource-cbr-policy-v3"
description: |-
  Manages a CBR Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CBR policy you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-backup-recovery/api-ref/cbr_apis/policies)

# opentelekomcloud_cbr_policy_v3

Manages a V3 CBR policy resource within OpenTelekomCloud.

## Example usage

```hcl
resource "opentelekomcloud_cbr_policy_v3" "policy" {
  name           = "some-policy"
  operation_type = "backup"

  trigger_pattern = [
    "FREQ=WEEKLY;BYDAY=MO,TU,WE,TH,FR,SA,SU;BYHOUR=14;BYMINUTE=00"
  ]
  operation_definition {
    day_backups   = 1
    week_backups  = 2
    year_backups  = 3
    month_backups = 4
    max_backups   = 10
    timezone      = "UTC+03:00"
  }

  enabled = "false"
}
```

## Argument reference

The following arguments are supported:

* `enabled` - (Optional) Whether to enable the policy. Default value is `true`.

* `name` - (Required) Specifies the policy name. The value consists of 1 to 64 characters
  and can contain only letters, digits, underscores (_), and hyphens (-).

* `operation_definition` - (Optional) Scheduling parameter. See reference below.

* `operation_type` - (Required) Policy type. Enumeration values: `backup`, `replication`.

* `trigger_pattern` - (Required) Scheduling rule. In the replication policy, you are advised
  to set one time point for one day. A maximum of 24 rules can be configured. The scheduling
  rule complies with iCalendar RFC 2445, but it supports only parameters `FREQ`, `BYDAY`, `BYHOUR`,
  `BYMINUTE`, and `INTERVAL`. `FREQ` can be set only to `WEEKLY` and `DAILY`.

The `operation_definition` block contains:

* `day_backups` - (Optional) Specifies the number of retained daily backups. The latest
  backup of each day is saved in the long term. This parameter can be effective together
  with the maximum number of retained backups specified by `max_backups`. The value ranges
  from `0` to `100`. If this parameter is configured, `timezone` is mandatory.

* `week_backups` - (Optional) Specifies the number of retained weekly backups. The latest
  backup of each week is saved in the long term. This parameter can be effective together
  with the maximum number of retained backups specified by `max_backups`. The value ranges
  from `0` to `100`. If this parameter is configured, `timezone` is mandatory.

* `month_backups` - (Optional) Specifies the number of retained monthly backups. The latest
  backup of each month is saved in the long term. This parameter can be effective together
  with the maximum number of retained backups specified by `max_backups`. The value ranges from
  `0` to `100`. If this parameter is configured, `timezone` is mandatory.

* `year_backups` - (Optional) Specifies the number of retained yearly backups. The latest
  backup of each year is saved in the long term. This parameter can be effective together
  with the maximum number of retained backups specified by `max_backups`. The value ranges
  from `0` to `100`. If this parameter is configured, `timezone` is mandatory.

* `timezone` - (Required) Time zone where the user is located, for example, `UTC+00:00`.

* `max_backups` - (Optional) Maximum number of retained backups. The value can be `-1` or ranges
  from `0` to `99999`. If the value is set to `-1`, the backups will not be cleared even though
  the configured retained backup quantity is exceeded. If this parameter and `retention_duration_days`
  are both left blank, the backups will be retained permanently.

* `retention_duration_days` - (Optional) Duration of retaining a backup, in days.
  The maximum value is `99999`. `-1` indicates that the backups will not be cleared based on
  the retention duration. If this parameter and `max_backups` are left blank at the same time,
  the backups will be retained permanently.

## Attributes Reference

The following attributes are exported:

* `enabled` - See Argument Reference above.

* `name` - See Argument Reference above.

* `operation_type` - See Argument Reference above.

* `trigger_pattern` - See Argument Reference above.

* `region` - Specifies the region of the CBRv3 policy.
