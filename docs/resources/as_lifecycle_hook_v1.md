---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_lifecycle_hook_v1"
sidebar_current: "docs-opentelekomcloud-resource-as-lifecycle-hook-v1"
description: |-
  Manages an AS Lifecycle Hook v1 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for AS lifecycle hook you can get at
[documentation portal](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/lifecycle_hooks)

# opentelekomcloud_as_lifecycle_hook_v1

Manages a V1 AS Lifecycle Hook resource within OpenTelekomCloud.

## Example Usage

### Autoscaling Lifecycle Hook

```hcl
resource "opentelekomcloud_as_lifecycle_hook_v1" "hth_aslifecyclehook" {
  scaling_lifecycle_hook_name = "hth_aslifecyclehook"
  scaling_group_id    = "4579f2f5-cbe8-425a-8f32-53dcb9d9053a"
  scaling_lifecycle_hook_type = "INSTANCE_TERMINATING"
  default_result = "ABANDON"
  default_timeout = 3600
  notification_topic_urn = "urn:smn:regionId:b53e5554fad0494d96206fb84296510b:gsh"
  notification_metadata = "Some customized notification"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_lifecycle_hook_name` - (Required, String, ForceNew) The name of the AS Lifecycle Hook. The name can contain letters, digits, underscores(_), and hyphens(-),and cannot exceed 32 characters.

* `scaling_group_id` - (Required, String, ForceNew) The AS group ID. Changing this creates a new AS lifecycle hook.

* `scaling_lifecycle_hook_type` - (Required, String) The lifecycle hook type. The values can be `INSTANCE_TERMINATING`, and `INSTANCE_LAUNCHING`. 
  - INSTANCE_TERMINATING: The hook suspends the instance when it is terminated.
  - INSTANCE_LAUNCHING: The hook suspends the instance when it is started.

* `default_result` - (Optional, String) The default lifecycle hook callback operation.  This operation is performed when the timeout duration expires. The values
  can be `ABANDON` (default value), and `CONTINUE`.
	- ABANDON:
	  If an instance is starting, ABANDON indicates that your customized operations failed, and the instance will be terminated.
	  In such a case, the scaling action fails, and you must create a new instance.
	  If an instance is stopping, ABANDON allows instance termination BUT stops other lifecycle hooks.
	- CONTINUE:
	  If an instance is starting, CONTINUE indicates that your customized operations are successful and the instance can be used.
	  If an instance is stopping, CONTINUE allows instance termination AND the completion of other lifecycle hooks.

* `default_timeout` - (Optional, Integer) the lifecycle hook timeout duration, which ranges from 60 to 86400 seconds. The default value is 3600.

* `notification_topic_urn` - (Required, String) The URN of an SMN topic. This parameter specifies a notification object for a lifecycle hook. When an instance is suspended by the lifecycle hook, the SMN service sends a notification to the object. This notification contains the basic instance information, your customized notification content, and the token for controlling lifecycle operations.

* `notification_metadata` - (Optional, String) A customized notification, which contains no more than 256 characters. The message cannot contain the following characters: <>&'(){}.

## Attributes Reference

The following extra attributes are exported:

* `notification_topic_name` - (String) Name of the associated topic in SMN..

* `create_time` - (String) Time of creation of the autoscaling lifecycle hook.
