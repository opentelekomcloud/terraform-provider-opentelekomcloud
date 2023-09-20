---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated domain you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/managing_websites_protected_in_dedicated_mode/index.html).

# opentelekomcloud_waf_dedicated_domain_v1

Manages a WAF dedicated domain resource within OpenTelekomCloud.

-> **Note:** For this resource region must be set in environment variable `OS_REGION_NAME` or in `clouds.yaml`

## Example Usage

### Basic
```hcl
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "my_subnet"
}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain      = "www.mydom.com"
  keep_policy = false
  proxy       = true

  server {
    client_protocol = "HTTP"
    server_protocol = "HTTP"
    address         = "192.168.0.10"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }
}
```

### With certificate
```hcl
data "opentelekomcloud_vpc_subnet_v1" "shared_subnet" {
  name = "my_subnet"
}

resource "opentelekomcloud_waf_dedicated_certificate_v1" "certificate_1" {
  name    = "certificate_1"
  content = <<EOT
-----BEGIN CERTIFICATE-----
MIIFazCCA1OgAwIBAgIUN3w1KX8/T/HWVxZIOdHXPhUOnsAwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
...
dKvZbPEsygYRIjwyhHHUh/YXH8KDI/uu6u6AxDckQ3rP1BkkKXr5NPBGjVgM3ZI=
-----END CERTIFICATE-----
EOT
  key     = <<EOT
-----BEGIN PRIVATE KEY-----
MIIJQQIBADANBgkqhkiG9w0BAQEFAASCCSswggknAgEAAoICAQC+9uwFVenCdPD9
5LWSWMuy4riZW718wxBpYV5Y9N8nM7N0qZLLdpImZrzBbaBldTI+AZGI3Nupuurw
...
s9urs/Kk/tbQhsEvu0X8FyGwo0zH6rG8apTFTlac+v4mJ4vlpxSvT5+FW2lgLISE
+4sM7kp0qO3/p+45HykwBY5iHq3H
-----END PRIVATE KEY-----
EOT

}

resource "opentelekomcloud_waf_dedicated_domain_v1" "domain_1" {
  domain         = "www.mydom.com"
  certificate_id = opentelekomcloud_waf_dedicated_certificate_v1.certificate_1.id
  keep_policy    = false
  proxy          = false
  tls            = "TLS v1.1"
  cipher         = "cipher_1"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "192.168.0.20"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  }

  depends_on = [
    opentelekomcloud_waf_dedicated_certificate_v1.certificate_1
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, ForceNew) The region in which to create the dedicated mode domain resource. If omitted,
  the provider-level region will be used. Changing this setting will push a new domain.

* `domain` - (Required, ForceNew) Specifies the protected domain name or IP address (port allowed). For example,
  `www.example.com` or `*.example.com` or `www.example.com:89`. Changing this creates a new domain.

* `server` - (Required, ForceNew) The server configuration list of the domain. A maximum of 80 can be configured.
  The `server` block supports:

  + `client_protocol` - (Required, ForceNew) Protocol type of the client. Values are `HTTP` and `HTTPS`.
    Changing this creates a new server.

  + `server_protocol` - (Required, ForceNew) Protocol used by WAF to forward client requests to the server.
    Values are`HTTP` and `HTTPS`. Changing this creates a new server.

  + `vpc_id` - (Required, ForceNew) The id of the vpc used by the server. Changing this creates a server.

  + `type` - (Required, ForceNew) Server network type, IPv4 or IPv6. Valid values are: `ipv4` and `ipv6`. Changing
    this creates a new server.

  + `address` - (Required, ForceNew) IP address or domain name of the web server that the client accesses. For
    example, `192.168.1.1` or `www.example.com`. Changing this creates a new server.

  + `port` - (Required, ForceNew) Port number used by the web server. The value ranges from 0 to 65535. Changing this
    creates a new server.

* `certificate_id` - (Optional) Specifies the certificate ID. This parameter is mandatory when `client_protocol`
  is set to HTTPS.

* `policy_id` - (Optional) Specifies the policy ID associated with the domain. If not specified, a new policy
  will be created automatically.

* `proxy` - (Optional) Specifies whether a proxy is configured. Default value is `false`.

  -> **NOTE:** WAF forwards only HTTP/S traffic. So WAF cannot serve your non-HTTP/S traffic, such as UDP, SMTP, FTP,
  and basically all other non-HTTP/S traffic. If a proxy such as public network ELB (or Nginx) has been used, set
  proxy `true` to ensure that the WAF security policy takes effect for the real source IP address.

* `keep_policy` - (Optional) Specifies whether to retain the policy when deleting a domain name.
  Defaults to `true`.

* `protect_status` - (Optional) The protection status of domain, `0`: suspended, `1`: enabled.
  Default value is `1`.

* `tls` - (Optional) Specifies the minimum required TLS version.
  Values are:
  + `TLS v1.0`
  + `TLS v1.1`
  + `TLS v1.2`
  + `TLS v1.3`

* `cipher` - (Optional) Specifies the cipher suite of domain.
  Values are:
  + `cipher_1` - ECDHE-ECDSA-AES256-GCM-SHA384:HIGH:!MEDIUM:!LOW:!aNULL:!eNULL:!DES:!MD5:!PSK:!RC4:!kRSA:!SRP:!3DES:!DSS:!EXP:!CAMELLIA:@STRENGTH
  + `cipher_2` - EECDH+AESGCM:EDH+AESGCM
  + `cipher_3` - ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-SHA384:RC4:HIGH:!MD5:!aNULL:!eNULL:!NULL:!DH:!EDH
  + `cipher_4` - ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-SHA384:AES256-SHA256:RC4:HIGH:!MD5:!aNULL:!eNULL:!NULL:!EDH
  + `cipher_default` - ECDHE-RSA-AES256-SHA384:AES256-SHA256:RC4:HIGH:!MD5:!aNULL:!eNULL:!NULL:!DH:!EDH:!AESGCM

* `pci_3ds` - (Optional) Specifies the status of the PCI 3DS compliance certification check.
  Values are: `true` and `false`. This parameter must be used together with tls and cipher.

  -> **NOTE:** Tls must be set to `TLS v1.2`, and cipher must be set to `cipher_2`. The PCI 3DS compliance certification
  check cannot be disabled after being enabled.

* `pci_dss` - (Optional) Specifies the status of the PCI DSS compliance certification check.
  Values are: `true` and `false`. This parameter must be used together with tls and cipher.

  -> **NOTE:** Tls must be set to `TLS v1.2`, and cipher must be set to `cipher_2`.


## Attributes Reference

The following attributes are exported:

* `id` - ID of the domain.

* `certificate_name` - The name of the certificate used by the domain name.

* `access_status` - Whether a domain name is connected to WAF. Valid values are:
  + `0` - The domain name is not connected to WAF,
  + `1` - The domain name is connected to WAF.

* `protocol` - The protocol type of the client. The options are `HTTP` and `HTTPS`.

* `compliance_certification` - The compliance certifications of the domain, values are:
  + `pci_dss` - The status of domain PCI DSS, `true`: enabled, `false`: disabled.
  + `pci_3ds` - The status of domain PCI 3DS, `true`: enabled, `false`: disabled.

* `alarm_page` - The alarm page of domain. Valid values are:
  + `template_name` - The template of alarm page, values are: `default`, `custom` and `redirection`.
  + `redirect_url` - The redirection URL when `template_name` is set to `redirection`.

* `traffic_identifier` - The traffic identifier of domain. Valid values are:
  + `ip_tag` - The IP tag of traffic identifier.
  + `session_tag` - The session tag of traffic identifier.
  + `user_tag` - The user tag of traffic identifier.

* `created_at` - Timestamp when the dedicated WAF domain was created.

## Import

WAF dedicated domain can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_waf_dedicated_domain_v1.dom <id>
```
