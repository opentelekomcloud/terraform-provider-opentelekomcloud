---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_certificate_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lb-certificate-v3"
description: |-
Get ELBv3 certificate from OpenTelekomCloud
---

Up-to-date reference of API arguments for ELBv3 certificate you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/certificate/querying_certificates.html#listcertificates)

# opentelekomcloud_lb_certificate_v3

Use this data source to get the info about an existing ELBv3 certificate.

## Example Usage

```hcl
data "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name = var.certificate_id
}
```

## Argument Reference

* `id` - (Optional) Specifies the certificate ID.

* `name` - (Optional) Specifies the certificate name.

* `domain` - (Optional) The domain of the Certificate.

* `type`- (Optional) The type of certificate the container holds. Either `server` or `client`.

## Attributes Reference

In addition, the following attributes are exported:

* `private_key` - The private encrypted key of the Certificate, PEM format.

* `certificate` - The public encrypted key of the Certificate, PEM format.

* `description` - Provides supplementary information about the certificate.

* `updated_at` - Indicates the update time.

* `createa_at` - Indicates the creation time.

* `expire_time` - Indicates the expiration time.
