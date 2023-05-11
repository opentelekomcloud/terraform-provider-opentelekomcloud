---
subcategory: "Web Application Firewall (WAF)"
---

Up-to-date reference of API arguments for WAF certificate you can get at
`https://docs.otc.t-systems.com/web-application-firewall/api-ref/apis/certificate_management`.

# opentelekomcloud_waf_certificate_v1

Manages a WAF certificate resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_certificate_v1" "certificate_1" {
  name    = "cert_1"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The certificate name. The maximum length is 256 characters.
  Only digits, letters, underscores(`_`), and hyphens(`-`) are allowed.

* `content` - (Optional) The certificate content. Changing this creates a new certificate.

* `key` - (Optional) The private key. Changing this creates a new certificate.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the certificate.

* `name` - See Argument Reference above.

* `content` - See Argument Reference above.

* `key` - See Argument Reference above.

## Import

Certificates can be imported using the `id`, e.g.

```shell
terraform import opentelekomcloud_waf_certificate_v1.cert_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```

Or using `name`, e.g.

```shell
terraform import opentelekomcloud_waf_certificate_v1.cert_1 cert_1
```
