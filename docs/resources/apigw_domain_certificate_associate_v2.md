---
subcategory: "APIGW"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_apigw_domain_certificate_associate_v2"
sidebar_current: "docs-opentelekomcloud-resource-apigw-domain-certificate-associate-v2"
description: |-
  Manages an APIGW domain certificate association resource within T-Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for API Gateway SSL certificate management you can get at
[documentation portal](https://docs.otc.t-systems.com/api-gateway/api-ref/dedicated_gateway_apis_v2/ssl_certificate_management/index.html)

# opentelekomcloud_apigw_domain_certificate_associate_v2

Associates an SSL certificate with a custom domain of a dedicated gateway instance within T-Cloud Public (formerly OpenTelekomCloud).

## Example Usage

```hcl
variable "gateway_id" {}
variable "group_id" {}
variable "domain_id" {}
variable "cert_id" {}

resource "opentelekomcloud_apigw_domain_certificate_associate_v2" "assoc" {
  gateway_id     = var.gateway_id
  group_id       = var.group_id
  domain_id      = var.domain_id
  certificate_id = var.cert_id
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required, String, ForceNew) Specifies the ID of the dedicated gateway instance to which the
  domain belongs.
  Changing this will create a new resource.

* `group_id` - (Required, String, ForceNew) Specifies the ID of the API group to which the domain belongs.
  Changing this will create a new resource.

* `domain_id` - (Required, String, ForceNew) Specifies the ID of the custom domain to bind the certificate to.
  Changing this will create a new resource.

* `certificate_id` - (Required, String, ForceNew) Specifies the ID of the SSL certificate to bind to the domain.
  Changing this will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the association in the format `<gateway_id>/<group_id>/<domain_id>/<certificate_id>`.

* `bound` - Indicates whether the certificate is currently bound to the domain.

## Import

Associations can be imported using the ID of the dedicated gateway instance, the ID of the API group,
the domain ID and the certificate ID, separated by slashes, e.g.

```
$ terraform import opentelekomcloud_apigw_domain_certificate_associate_v2.assoc <gateway_id>/<group_id>/<domain_id>/<certificate_id>
```
