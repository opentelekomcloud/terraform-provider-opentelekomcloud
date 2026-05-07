---
subcategory: "Web Application Firewall (WAF)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_cloud_instance_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-cloud-instance-v1"
description: |-
  Manages a cloud WAF instance resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for cloud WAF you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/managing_your_subscriptions/index.html).

# opentelekomcloud_waf_cloud_instance_v1

Manages a postpaid cloud WAF instance resource within OpenTelekomCloud.

-> **Note:** This resource uses the current postpaid cloud WAF API only. Updating the resource in place is not supported.

## Example Usage

```hcl
resource "opentelekomcloud_waf_cloud_instance_v1" "cloud_1" {
  charging_mode = "postPaid"
  website       = "dt"
}
```

## Argument Reference

The following arguments are supported:

* `charging_mode` - (Optional, ForceNew, String) The charging mode of the cloud WAF.
  Defaults to `postPaid`. Only `postPaid` is currently supported, but this argument is kept in the schema for forward compatibility.

* `website` - (Required, ForceNew, String) Website to which the account belongs.
  The current API expects the cloud website value, for example `dt`.

* `enterprise_project_id` - (Optional, ForceNew, String) The ID of the enterprise project to which the cloud WAF belongs.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of the postpaid protected-domain WAF subscription item.

* `resource_spec_code` - The specification code returned by the cloud WAF API.

* `status` - The current resource status. The value can be:
  + `0` - Normal.
  + `1` - Frozen.
  + `2` - Deleted.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

Cloud WAF instances can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_waf_cloud_instance_v1.cloud_1 239a19866000479b8d22d4954980f790
```

-> **Note:** The API does not return the `website` argument on read. After import, keep `website` in your Terraform configuration to avoid drift.
The API also does not return `charging_mode`; keep it in configuration as `postPaid`.
