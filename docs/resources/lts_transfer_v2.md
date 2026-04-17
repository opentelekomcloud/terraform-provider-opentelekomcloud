---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_transfer_v2"
sidebar_current: "docs-opentelekomcloud-resource-lts-transfer-v2"
description: |-
  Manages a LTS Log Transfer resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log transfer you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/log_transfer/index.html)

# opentelekomcloud_lts_transfer_v2

Manages an LTS transfer task resource within OpenTelekomCloud.

## Example Usage

### Create an OBS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "obs_buket" {}

resource "opentelekomcloud_lts_transfer_v2" "test" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = var.obs_buket
      obs_dir_prefix_name = "dir_prefix_"
      obs_prefix_name     = "prefix_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }
  }
}
```

### Create a delegated OBS transfer task

```hcl
variable "lts_group_id" {}
variable "lts_stream_id" {}
variable "obs_buket" {}
variable "agency_domain_id" {}
variable "agency_domain_name" {}
variable "agency_name" {}
variable "agency_project_id" {}

resource "opentelekomcloud_lts_transfer_v2" "obs_agency" {
  log_group_id = var.lts_group_id

  log_streams {
    log_stream_id = var.lts_stream_id
  }

  log_transfer_info {
    log_transfer_type   = "OBS"
    log_transfer_mode   = "cycle"
    log_storage_format  = "RAW"
    log_transfer_status = "ENABLE"

    log_transfer_detail {
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_bucket_name     = var.obs_buket
      obs_dir_prefix_name = "dir_prefix_"
      obs_prefix_name     = "prefix_"
      obs_time_zone       = "UTC"
      obs_time_zone_id    = "Etc/GMT"
    }

    log_agency_transfer {
      agency_domain_id   = var.agency_domain_id
      agency_domain_name = var.agency_domain_name
      agency_name        = var.agency_name
      agency_project_id  = var.agency_project_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `log_group_id` - (Required, String, ForceNew) Log group ID.

  Changing this parameter will create a new resource.

* `log_streams` - (Required, List, ForceNew) The list of log streams.

  Changing this parameter will create a new resource.
  The [log_streams](#LtsTransfer_LogStreams) structure is documented below.

* `log_transfer_info` - (Required, List) Log transfer information.
  The [log_transfer_info](#LtsTransfer_LogTransferInfo) structure is documented below.

<a name="LtsTransfer_LogStreams"></a>
The `log_streams` block supports:

* `log_stream_id` - (Required, String, ForceNew) Log stream ID.

  Changing this parameter will create a new resource.

* `log_stream_name` - (Optional, String, ForceNew) Log stream name.

  Changing this parameter will create a new resource.

<a name="LtsTransfer_LogTransferInfo"></a>
The `log_transfer_info` block supports:

* `log_transfer_type` - (Required, String, ForceNew) Log transfer type.
  The valid values are `OBS`.

  Changing this parameter will create a new resource.

* `log_transfer_mode` - (Required, String, ForceNew) Log transfer mode.
  Value options are as follows:
    + `cycle`: Periodical transfer.

  Changing this parameter will create a new resource.

* `log_storage_format` - (Required, String) Log transfer format.
  Value options are as follows:
    + `JSON`: JSON format.
    + `RAW`: Raw log format.

* `log_transfer_status` - (Required, String) Log transfer status.
  Value options are as follows:
    + `ENABLE`: Log transfer is enabled.
    + `DISABLE`: Log transfer is disabled

* `log_agency_transfer` - (Optional, List, ForceNew) Information about agency which lets an account delegate resource management
   to other accounts.
  This parameter is mandatory if you transfer logs for another account.
  The [log_agency_transfer](#LtsTransfer_LogAgency) structure is documented below.

  Changing this parameter will create a new resource.

* `log_transfer_detail` - (Required, List) Log transfer details.
  The [log_transfer_detail](#LtsTransfer_LogTransferDetail) structure is documented below.

<a name="LtsTransfer_LogAgency"></a>
The `log_agency_transfer` block supports:

* `agency_domain_id` - (Required, String, ForceNew) Delegator account ID.

  Changing this parameter will create a new resource.

* `agency_domain_name` - (Required, String, ForceNew) Delegator account name.

  Changing this parameter will create a new resource.

* `agency_name` - (Required, String, ForceNew) The agency name created by the delegator.

  Changing this parameter will create a new resource.

* `agency_project_id` - (Required, String, ForceNew) Project ID of the delegator.

  Changing this parameter will create a new resource.

<a name="LtsTransfer_LogTransferDetail"></a>
The `log_transfer_detail` block supports:

* `obs_period` - (Optional, Int) Length of the transfer interval for an OBS transfer task.
  This parameter is mandatory when you create an OBS transfer task.
  The log transfer interval is specified by the combination of the values of **obs_period** and **obs_period_unit**,
  and must be set to one of the following: 2 min, 5 min, 30 min, 1 hour, 3 hours, 6 hours, and 12 hours.
  Value options are as follows:
    + `2`: 2 minutes, the `obs_period_unit` must be `min`.
    + `5`: 5 minutes, the `obs_period_unit` must be `min`.
    + `30`: 30 minutes, the `obs_period_unit` must be `min`.
    + `1`: 1 hour, the `obs_period_unit` must be `hour`.
    + `3`: 3 hours, the `obs_period_unit` must be `hour`.
    + `6`: 6 hours, the `obs_period_unit` must be `hour`.
    + `12`: 12 hours, the `obs_period_unit` must be `hour`.

* `obs_period_unit` - (Optional, String) Unit of the transfer interval for an OBS transfer task.
  This parameter is mandatory when you create an OBS transfer task.
  The log transfer interval is specified by the combination of the values of `obs_period` and `obs_period_unit`,
  and must be set to one of the following: `2 min`, `5 min`, `30 min`, `1 hour`, `3 hours`, `6 hours`, and `12 hours`.
  Value options are as follows:
    + `min`: minute.
    + `hour`: hour.

* `obs_bucket_name` - (Optional, String) OBS bucket name.
  This parameter is mandatory when you create an OBS transfer task.

* `obs_transfer_path` - (Optional, String) OBS bucket path, which is the log transfer destination.

* `obs_dir_prefix_name` - (Optional, String) Custom transfer path of an OBS transfer task.

* `obs_prefix_name` - (Optional, String) Transfer file prefix of an OBS transfer task.

* `obs_eps_id` - (Optional, String) Enterprise project ID of an OBS transfer task.

* `obs_encrypted_enable` - (Optional, Bool) Whether OBS bucket encryption is enabled.

* `obs_encrypted_id` - (Optional, String) KMS key ID for an OBS transfer task.
  This parameter is mandatory if encryption is enabled for the target OBS bucket.

* `obs_time_zone` - (Optional, String) Time zone for an OBS transfer task.
  If this parameter is specified, `obs_time_zone_id` must also be specified.

* `obs_time_zone_id` - (Optional, String) ID of the time zone for an OBS transfer task.
  If this parameter is specified, `obs_time_zone` must also be specified.

* `tags` - (Optional, List) The list of tag fields will be delivered when transferring.
  This field must contain the following host information: `hostIP`, `hostId`, `hostName`, `pathFile`, and `collectTime`.
  The common fields include `logStreamName`, `regionName`, `logGroupName`, and `projectId`, which are optional.
  The transfer tag: `streamTag`, which is optional.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `log_group_name` - Log group name.

* `region` - Shows the region in the log group resource created.

## Import

The LTS transfer task can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_lts_transfer_v2.test <id>
```
