---
subcategory: "Virtual Private Network (VPN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_vpnaas_ike_policy_v2"
sidebar_current: "docs-opentelekomcloud-resource-vpnaas-ike-policy-v2"
description: |-
Manages a VPNAAS IKE Policy resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for VPNAAS ike policy you can get at
[documentation portal](https://docs.otc.t-systems.com/virtual-private-network/api-ref/native_openstack_apis/ike_policy_management)

# opentelekomcloud_vpnaas_ike_policy_v2

Manages a V2 IKE policy resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpnaas_ike_policy_v2" "policy_1" {
  name = "my_policy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create a VPN service. If omitted, the
  `region` argument of the provider is used. Changing this creates a new service.

* `name` - (Optional) The name of the policy.

* `tenant_id` - (Optional) The owner of the policy. Required if admin wants to
  create a service for another policy. Changing this creates a new policy.

* `description` - (Optional) The human-readable description for the policy.

* `auth_algorithm` - (Optional) The authentication hash algorithm. Valid values are `md5`,
  `sha1`, `sha2-256`, `sha2-384`, `sha2-512`. Default is `sha1`.

* `encryption_algorithm` - (Optional) The encryption algorithm. Valid values are `3des`, `aes-128`, `aes-192` and so on.
  The default value is `aes-128`.

* `pfs` - (Optional) The perfect forward secrecy mode. Valid values are `group1`, `group2`, `group5` and so on.
  Default is `group5`.

* `phase1_negotiation_mode` - (Optional) The IKE mode. Valid values are `main` and `aggressive`. Default is `main`.

* `ike_version` - (Optional) The IKE mode. Valid values are `v1` and `v2`. Default is `v1`.

* `lifetime` - (Optional) The lifetime of the security association. Consists of Unit and Value.
  * `unit` - (Optional) The units for the lifetime of the security association. A valid value is `seconds`. Default is `seconds`.
  * `value` - (Optional) The value for the lifetime of the security association. Default is `3600`.

* `value_specs` - (Optional) Map of additional options.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.

* `name` - See Argument Reference above.

* `tenant_id` - See Argument Reference above.

* `description` - See Argument Reference above.

* `auth_algorithm` - See Argument Reference above.

* `encapsulation_mode` - See Argument Reference above.

* `encryption_algorithm` - See Argument Reference above.

* `pfs` - See Argument Reference above.

* `transform_protocol` - See Argument Reference above.

* `lifetime` - See Argument Reference above.
  * `unit` - See Argument Reference above.
  * `value` - See Argument Reference above.

* `value_specs` - See Argument Reference above.

## Import

Services can be imported using the `id`, e.g.

```
terraform import opentelekomcloud_vpnaas_ike_policy_v2.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
