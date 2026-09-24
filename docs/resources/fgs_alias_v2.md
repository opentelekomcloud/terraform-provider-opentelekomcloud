---
subcategory: "FunctionGraph"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_fgs_alias_v2"
sidebar_current: "docs-opentelekomcloud-resource-fgs-alias-v2"
description: |-
  Manages an FGS Alias resource within T-Cloud Public (formerly OpenTelekomCloud).
---

Up-to-date reference of API arguments for FGS aliases you can get at
[documentation portal](https://docs.otc.t-systems.com/function-graph/api-ref/api/versions_and_aliases/index.html)

# opentelekomcloud_fgs_alias_v2

Manages an alias of a function version within T-Cloud Public (formerly OpenTelekomCloud).

## Example Usage

```hcl
variable "function_urn" {}
variable "version" {}

resource "opentelekomcloud_fgs_alias_v2" "test" {
  function_urn = var.function_urn
  name         = "test"
  version      = var.version
  description  = "test alias"
}
```

## Argument Reference

The following arguments are supported:

* `function_urn` - (Required, String, ForceNew) Specifies the URN of the function.
  Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the alias. Can contain letters, digits, hyphens (-), and underscores (_). Must start with a letter, and end with a letter or digit. Length: `1` to `64`.
  Changing this will create a new resource.

* `version` - (Required, String) Specifies the version corresponding to the alias.

* `description` - (Optional, String) Specifies the description of the alias.

* `additional_version_weights` - (Optional, Map[String]Integer) Specifies the weights of additional versions for traffic distribution. The key is the version and the value is the weight.

* `additional_version_strategy` - (Optional, List) Specifies the strategies of additional versions for traffic distribution. The [object](#additional_version_strategy) structure is documented below.

<a name="additional_version_strategy"></a>
The `additional_version_strategy` block supports the following parameters:

* `version` - (Required, String) Specifies the version to which the strategy applies.

* `combine_type` - (Optional, String) Specifies the rule aggregation mode. `and`: All rules are met. `or`: Any rule is met.

* `rules` - (Optional, List) Specifies the rule of the strategy.
  The [object](#additional_version_strategy_rules) structure is documented below.

<a name="additional_version_strategy_rules"></a>
The `rules` block supports the following parameters:

* `rule_type` - (Required, String) Specifies the rule type of the strategy. Accepted values: `Header`.

* `param` - (Optional, String) Specifies the rule parameter name, which can contain only letters, digits, underscores (_), and hyphens (-).

* `op` - (Optional, String) Specifies the rule matching operator. Accepted values: `=`, `in`.

* `value` - (Optional, String) Specifies the rule value. If op is set to `in`, the value is a multi-value character string separated by commas (,).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in the format of `function_urn/alias_name`.

* `alias_urn` - Specifies the URN of the alias.

* `last_modified` - Specifies the last modified time of the alias.

## Import

The configurations can be imported using their related `function_urn/alias_name`, e.g.

```bash
$ terraform import opentelekomcloud_fgs_alias_v2.test <function_urn>/<alias_name>
```
