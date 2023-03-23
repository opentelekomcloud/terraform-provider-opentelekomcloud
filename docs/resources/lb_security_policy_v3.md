---
subcategory: "Dedicated Load Balancer (DLB)"
---

# opentelekomcloud_lb_security_policy_v3

Manages a Dedicated Load Balancer Security Policy resource within OpenTelekomCloud.

## Example Usage Basic

```hcl
resource "opentelekomcloud_lb_security_policy_v3" "policy_1" {
  name        = "elb-security-policy"
  description = "This is security policy"
  protocols   = ["TLSv1", "TLSv1.1"]
  ciphers     = ["ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA"]
}
```

## Example Usage Security policy assigned to ELB listener

```hcl
resource "opentelekomcloud_vpc_v1" "vpc_v1" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "opentelekomcloud_vpc_subnet_v1" "subnet_v1" {
  name   = var.subnet_name
  cidr   = var.subnet_cidr
  vpc_id = opentelekomcloud_vpc_v1.vpc_v1.id

  gateway_ip    = var.subnet_gateway_ip
  ntp_addresses = "10.100.0.33,10.100.0.34"
}

resource "opentelekomcloud_lb_loadbalancer_v3" "loadbalancer_1" {
  name        = "loadbalancer_1"
  router_id   = opentelekomcloud_vpc_subnet_v1.subnet_v1.vpc_id
  network_ids = [opentelekomcloud_vpc_subnet_v1.subnet_v1.network_id]

  availability_zones = [var.az]
}

resource "opentelekomcloud_lb_certificate_v3" "certificate_1" {
  name        = "certificate_1"
  type        = "server"
  private_key = var.private_key
  certificate = var.certificate
}

resource "opentelekomcloud_lb_listener_v3" "listener_1" {
  name                      = "listener_1"
  description               = "some interesting description"
  loadbalancer_id           = opentelekomcloud_lb_loadbalancer_v3.loadbalancer_1.id
  protocol                  = "HTTPS"
  protocol_port             = 443
  default_tls_container_ref = opentelekomcloud_lb_certificate_v3.certificate_1.id
  security_policy_id        = opentelekomcloud_lb_security_policy_v3.policy_1.id

  advanced_forwarding = true
  sni_match_algo      = "wildcard"

  insert_headers {
    forwarded_host = true
  }
}
resource "opentelekomcloud_lb_security_policy_v3" "policy_1" {
  name      = "security-policy"
  protocols = ["TLSv1", "TLSv1.1"]
  ciphers   = ["ECDHE-ECDSA-AES128-SHA", "ECDHE-RSA-AES128-SHA"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Specifies the security policy name.

* `description` - (Optional) Provides supplementary information about the security policy.

* `protocols` - (Required) Lists the TLS protocols supported by the custom security policy.
* Possible values: `TLSv1`, `TLSv1.1`, `TLSv1.2`, and `TLSv1.3`.

* `ciphers` - (Required) Lists the cipher suites supported by the custom security policy.
* The protocol and cipher suite must match. At least one cipher suite must match the protocol.
* Possible values:
  `ECDHE-RSA-AES256-GCM-SHA384`, `ECDHE-RSA-AES128-GCM-SHA256`,`ECDHE-ECDSA-AES256-GCM-SHA384`,
  `ECDHE-ECDSA-AES128-GCM-SHA256`, `AES128-GCM-SHA256`, `AES256-GCM-SHA384`, `ECDHE-ECDSA-AES128-SHA256`,
  `ECDHE-RSA-AES128-SHA256`, `AES128-SHA256,AES256-SHA256`, `ECDHE-ECDSA-AES256-SHA384`, `ECDHE-RSA-AES256-SHA384`,
  `ECDHE-ECDSA-AES128-SHA`, `ECDHE-RSA-AES128-SHA`, `ECDHE-RSA-AES256-SHA`, `ECDHE-ECDSA-AES256-SHA`,
  `AES128-SHA`, `AES256-SHA`, `CAMELLIA128-SHA`, `DES-CBC3-SHA`, `CAMELLIA256-SHA`, `ECDHE-RSA-CHACHA20-POLY1305`,
  `ECDHE-ECDSA-CHACHA20-POLY1305`, `TLS_AES_128_GCM_SHA256`, `TLS_AES_256_GCM_SHA384`, `TLS_CHACHA20_POLY1305_SHA256`,
  `TLS_AES_128_CCM_SHA256`, `TLS_AES_128_CCM_8_SHA256`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the policy.

* `project_id` - The project ID of the custom security policy.

* `listeners` - The listeners that use the custom security policies.

* `created_at` - The time when the custom security policy was created.

* `updated_at` - The time when the custom security policy was updated.

## Import

Load Balancer Policy can be imported using the Policy ID, e.g.:

```shell
terraform import opentelekomcloud_lb_security_policy_v3.this 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
