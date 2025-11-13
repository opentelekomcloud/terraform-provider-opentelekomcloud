---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_instance_template_v3"
sidebar_current: "docs-opentelekomcloud-datasource-gemini-instance-template-v3"
description: |-
  Get GeminiDB instance parameter settings from OpenTelekomCloud
---

# opentelekomcloud_gemini_instance_template_v3

Use this data source to query GeminiDB instance parameter settings.

## Example Usage
```hcl
variable "instance_id" {}

data "opentelekomcloud_gemini_instance_template_v3" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The parameter template ID.

* `region` - The parameter region.

* `datastore_version_name` - Indicates the database version name.

* `datastore_name` - Indicates the database name.

* `created_at` - Indicates the creation time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `updated_at` - Indicates the update time in the **yyyy-MM-ddTHH:mm:ssZ** format.

* `mode` - Indicates the instance type.

* `configuration_parameters` - Indicates the list of parameters defined by users based on a default parameter template.

  The [configuration_parameters](#configuration_parameters_struct) structure is documented below.

<a name="configuration_parameters_struct"></a>
The `configuration_parameters` block supports:

* `name` - Indicates the parameter name.

* `value` - Indicates the parameter value.

* `restart_required` - Indicates whether the instance needs to be restarted after the parameter is changed.
    + **false**: the instance does not need to be restarted.
    + **true**: the instance needs to be restarted.

* `readonly` - Indicates whether the parameter is read-only.
    + **false**: the parameter is not read-only.
    + **true**: the parameter is read-only.

* `value_range` - Indicates the value range. For example, the value of the Integer type ranges from 0 to 1,
  and the value of the Boolean type is true or false.

* `type` - Indicates the parameter type. The value can be **string**, **integer**, **boolean**, **list**, or **float**.

* `description` - Indicates the parameter description.
