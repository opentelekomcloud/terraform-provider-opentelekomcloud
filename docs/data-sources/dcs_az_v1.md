---
subcategory: "Distributed Cache Service (DCS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_dcs_az_v1"
sidebar_current: "docs-opentelekomcloud-datasource-dcs-az-v1"
description: |-
Get DCS AZ information from OpenTelekomCloud
---

Up-to-date reference of API arguments for DCS AZ you can get at
[documentation portal](https://docs.otc.t-systems.com/distributed-cache-service/api-ref/apis_v2_recommended/other_apis/querying_azs.html#listavailablezones)

# opentelekomcloud_dcs_az_v1

Use this data source to get the ID of an available DCS AZ from OpenTelekomCloud.

## Example Usage

### Query AZ `id` by providing `name` argument

```hcl
data "opentelekomcloud_dcs_az_v1" "az1" {
  name = "eu-de-01"
}
```

### Query AZ `id` by providing `port` and `code` arguments

```hcl
data "opentelekomcloud_dcs_az_v1" "az2" {
  port = "8003"
  code = "eu-de-02"
}
```

### Query AZ `id` by providing all arguments

```hcl
data "opentelekomcloud_dcs_az_v1" "az2" {
  name = "eu-de-02"
  port = "8003"
  code = "eu-de-02"
}
```

## Argument Reference

* `name` - (Optional, String) Indicates the name of an AZ.

* `code` - (Optional, String) Indicates the code of an AZ.

* `port` - (Optional, String) Indicates the port number of an AZ.


## Attributes Reference

`id` is set to the ID of the found AZ. In addition, the following attributes are exported:

* `name` - See Argument Reference above.

* `code` - See Argument Reference above.

* `port` - See Argument Reference above.
