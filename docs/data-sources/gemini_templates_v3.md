---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_templates_v3"
sidebar_current: "docs-opentelekomcloud-datasource-gemini-templates-v3"
description: |-
  Get available GeminiDB parameter templates from OpenTelekomCloud
---

# opentelekomcloud_gemini_templates_v3

Use this data source to get the list of GeminiDB parameter templates, including all of the default and custom parameter templates.

## Example Usage
```hcl
data "opentelekomcloud_gemini_templates_v3" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `region` - The data source region.

* `templates` - Indicates the list of parameter templates.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - Indicates the ID of the parameter template.

* `name` - Indicates the name of the parameter template.

* `description` - Indicates the description of parameter template.

* `datastore_version_name` - Indicates the database version name.

* `datastore_name` - Indicates the database name.

* `created_at` - Indicates the creation time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `updated_at` - Indicates the update time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `mode` - Indicates the instance type.

* `user_defined` - Indicates whether the parameter template is a custom template.
    + **false**: default parameter template.
    + **true**: custom template.
