---
subcategory: "Enterprise Router (ER)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_er_availability_zones_v3"
sidebar_current: "docs-opentelekomcloud-datasource-er-availability-zones-v3"
description: |-
  Query AZs where enterprise routers can be created within OpenTelekomCloud
---

# opentelekomcloud_er_availability_zones_v3

Use this data source to query availability zones where ER instances can be created within OpenTelekomCloud.

## Example Usage

```hcl
data "opentelekomcloud_er_availability_zones_v3" "zones" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `id` - The data source ID.

* `names` - The names of availability zone.

* `region` - The region where resources are located.
