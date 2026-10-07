---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_domain_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-domain-v2"
description: |-
  Manages an APIGW domain resource within T-Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for API Gateway domain service you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/domain_name_management/index.html)

# opentelekomcloud_apigw_domain_v2

Manages an API Gateway domain resource within T-Cloud Public (formerly OpenTelekomCloud).

## Example Usage

```hcl
variable "gateway_id" {}
variable "group_id" {}
variable "domain_name" {}

resource "opentelekomcloud_apigw_domain_v2" "test" {
  gateway_id             = var.gateway_id
  group_id               = var.group_id
  name                   = var.domain_name
  min_ssl_version        = "TLSv1.2"
  http_redirect_to_https = false
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway instance to which the
  domain belongs.
  Changing this will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the API group to which the domain belongs.
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the custom domain name.
  It can contain a maximum of `255` characters and must comply with domain name specifications.
  Changing this will create a new resource.

* `min_ssl_version` - (Optional, String) Specifies the minimum SSL version.
  The valid values are `TLSv1.1` and `TLSv1.2`.

* `http_redirect_to_https` - (Optional, Bool) Specifies whether to enable HTTP redirection to HTTPS.
  The value **false** means disable and **true** means enable. The default value is **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the custom domain.

* `status` - The CNAME resolution status.
  The valid values are as follows:
  * `1`: not resolved.
  * `2`: resolving.
  * `3`: resolved.
  * `4`: resolution failed.

## Import

Domains can be imported using the ID of the dedicated gateway instance, the ID of the API group and the domain ID,
separated by slashes, e.g.

```
$ terraform import opentelekomcloud_apigw_domain_v2.test <gateway_id>/<group_id>/<domain_id>
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below:

```hcl
resource "opentelekomcloud_apigw_domain_v2" "test" {
  # ...

  lifecycle {
    ignore_changes = [
      name,
      http_redirect_to_https
    ]
  }
}
```
