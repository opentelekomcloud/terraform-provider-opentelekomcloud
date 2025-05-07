---
subcategory: "Cloud Firewall (CFW)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_cfw_domain_name_group_v1"
sidebar_current: "docs-opentelekomcloud-resource-cfw-domain-name-group-v1"
description: |-
  Manages a CFW domain name group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for CFW domain name group you can get at
[documentation portal](https://docs.otc.t-systems.com/cloud-firewall/api-ref/api/domain_name_resolution_and_domain_name_group_management/index.html)

# opentelekomcloud_cfw_domain_name_group_v1

Manages a CFW domain name Group resource within OpenTelekomCloud.

## Example Usage:
```hcl
variable firewall_id {}
variable object_id {}

resource "opentelekomcloud_cfw_domain_name_group_v1" "group_1" {
  firewall_id = var.firewall_id
  object_id   = var.object_id
  name        = "test-acc-tf-domain-group"
  domain_names {
    domain_name = "www.testacctf.com"
  }
}
```

## Argument Reference

The following arguments are supported:
* `firewall_id` - (Required, String, ForceNew) Specifies the Firewall ID.

* `object_id` - (Required, String, ForceNew) Specifies the protected object ID, which is used to distinguish between Internet border protection and VPC border protection after a cloud firewall is created. If the value of type is 0, the protected object ID belongs to the Internet border. If the value of type is 1, the protected object ID belongs to the VPC border.

* `name` - (Required, String) Specifies the CFW domain name group name. The CFW domain name group name of the same type is unique in the same firewall.

* `description` - (Optional, String) Specifies the description of the domain name group.

* `domain_names` - (Required, List, ForceNew) Specifies the domain name information list. The [domain_names](#domainnames) structure is documented below.

* `domain_set_type` - (Optional, Integer, ForceNew) Specifies the domain name group typ: `0` (application domain name group), `1` (network domain name group).

<a name="domainnames"></a>
The `domain_names` block supports:

* `domain_name` - (Required, String, ForceNew) Specifies the domain name, for example, www.test.com.
* `description` - (Optional, String, ForceNew) Specifies the domain name description.

## Attributes Reference

In addition to the arguments mentioned above, the following attributes are exported:

* `id` - Indicates the domain name group ID.
* `domain_names` - Indicates the domain name information list. The `domain_name` structure is as follows:
    - `domain_name` - See Argument Reference above.
    - `description` - See Argument Reference above.
    - `domain_address_id` - Indicates the domain name ID.
* `ref_count` - Indicates the number of times a domain name group is referenced by rules.
* `config_status` - Indicates the configuration status: `-1` (unconfigured), `0` (configuration failed), `1` (configuration succeeded), `2` (configuring), 3 (normal), or `4` (abnormal).
* `rules` - Indicates the used rule ID list. The `rules` structure is as follows:
    - `id` - Indicates the rule ID.
    - `name` - Indicates the rule name.

## Import

CFW Domain name Group V1 resource can be imported using the firewall ID, `firewall_id`, the object ID, `object_id`, and the name of the group, `name`, e.g.

```shell
terraform import opentelekomcloud_cfw_domain_name_group_v1.group_1 <firewall_id>/<object_id>/<name>
```
