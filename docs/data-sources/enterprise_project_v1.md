---
subcategory: "Enterprise Project Management Service (EPS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_enterprise_project_v1"
sidebar_current: "docs-opentelekomcloud-datasource-enterprise-project-v1"
description: "Get an enterprise project details from T-Cloud Public (former OpenTelekomCloud)"
---

# opentelekomcloud_enterprise_project_v1

Use this data source to get an enterprise project from T-Cloud Public (former OpenTelekomCloud)

## Example Usage

```hcl
data "opentelekomcloud_enterprise_project_v1" "test" {
  name = "default"
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the enterprise project name. Fuzzy search is supported.

* `id` - (Optional, String) Specifies the ID of an enterprise project. The value 0 indicates enterprise project default.

* `status` - (Optional, Int) Specifies the status of an enterprise project.
    + 1 indicates Enabled.
    + 2 indicates Disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `description` - Provides supplementary information about the enterprise project.

* `created_at` - Specifies the time (UTC) when the enterprise project was created. Example: 2018-05-18T06:49:06Z

* `updated_at` - Specifies the time (UTC) when the enterprise project was modified. Example: 2018-05-28T02:21:36Z
