---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_host_access_v3"
sidebar_current: "docs-opentelekomcloud-resource-lts-host-access-v3"
description: |-
  Manages a LTS Host Access resource within OpenTelekomCloud.
---

# opentelekomcloud_lts_host_access_v3

Manages an LTS host access resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "stream_id" {}
variable "host_group_id" {}

resource "opentelekomcloud_lts_host_access_v3" "acc" {
  name           = "access-demo"
  log_group_id   = var.group_id
  log_stream_id  = var.stream_id
  host_group_ids = [var.host_group_id]

  access_config {
    paths       = ["/var/log/*"]
    black_paths = ["/var/log/*/a.log"]

    single_log_format {
      mode = "system"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required, String, ForceNew) Specifies the host access name. The name consists of `1` to `64` characters.
  Only letters, digits, underscores (_), and periods (.) are allowed, and the period cannot be the first or last character.
  Changing this parameter will create a new resource.

* `log_group_id` - (Required, String, ForceNew) Specifies the log group ID.
  Changing this parameter will create a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the log stream ID.
  Changing this parameter will create a new resource.

* `access_config` - (Required, List) Specifies the configurations of host access.
  The [access_config](#HostAccessConfigDetail) structure is documented below.

* `host_group_ids` - (Optional, List) Specifies the log access host group ID list.

* `tags` - (Optional, Map) Specifies the key/value to attach to the host access.

<a name="HostAccessConfigDetail"></a>
The `access_config` block supports:

* `paths` - (Required, List) Specifies the collection paths.

  + A path must start with `/` or `Letter:\`.
  + A path cannot contain only slashes (/). The following special characters are not allowed: <>'|"
  + A path cannot start with `/**` or `/*`.
  + Only one double asterisk (**) can be contained in a path.
  + Up to 10 paths can be specified.

* `black_paths` - (Optional, List) Specifies the collection path blacklist.

  + A path must start with `/` or `Letter:\`.
  + A path cannot contain only slashes (/). The following special characters are not allowed: <>'|"
  + A path cannot start with `/**` or `/*`.
  + Only one double asterisk (**) can be contained in a path.
  + Up to 10 paths can be specified.

  -> If you blacklist a file or directory that has been set as a collection path, the blacklist settings
    will be used and the file or files in the directory will be filtered out.

* `single_log_format` - (Optional, List) Specifies the configuration single-line logs. Each log line is displayed as a
  single log event. The [single_log_format](#HostAccessConfigSingleLogFormat) structure is documented below.

* `multi_log_format` - (Optional, List) Specifies the configuration multi-line logs. Multiple lines of exception log events
  can be displayed as a single log event. This is helpful when you check logs to locate problems.
  The [multi_log_format](#HostAccessConfigMultiLogFormat) structure is documented below.

* `windows_log_info` - (Optional, List) Specifies the configuration of Windows event logs.
  The [windows_log_info](#HostAccessConfigWindowsLogInfo) structure is documented below.

* `binary_collect` - (Optional, Bool) Specifies whether collect in binary format. Default is **false**.

* `log_split` - (Optional, Bool) Specifies whether to split log. Default is false.

<a name="HostAccessConfigSingleLogFormat"></a>
The `single_log_format` blocks supports:

* `mode` - (Required, String) Specifies mode of single-line log format. The options are as follows:
  + `system`: the system time.
  + `wildcard`: the time wildcard.

* `value` - (Optional, String) Specifies value of single-line log format.
  + If mode is `system`, the value is the current timestamp, the number of milliseconds elapsed since January 1, 1970 UTC.
  + If mode is `wildcard`, the value is `required` and is a time wildcard, which is used to look for the log printing
    time as the beginning of a log event. If the time format in a log event is `2019-01-01 23:59:59`, the time wildcard is
    `YYYY-MM-DD hh:mm:ss`. If the time format in a log event is `19-1-1 23:59:59`, the time wildcard is `YY-M-D hh:mm:ss`.

<a name="HostAccessConfigMultiLogFormat"></a>
The `multi_log_format` blocks supports:

* `mode` - (Required, String) Specifies mode of multi-line log format. The options are as follows:
  + `time`: the time wildcard.
  + `regular`: the regular expression.

* `value` - (Required, String) Specifies value of multi-line log format.
  + If mode is `regular`, the value is a regular expression.
  + If mode is `time`, the value is a time wildcard, which is used to look for the log printing time as the beginning
    of a log event. If the time format in a log event is `2019-01-01 23:59:59`, the time wildcard is `YYYY-MM-DD hh:mm:ss`.
    If the time format in a log event is `19-1-1 23:59:59`, the time wildcard is `YY-M-D hh:mm:ss`.

-> The time wildcard and regular expression will look for the specified pattern right from the beginning of each log line.
  If no match is found, the system time, which may be different from the time in the log event, is used. In general cases,
  you are advised to select `Single-line` for Log Format and `system` time for Log Time.

<a name="HostAccessConfigWindowsLogInfo"></a>
The `windows_log_info` block supports:

* `categories` - (Required, List) Specifies the types of Windows event logs to collect. The valid values are
  `Application`, `System`, `Security` and `Setup`.

* `event_level` - (Required, List) Specifies the Windows event severity. The valid values are `information`, `warning`,
   `error`, `critical` and `verbose`.  Only Windows Vista or later is supported.

* `time_offset_unit` - (Required, String) Specifies the collection time offset unit. The valid values are
  `day`, `hour` and `sec`.

* `time_offset` - (Required, Int) Specifies the collection time offset. This time takes effect only for the first
  time to ensure that the logs are not collected repeatedly.

  + When `time_offset_unit` is set to `day`, the value ranges from `1` to `7` days.
  + When `time_offset_unit` is set to `hour`, the value ranges from `1` to `168` hours.
  + When `time_offset_unit` is set to `sec`, the value ranges from `1` to `604,800` seconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the host access.

* `created_at` - The creation time of the Host access, in RFC3339 format.

* `access_type` - The log access type.

* `log_group_name` - The log group name.

* `log_stream_name` - The log stream name.

* `region` - Shows the region in the host access resource created.

## Import

The host access can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_lts_host_access_v3.test <id>
```
