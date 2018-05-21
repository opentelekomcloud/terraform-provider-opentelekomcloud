---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_elb_loadbalancer"
sidebar_current: "docs-opentelekomcloud-resource-elb-loadbalancer"
description: |-
  Manages an Elastic loadbalancer resource within OpenTelekomCloud.
---

# opentelekomcloud\_elb\_loadbalancer

Manages an Elastic loadbalancer resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"
  type = "External"
  bandwidth = 5
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Human-readable name for the Loadbalancer. Does not have
    to be unique.

* `vpc_id` - (Optional) Required for admins.

* `type` - 

* `bandwidth` - 

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.
* `vpc_id` - See Argument Reference above.
* `type` - See Argument Reference above.
* `bandwidth` - See Argument Reference above.
