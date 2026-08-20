---
subcategory: "Log Tank Service (LTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_lts_structuring_custom_templates_v3"
sidebar_current: "docs-opentelekomcloud-datasource-lts-structuring-custom-templates-v3"
description: |-
  Queries LTS custom structuring templates.
---

Up-to-date reference of API arguments for LTS custom structuring templates is available in the
[documentation portal](https://docs.otc.t-systems.com/log-tank-service/api-ref/apis/cloud_structuring/querying_a_structuring_template.html).

# opentelekomcloud_lts_structuring_custom_templates_v3

Use this data source to query custom LTS structuring templates.

## Example Usage

```hcl
data "opentelekomcloud_lts_structuring_custom_templates_v3" "all" {}
```

To query a specific template:

```hcl
variable "template_id" {}

data "opentelekomcloud_lts_structuring_custom_templates_v3" "selected" {
  id = var.template_id
}
```

## Argument Reference

* `id` - (Optional, String) Specifies the template ID. When omitted, all custom structuring templates in the project
  are returned.

## Attribute Reference

* `templates` - List of custom structuring templates matching the filter.

  The [templates](#templates_struct) structure is documented below.

* `region` - Region in which the templates were queried.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - Template ID.
* `project_id` - Project ID.
* `template_name` - Template name.
* `template_type` - Structuring type.
* `demo_log` - Sample log event.
* `demo_fields` - Structured sample fields.
* `tag_fields` - Structured tag fields.
* `rule` - Structuring rule.
* `demo_label` - Sample log label.
* `created_at` - Template creation time in RFC3339 format.

The `demo_fields` block supports:

* `field_name` - Field name.
* `content` - Example field value.
* `type` - Field data type.
* `is_analysis` - Whether quick analysis is enabled.
* `index` - Field sequence number.
* `relation` - Hierarchical relationship between fields.
* `user_defined_name` - Custom field alias.

The `tag_fields` block supports:

* `field_name` - Field name.
* `content` - Example field value.
* `type` - Field data type.
* `is_analysis` - Whether quick analysis is enabled.
* `index` - Field sequence number.

The `rule` block supports:

* `type` - Structuring type.
* `param` - Type-specific structuring rule serialized as JSON.
