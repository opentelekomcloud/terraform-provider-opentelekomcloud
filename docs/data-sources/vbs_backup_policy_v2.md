---
subcategory: "Volume Backup Service (VBS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vbs_backup_policy_v2"
sidebar_current: "docs-opentelekomcloud-datasource-vbs-backup-policy-v2"
description: |-
  Get details about a specific VBS backup policy from OpenTelekomCloud
---

Up-to-date reference of API arguments for VBS backup policy you can get at
[documentation portal](https://docs.otc.t-systems.com/volume-backup-service/api-ref/api_description/backup_policies/querying_backup_policies.html#en-us-topic-0043410558)

# opentelekomcloud_vbs_backup_policy_v2

Use this data source to get details about a specific VBS backup policy.

## Example Usage

```hcl
variable "policy_name" {}
variable "policy_id" {}

data "opentelekomcloud_vbs_backup_policy_v2" "policies" {
  name = var.policy_name
  id   = var.policy_id
}
```


## Argument Reference

The arguments of this data source act as filters for querying the available VBS backup policy.
The given filters must match exactly one VBS backup policy whose data will be exported as attributes.

* `id` - (Optional) The ID of the specific VBS backup policy to retrieve.

* `name` - (Optional) The name of the specific VBS backup policy to retrieve.

* `status` - (Optional) The status of the specific VBS backup policy to retrieve. The values can be ON or OFF.

* `filter_tags` - (Optional) Represents the list of tags. Backup policy with these tags will be filtered.

* `key` - (Required) Specifies the tag key. Tag keys must be unique.

* `values` - (Required) Specifies the List of tag values. This list can have a maximum of 10 values and all be unique.



## Attributes Reference

The following attributes are exported:

* `id` - See Argument Reference above.

* `name` - See Argument Reference above.

* `status` - See Argument Reference above.

* `start_time` - Specifies the start time of the backup job.The value is in the HH:mm format.

* `retain_first_backup` - Specifies whether to retain the first backup in the current month.

* `rentention_num` - Specifies number of retained backups.

* `frequency` - Specifies the backup interval. The value is in the range of 1 to 14 days.

* `policy_resource_count` - Specifies the number of volumes associated with the backup policy.

The `tags` block contains:

* `key` - Specifies the tag key.

* `value` - Specifies the tag value.
