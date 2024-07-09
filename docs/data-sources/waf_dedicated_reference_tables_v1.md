---
subcategory: "Dedicated Web Application Firewall (WAFD)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_waf_dedicated_reference_tables_v1"
sidebar_current: "docs-opentelekomcloud-datasource-waf-dedicated-reference-tables-v1"
description: |-
Get a list of WAF reference tables from OpenTelekomCloud
---

Up-to-date reference of API arguments for WAF reference table you can get at
[documentation portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/querying_the_reference_table_list.html)

# opentelekomcloud_waf_dedicated_reference_tables_v1

Use this data source to get a list of OpenTelekomCloud WAF reference tables.

## Example Usage

```hcl
data "opentelekomcloud_waf_dedicated_reference_tables_v1" "table" {
  name = "reference_table_1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to create the WAF reference table resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The name of the reference table. The value is case-sensitive and matches exactly.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `tables` - A list of WAF reference tables.

The `tables` block supports:

* `id` - The id of the reference table.

* `name` - The name of the reference table. The maximum length is 64 characters.

* `type` - The type of the reference table, The options are: `url`, `user-agent`, `ip`, `params`, `cookie`, `referer`
  and `header`.

* `conditions` - The conditions of the reference table.

* `description` - The description of the reference table.

* `created_at` - The time when reference table was created.
