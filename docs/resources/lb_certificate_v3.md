---
subcategory: "Dedicated Load Balancer (DLB)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lb_certificate_v3"
sidebar_current: "docs-opentelekomcloud-resource-lb-certificate-v3"
description: |-
  Manages a LB Certificate resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for DLB certificate you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-load-balancing/api-ref/apis_v3/certificate)

# opentelekomcloud_lb_certificate_v3

Manages a V3 certificate resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  description = "terraform test certificate"
  domain      = "www.elb.com"
  private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIBUwIBADANBgkqhkiG9w0BAQEFAASCAT0wggE5AgEAAkEAu+qgVpV6mqbaGW1Q
n6eDPzhwentQPPiXwG1665M9+gjW4pUQ0RudBc0fkUU/O+Q0UMT8ZV/I2hSenCVy
JoyPEwIDAQABAkAbyksEAv8qt9oxQHVX5xIF23bm5i2rlqf6kTZIeHIF89/NNJ2E
sejiqFIWqPc5a00Scn+ymdCvjC25JVyup9cBAiEA4a+7WhPmgS54yNHjwkG2pflz
cfH1V7qPqlBKIGLwZbMCIQDVKCsZ6eoNdQoLVmK0zii8XDCgL8HWMrm/bytbYM9B
IQIgVdcAXKebEeF6IW/rwDQ8Y2644UsVdTPJdw8o0p6vLw8CIDqm191EiPt09fOS
rIxVoc3ajCK3oV2ADa5IN6ToKX8hAiBPuNCCIYcZz0tAzWX7I1OYMI3UhJjtrESg
mYFrsJ4gHw==
-----END RSA PRIVATE KEY-----
EOT

  certificate = <<EOT
-----BEGIN CERTIFICATE-----
MIIB4TCCAYugAwIBAgIUPXCpWJCiy5mI79NIfenl5KNWPzkwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMTExMDIxMDM3MjBaFw0yMTEy
MDIxMDM3MjBaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEF
AANLADBIAkEAu+qgVpV6mqbaGW1Qn6eDPzhwentQPPiXwG1665M9+gjW4pUQ0Rud
Bc0fkUU/O+Q0UMT8ZV/I2hSenCVyJoyPEwIDAQABo1MwUTAdBgNVHQ4EFgQUtItI
IAXZDIEfuvCX7AY3s//wlI8wHwYDVR0jBBgwFoAUtItIIAXZDIEfuvCX7AY3s//w
lI8wDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAANBAEkgP/JlpVKc4j+Z
KRcMa7RAXYJqCbRxtpqRU7OOAhDmBnldtS5CTMoh1r7TOGMfM1Npa+kGV5QnjRzI
FzFSymo=
-----END CERTIFICATE-----
EOT
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V3 ELB client.
  An ELB client is needed to create an LB certificate. If omitted, the
  `region` argument of the provider is used. Changing this creates a new
  LB certificate.

* `name` - (Optional) Specifies the certificate name. Only letters,
  digits, underscores, and hyphens are allowed.

* `description` - (Optional) Provides supplementary information about the certificate.

* `domain` - (Optional) The domain of the Certificate.

* `private_key` - (Optional) The private encrypted key of the Certificate, PEM format.
  Required for certificates of type `server`.

* `certificate` - (Required) The public encrypted key of the Certificate, PEM format.

* `type`- (Optional) The type of certificate the container holds. Either `server` or `client`.
  Defaults to `server` if not set. Changing this creates a new LB certificate.

## Attributes Reference

In addition, the following attributes are exported:

* `updated_at` - Indicates the update time.

* `createa_at` - Indicates the creation time.

* `expire_time` - Indicates the expiration time.

## Import

Certificates can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_lb_certificate_v3.certificate_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
