---
subcategory: "VPC Endpoint (VPCEP)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpcep_public_service_v1"
sidebar_current: "docs-opentelekomcloud-datasource-vpcep-public-service-v1"
description: |-
Get details about a specific VPCEP public service from OpenTelekomCloud
---

Up-to-date reference of API arguments for VPCEP public service you can get at
[documentation portal](https://docs.otc.t-systems.com/vpc-endpoint/api-ref/apis/apis_for_managing_vpc_endpoint_services/querying_public_vpc_endpoint_services.html)

# opentelekomcloud_vpcep_public_service_v1

Use this data source to get details about a specific VPCEP public service.

## Example Usage

```hcl
data "opentelekomcloud_vpcep_public_service_v1" "obs" {
  name = "com.t-systems.otc.eu-de.obs"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Specifies the unique ID of the public VPC endpoint service.

* `name` - (Optional) Specifies the name of the public VPC endpoint service. The value is not case-sensitive and supports fuzzy match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `owner` - Specifies the owner of the VPC endpoint service.

* `service_type` - Specifies the type of the VPC endpoint service.

* `created_at` - Specifies the creation time of the VPC endpoint service.

* `is_charge` - Specifies whether the associated VPC endpoint carries a charge.
