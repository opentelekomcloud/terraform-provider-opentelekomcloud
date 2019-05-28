---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_certificate_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-certificate-v1"
description: |-
  Manages a V1 WAF certificate resource within OpenTelekomCloud.
---

# opentelekomcloud_waf_certificate_v1

Manages a WAF certificate resource within OpenTelekomCloud.

## Example Usage

```hcl

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_1"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The certificate name. The maximum length is 256 characters. Only digits, letters, underscores(_), and hyphens(-) are allowed.

* `content` - (Optional) The certificate content. Changing this creates a new certificate.

* `key` - (Optional) The private key. Changing this creates a new certificate.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the certificate.

* `name` -  See Argument Reference above.

* `content` - See Argument Reference above.

* `key` - See Argument Reference above.

## Import

Certificates can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_waf_certificate_v1.cert_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
