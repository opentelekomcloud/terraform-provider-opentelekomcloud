---
subcategory: "Virtual Private Network (VPN)"
---

Up-to-date reference of API arguments for VPNAAS ipsec policy service you can get at
`https://docs.otc.t-systems.com/virtual-private-network/api-ref/native_openstack_apis/ipsec_policy_management`.

# opentelekomcloud_vpnaas_ipsec_policy_v2

Manages a V2 IPSec policy resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_vpnaas_ipsec_policy_v2" "policy_1" {
  name = "my_policy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
  A Networking client is needed to create an IPSec policy. If omitted, the
  `region` argument of the provider is used. Changing this creates a new policy.

* `name` - (Optional) The name of the policy.

* `tenant_id` - (Optional) The owner of the policy. Required if admin wants to
  create a policy for another project. Changing this creates a new policy.

* `description` - (Optional) The human-readable description for the policy.

* `auth_algorithm` - (Optional) The authentication hash algorithm. Valid values are `md5`, `sha1`, `sha2-256`, `sha2-384`, `sha2-512`.
  Default is `sha1`.

* `encapsulation_mode` - (Optional) The encapsulation mode. Default is `tunnel`.

* `encryption_algorithm` - (Optional) The encryption algorithm. Valid values are `3des`, `aes-128`, `aes-192` and so on.
  The default value is `aes-128`.

* `pfs` - (Optional) The perfect forward secrecy mode. Valid values are `group1`, `group2`, `group5`, `group14`,
  `group15`, `group16`, `group19`, `group20`, `group21` or `disable` Default is `group5`.

* `transform_protocol` - (Optional) The transform protocol. Valid values are `esp`, `ah` and `ah-esp`. Default is `esp`.

* `lifetime` - (Optional) The lifetime of the security association. Consists of Unit and Value.
  - `unit` - (Optional) The units for the lifetime of the security association. Default is `seconds`.
  - `value` - (Optional) The value for the lifetime of the security association. Must be a positive integer. Default is `3600`.

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
  - `unit` - See Argument Reference above.
  - `value` - See Argument Reference above.

* `value_specs` - See Argument Reference above.


## Import

Policies can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_vpnaas_ipsec_policy_v2.policy_1 832cb7f3-59fe-40cf-8f64-8350ffc03272
```
