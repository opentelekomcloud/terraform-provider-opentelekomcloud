---
subcategory: "Cloud Eye (CES)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_ces_dashboard_v2"
sidebar_current: "docs-opentelekomcloud-resource-ces-dashboard-v2"
description: |-
  Manages a CES Dashboard v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CES dashboard you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-eye/api-ref/api_v2/dashboards/index.html)

# opentelekomcloud_ces_dashboard_v2

Manages a CES Dashboard v2 resource within OpenTelekomCloud.

## Example Usage

### Basic dashboard

```hcl
resource "opentelekomcloud_ces_dashboard_v2" "dashboard" {
  name           = "my-dashboard"
  row_widget_num = 1
  is_favorite    = true
}
```

### Copy an existing dashboard

```hcl
resource "opentelekomcloud_ces_dashboard_v2" "base" {
  name           = "base-dashboard"
  row_widget_num = 2
}

resource "opentelekomcloud_ces_dashboard_v2" "copy" {
  name           = "copied-dashboard"
  dashboard_id   = opentelekomcloud_ces_dashboard_v2.base.id
  row_widget_num = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the dashboard name. The value can be a string of `1` to `128`
  characters that can consist of letters, digits, underscores (_), and hyphens (-).

* `row_widget_num` - (Optional, Int) Specifies the graph display mode. The value can be:
  + **0**: Graphs are displayed in a customizable position.
  + **1**: One graph is displayed per row.
  + **2**: Two graphs are displayed per row.
  + **3**: Three graphs are displayed per row.

  Defaults to `0`.

* `is_favorite` - (Optional, Bool) Specifies whether to add the dashboard to favorites.

* `dashboard_id` - (Optional, String, ForceNew) Specifies the ID of an existing dashboard to copy.
  If omitted, a new empty dashboard is created. Changing this creates a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The dashboard ID.

* `creator_name` - The creator of the dashboard.

* `created_at` - The creation time of the dashboard in RFC3339 format.

* `region` - The region in which the dashboard is created.

## Import

CES dashboards v2 can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_ces_dashboard_v2.dashboard db16564943172807wjOmoLyn
```
