---
subcategory: "Distributed Cache Service (DCS)"
---

Up-to-date reference of API arguments for DCS AZ information you can get at
[https://docs.otc.t-systems.com/distributed-cache-service/api-ref](https://docs.otc.t-systems.com/distributed-cache-service/api-ref/apis_v2_recommended/network_security/downloading_the_ssl_certificate_of_an_instance.html#downloadsslcert).

# opentelekomcloud_dcs_certificate_v2

Use this data source to get the certificate of OpenTelekomCloud DCS instance.

~>
    SSL certificate download is available only for DCS 6.0 instances with enabled SSL.

## Example Usage

```hcl
variable "dcs_id" {}
data "opentelekomcloud_dcs_certificate_v2" "cert" {
  instance_id = var.dcs_id
}
```

## Argument Reference

* `instance_id` - (Required) A DCS instance ID.

## Attributes Reference

All above argument parameters can be exported as attribute parameters along with attribute reference.

* `file_name` - SSL certificate file name.

* `link` - Download link of the SSL certificate.

* `bucket_name` - Name of the OBS bucket for storing the SSL certificate.

* `certificate` - SSL certificate of an instance.
