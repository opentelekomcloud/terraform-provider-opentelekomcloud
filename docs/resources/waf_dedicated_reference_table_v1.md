---
subcategory: "Dedicated Web Application Firewall (WAFD)"
---

Up-to-date reference of API arguments for WAF dedicated Precise Protection rule you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/web-application-firewall-dedicated/api-ref/apis/rule_management/creating_a_reference_table.html).


# opentelekomcloud_waf_dedicated_reference_table_v1

Manages a WAF Dedicated Reference Table resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_waf_dedicated_reference_table_v1" "table" {
  name = "%s"
  type = "url"

  conditions = [
    "/admin",
    "/manage"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the WAF reference table resource. If omitted,
  the provider-level region will be used. Changing this setting will push a new reference table.

* `name` - (Required, String) The name of the reference table. Only letters, digits, and underscores(_) are allowed. The
  maximum length is 64 characters.

* `type` - (Required, String, ForceNew) The type of the reference table, The options are `url`, `user-agent`, `ip`,
  `params`, `cookie`, `referer` and `header`. Changing this setting will push a new reference table.

* `conditions` - (Required, List) The conditions of the reference table. The maximum length is 30. The maximum length of
  condition is 2048 characters.

* `description` - (Optional, String) The description of the reference table. The maximum length is 128 characters.
  Currently, could be set only on update.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the reference table.

* `created_at` - The time when reference table was created.

## Import

Dedicated WAF Reference Table can be imported using `id`, e.g.

```bash
$ terraform import opentelekomcloud_waf_dedicated_reference_table_v1.table <id>
```
