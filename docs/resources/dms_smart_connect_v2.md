---
subcategory: "Distributed Message Service (DMS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dms_smart_connect_v2"
sidebar_current: "docs-opentelekomcloud-resource-dms-smart-connect-v2"
description: |-
  Manages an up-to-date DMS Smart Connect v2 resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DMS instance you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-message-service/api-ref/apis_v2_recommended/smart_connect/index.html)

# opentelekomcloud_dms_smart_connect_v2

Manage DMS smart connect v2 resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "opentelekomcloud_dms_smart_connect_v2" "test" {
  instance_id       = var.instance_id
  storage_spec_code = "dms.physical.storage.ultra.v2"
  bandwidth         = "100MB"
  node_count        = 2
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS instance.

  Changing this parameter will create a new resource.

* `storage_spec_code` - (Required, String, ForceNew) Specifies the storage specification code of the connector.

  Changing this parameter will create a new resource.

* `bandwidth` - (Optional, String, ForceNew) Specifies the bandwidth of the connector.

  Changing this parameter will create a new resource.

* `node_count` - (Optional, Int, ForceNew) Specifies the node count of the connector. Defaults to 2 and minimum is 2.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID.

* `region` - The DMS instance region
