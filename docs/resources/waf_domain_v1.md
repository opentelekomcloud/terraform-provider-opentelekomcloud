---
subcategory: "Web Application Firewall (WAF)"
---

# opentelekomcloud_waf_domain_v1

Manages a WAF domain resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
  name    = "cert_1"
  content = "-----BEGIN CERTIFICATE-----MIIDIjCCAougAwIBAgIJALV96mEtVF4EMA0GCSqGSIb3DQEBBQUAMGoxCzAJBgNVBAYTAnh4MQswCQYDVQQIEwJ4eDELMAkGA1UEBxMCeHgxCzAJBgNVBAoTAnh4MQswCQYDVQQLEwJ-----END CERTIFICATE-----"
  key     = "-----BEGIN RSA PRIVATE KEY-----MIICXQIBAAKBgQDFPN9ojPndxSC4E1pqWQVKGHCFlXAAGBOxbGfSzXqzsoyacotueqMqXQbxrPSQFATeVmhZPNVEMdvcAMjYsV/mymtAwVqVA6q/OFdX/b3UHO+b/VqLo3J5SrM-----END RSA PRIVATE KEY-----"
}

resource "opentelekomcloud_waf_domain_v1" "domain_1" {
  hostname = "www.example.com"
  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTPS"
    address         = "80.158.42.162"
    port            = "443"
  }
  certificate_id  = opentelekomcloud_waf_certificate_v1.certificate_1.id
  proxy           = true
  sip_header_name = "default"
  sip_header_list = ["X-Forwarded-For"]
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) The domain name. For example, `www.example.com` or `*.example.com`.
  Changing this creates a new domain.

* `certificate_id` - (Optional) The certificate ID. This parameter is mandatory when
  `front_protocol`/`client_protocol` is set to `HTTPS`.

* `server` - (Required) Array of server object. The server object structure is documented below.

* `proxy` - (Required) Specifies whether a proxy is configured.

* `policy_id` - The policy ID associate with the domain. Changing this create a new domain.

* `sip_header_name` - (Optional) The type of the source IP header. This parameter is required only when proxy is set to `true`.
  The options are as follows: `default`, `cloudflare`, `akamai`, and `custom`.

* `sip_header_list` - (Optional) Array of HTTP request header for identifying the real source IP address.
  This parameter is required only when proxy is set to `true`.
  * If `sip_header_name` is `default`, `sip_header_list` is `["X-Forwarded-For"]`.
  * If `sip_header_name` is `cloudflare`, `sip_header_list` is `["CF-Connecting-IP", "X-Forwarded-For"]`.
  * If `sip_header_name` is `akamai`, `sip_header_list` is `["True-Client-IP"]`.
  * If `sip_header_name` is `custom`, you can customize a value.

The `server` block supports:

* `client_protocol` - (Optional) Protocol type of the client. The options are HTTP and HTTPS.
  Required if `front_protocol` is not set

* `front_protocol` **DEPRECATED** - (Optional)  Same as `client_protocol`. Required if `client_protocol` is not set.

* `server_protocol` - (Optional) Protocol used by WAF to forward client requests to the server.
  The options are HTTP and HTTPS. Required if `back_protocol` is not set.

* `back_protocol` **DEPRECATED** - (Optional) Same as `server_protocol`. Required if `server_protocol` is not set.

* `address` - (Required) IP address or domain name of the web server that the client accesses.
  For example, `192.168.1.1` or `www.bla-bla.com`.

* `port` - (Required) Port number used by the web server. The value ranges from `0` to `65535`, for example, `8080`.

* `cipher` - (Optional) Cipher suite to use with TLS. Possible values are:
  * `cipher_default` - Default cipher suite: Good browser compatibility, most clients supported, sufficient for most scenarios
  * `cipher_1` - Cipher suite 1: Recommended configuration, the best combination of compatibility and security
  * `cipher_2` - Cipher suite 2: Strict compliance with forward secrecy requirements of PCI DSS and excellent protection, but older browsers may be unable to access the websites
  * `cipher_3` - Cipher suite 3: Support for ECDHE, DHE-GCM, and RSA-AES-GCM algorithms but not CBC

* `tls` - (Optional) Minimum TLS version for accessing the protected domain name  if `client_protocol` is set to `HTTPS`.
  Possible values are: `TLS v1.1` and `TLS v1.2`.

-> `сipher_2`  is not supported if `TLS v1.1` is selected.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` -  ID of the domain.

* `access_code` - The access code.

* `cname` - The CNAME value.

* `txt_code` - The TXT record. This attribute is returned only when proxy is set to `true`.

* `sub_domain` - The subdomain name. This attribute is returned only when proxy is set to `true`.

* `protect_status` - The WAF mode. `-1`: `bypassed`, `0`: `disabled`, `1`: `enabled`.

* `access_status` - Whether a domain name is connected to WAF. `0`: The domain name is not connected to WAF,
  `1`: The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are `HTTP`, `HTTPS`, and `HTTP&HTTPS`.

* `auto_policy_id` - ID of the policy automatically created for the domain.

## Import

Domains can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_waf_domain_v1.dom_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```
