---
subcategory: "GeminiDB"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_gemini_template_v3"
sidebar_current: "docs-opentelekomcloud-resource-gemini-template-v3"
description: |-
  Manages a GeminiDB Parameter Template resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for GeminiDB you can get at
[documentation portal](https://docs.otc.t-systems.com/geminidb/api-ref/apis_v3/parameter_templates/index.html).

# opentelekomcloud_gemini_template_v3

Manages a GeminiDB Parameter Template resource within OpenTelekomCloud.

## Example Usage

### Basic Parameter Template
```hcl
resource "opentelekomcloud_gemini_template_v3" "template" {
  name           = "cassandra_template"
  description    = "Custom Cassandra configuration"
  instance_type  = "cassandra"
  engine_version = "3.11"

  parameters {
    name  = "write_request_timeout_in_ms"
    value = "7000"
  }

  parameters {
    name  = "read_request_timeout_in_ms"
    value = "8000"
  }

  parameters {
    name  = "slow_query_log_timeout_in_ms"
    value = "15000"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the parameter template name. The template name can include
  1 to 64 characters and can contain only uppercase letters, lowercase letters, digits, hyphens (-),
  underscores (_), and periods (.). Changing this creates a new resource.

* `instance_type` - (Required, String, ForceNew) Specifies the database type. Currently, only `cassandra`
  is supported. Changing this creates a new resource.

* `engine_version` - (Required, String, ForceNew) Specifies the database version. Currently, only `3.11`
  is supported for GeminiDB Cassandra. Changing this creates a new resource.

* `parameters` - (Required, Set) Specifies the parameter values. The structure is documented below.

* `description` - (Optional, String) Specifies the parameter template description. It can contain a
  maximum of 256 characters. The following special characters are not allowed: `>!<"&'=`

The `parameters` block supports:

* `name` - (Required, String) Specifies the parameter name.

* `value` - (Required, String) Specifies the parameter value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (parameter template ID).

* `region` - The region in which the parameter template is created.

* `created_at` - The creation time in the `yyyy-MM-ddTHH:mm:ssZ` format.

* `updated_at` - The update time in the `yyyy-MM-ddTHH:mm:ssZ` format.

* `parameters` - The parameter list. In addition to the arguments above, the following attributes are exported:
    * `need_restart` - Indicates whether the instance needs to be restarted after the parameter is modified.
    * `readonly` - Indicates whether the parameter is read-only.
    * `value_range` - The value range of the parameter.
    * `data_type` - The data type of the parameter.
    * `description` - The parameter description.
