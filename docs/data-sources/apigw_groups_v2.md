---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_groups_v2"
sidebar_current: "docs-opentelekomcloud-datasource-apigw-groups-v2"
description: |-
  Get the all APIGW gateway groups from OpenTelekomCloud
---

Up-to-date reference of API arguments for API Gateway groups service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/api_group_management/querying_api_groups.html)

# opentelekomcloud_apigw_groups_v2

Use this data source to query the group list under the APIGW instance within OpenTelekomCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group_name" {}

data "opentelekomcloud_apigw_groups_v2" "test" {
  instance_id = var.instance_id
  name        = var.group_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies an ID of the APIGW dedicated instance to which the API group belongs.

* `group_id` - (Optional, String) Specifies the API group ID used to query.

* `name` - (Optional, String) Specifies the API group name used to query.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Data source ID.

* `region` -  The region in which to query the data source.

* `groups` - All groups that match the filter parameters.
  The [groups](#APIG_Groups) structure is documented below.

<a name="APIG_Groups"></a>
The `groups` block supports:

* `id` - The API group ID.

* `name` - The API group name.

* `status` - The current status of the API group.
  The valid values are as follows:
  + **1**: Normal.

* `sl_domain` - The subdomain name assigned by the system by default.

* `created_at` - The creation time of the API group.

* `updated_at` - The latest update time of the API group.

* `on_sell_status` - Whether it has been listed on the cloud store.
  The valid values are as follows:
  + **1**: Listed.
  + **2**: Not listed.
  + **3**: Under review.

* `url_domains` - List of independent domains bound on the API group.
  The [url_domains](#APIG_Groups_urlDomains) structure is documented below.

* `sl_domains` - List of subdomain names assigned by the system by default.

* `description` - The description of the API group.

* `is_default` - Indicates whether the API group is the default group.

* `environment` - The array of one or more environments of the API group.
  The [environment](#APIG_Groups_environment_attr) structure is documented below.

<a name="APIG_Groups_urlDomains"></a>
The `url_domains` block supports:

* `id` - The domain ID.

* `name` - The domain name.

* `cname_status` - CNAME resolution status of the domain name.
  The valid values are as follows:
  + **1**: Not resolved.
  + **2**: Resolving.
  + **3**: Resolved.
  + **4**: Resolution failed.

* `ssl_id` - The SSL certificate ID.

* `ssl_name` - The SSL certificate name.

* `min_ssl_version` - Minimum SSL version. The default is **TLSv1.1**.
  The valid values are as follows:
  + **TLSv1.1**
  + **TLSv1.2**

* `verified_client_certificate_enabled` - Whether to enable client certificate verification.
  This parameter is available only when a certificate is bound. It is enabled by default if trusted_root_ca exists,
  and disabled if trusted_root_ca does not exist. The default is **false**.

* `is_has_trusted_root_ca` - Whether a trusted root certificate (CA) exists. The value is true
  if trusted_root_ca exists in the bound certificate. The default is **false**.

<a name="APIG_Groups_environment_attr"></a>
The `environment` block supports:

* `variable` - The array of one or more environment variables.
  The [variable](#APIG_Groups_environment_variable_attr) structure is documented below.

* `environment_id` - The ID of the environment to which the variables belong.

<a name="APIG_Groups_environment_variable_attr"></a>
The `variable` block supports:

* `name` - The variable name.

* `value` - The variable value.

* `id` - The variable ID.
