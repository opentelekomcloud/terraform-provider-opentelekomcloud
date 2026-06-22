---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_log_configuration_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-log-configuration-v1"
description: |-
  Manages a CFW Log Configuration resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW log configuration you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/log_management/index.html)

# opentelekomcloud_cfw_log_configuration_v1

Manages a CFW Log Configuration resource within OpenTelekomCloud.

## Example Usage: Enabling LTS log configuration for a CFW firewall

```hcl
variable firewall_id {}
variable group_id {}

resource "opentelekomcloud_cfw_log_configuration_v1" "log_config_1" {
  fw_instance_id   = var.firewall_id
  lts_enable       = 1
  lts_log_group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `fw_instance_id` - (Required, String, ForceNew) Specifies the Cloud Firewall instance ID.

* `lts_enable` - (Required, Integer) Specifies whether to enable the LTS service. Valid values are `0` (disable) and `1` (enable).

* `lts_log_group_id` - (Required, String) Specifies the LTS log group ID. The value can be obtained by calling the LTS API for querying all log groups. The format is `log_groups.log_group_id` (the period `[.]` is used to separate different levels of objects).

* `lts_attack_log_stream_id` - (Optional, String) Specifies the attack log stream ID. The value can be obtained by calling the LTS API for querying all log streams in a specified log group.

* `lts_attack_log_stream_enable` - (Optional, Integer) Specifies whether to enable the attack log stream. Valid values are `0` (disable, default) and `1` (enable).

* `lts_access_log_stream_id` - (Optional, String) Specifies the access control log stream ID. The value can be obtained by calling the LTS API for querying all log streams in a specified log group.

* `lts_access_log_stream_enable` - (Optional, Integer) Specifies whether to enable the access control log stream. Valid values are `0` (disable, default) and `1` (enable).

* `lts_flow_log_stream_id` - (Optional, String) Specifies the traffic log stream ID. The value can be obtained by calling the LTS API for querying all log streams in a specified log group.

* `lts_flow_log_stream_enable` - (Optional, Integer) Specifies whether to enable the traffic log stream. Valid values are `0` (disable, default) and `1` (enable).

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the Cloud Firewall instance ID (same as `fw_instance_id`).
