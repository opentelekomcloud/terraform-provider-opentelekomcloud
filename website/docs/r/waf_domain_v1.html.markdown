---
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_domain_v1"
sidebar_current: "docs-opentelekomcloud-resource-waf-domain-v1"
description: |-
  Manages a V1 WAF domain resource within OpenTelekomCloud.
---

# opentelekomcloud_waf_domain_v1

Manages a WAF domain resource within OpenTelekomCloud.

## Example Usage

```hcl

resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
	name = "cert_1"
	content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
	key = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
	hostname = "www.b.com"
	server {
		front_protocol = "HTTPS"
		back_protocol = "HTTP"
		address = "80.158.42.162"
		port = "8080"
	}
	certificate_id = "${opentelekomcloud_waf_certificate_v1.certificate_1.id}"
	proxy = "true"
	sip_header_name = "default"
	sip_header_list = ["X-Forwarded-For"]
}

```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The domain name. For example, www.example.com or *.example.com. Changing this creates a new domain.

* `certificate_id` - (Optional) The certificate ID. This parameter is mandatory when front_protocol is set to HTTPS.

* `server` - (Required) Array of server object. The server object structure is documented below.

* `proxy` - (Required) Specifies whether a proxy is configured.

* `policy_id` - The policy ID associate with the domain. Changing this create a new domain.

* `sip_header_name` - (Optional) The type of the source IP header. This parameter is required only when proxy is set to true. The options are as follows: default, cloudflare, akamai, and custom.

* `sip_header_list` - (Optional) Array of HTTP request header for identifying the real source IP address. This parameter is required only when proxy is set to true.

	* If `sip_header_name` is default, `sip_header_list` is ["X-Forwarded-For"].
	* If `sip_header_name` is cloudflare, `sip_header_list` is ["CF-Connecting-IP", "X-Forwarded-For"].
	* If `sip_header_name` is akamai, `sip_header_list` is ["True-Client-IP"].
	* If `sip_header_name` is custom, you can customize a value.

The `server` block supports:

* `front_protocol` - (Required) Protocol type of the client. The options are HTTP and HTTPS.

* `back_protocol` - (Required) Protocol used by WAF to forward client requests to the server. The options are HTTP and HTTPS.

* `address` - (Required) IP address or domain name of the web server that the client accesses. For example, 192.168.1.1 or www.a.com.

* `port` - (Required) Port number used by the web server. The value ranges from 0 to 65535, for example, 8080.


## Attributes Reference

The following attributes are exported:

* `id` -  ID of the domain.

* `access_code` - The acccess code.

* `cname` - The CNAME value.

* `txt_code` - The TXT record. This attribute is returned only when proxy is set to true.

* `sub_domain` - The subdomain name. This attribute is returned only when proxy is set to true.

* `protect_status` - The WAF mode. -1: bypassed, 0: disabled, 1: enabled.

* `access_status` - Whether a domain name is connected to WAF. 0: The domain name is not connected to WAF, 1: The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are HTTP, HTTPS, and HTTP&HTTPS.

## Import

Domains can be imported using the `id`, e.g.

```
$ terraform import opentelekomcloud_waf_domain_v1.dom_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
