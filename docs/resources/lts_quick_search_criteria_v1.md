---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_quick_search_criteria_v1"
sidebar_current: "docs-opentelekomcloud-resource-lts-quick-search-criteria-v1"
description: |-
  Manages an LTS quick search criteria resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for LTS log group you can get at
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/quick_search/index.html)

# opentelekomcloud_lts_quick_search_criteria_v1

Manages an LTS quick search criteria resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "opentelekomcloud_lts_quick_search_criteria_v1" "qsc" {
  log_group_id  = var.log_group_id
  log_stream_id = var.log_stream_id

  criteria = "content:test"
  name     = "search_criteria_test"
  type     = "ORIGINALLOG"
}
```

## Argument Reference

The following arguments are supported:

* `log_group_id` - (Required, String, ForceNew) Specifies the ID of a log group.

  Changing this parameter will create a new resource.

* `log_stream_id` - (Required, String, ForceNew) Specifies the ID of a log stream.

  Changing this parameter will create a new resource.

* `criteria` - (Required, String, ForceNew) Specifies the content of search criteria.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the search criteria. The name can only contain English
  letters, numbers, hyphens, underscores, and periods. It cannot start with a period or underscore
  or end with a period.

  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the search criteria. Available types are
  `ORIGINALLOG` (for raw logs).

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `region` - Shows the region in the search criteria resource created.


## Import

The search criteria can be imported using the group ID, stream ID, and resource ID separated by the slashes, e.g.

```bash
$ terraform import opentelekomcloud_lts_quick_search_criteria_v1.test <log_group_id>/<log_stream_id>/<id>
```
