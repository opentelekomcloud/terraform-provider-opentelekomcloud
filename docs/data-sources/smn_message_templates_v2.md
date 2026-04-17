---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_message_templates_v2"
sidebar_current: "docs-opentelekomcloud-datasource-smn-message-templates-v2"
description: |-
  Get details about an SMN Message templates resources within OpenTelekomCloud.
---

# opentelekomcloud_smn_message_templates_v2

Use this data source to get the list of SMN message templates.

## Example Usage

```hcl
variable "name" {}

data "opentelekomcloud_smn_message_templates_v2" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the message template.

* `protocol` - (Optional, String) Specifies the protocol of the message template.

* `template_id` - (Optional, String) Specifies the message template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of message templates.
  The [templates](#SmnMessageTemplate_MessageTemplate) structure is documented below.

<a name="SmnMessageTemplate_MessageTemplate"></a>
The `templates` block supports:

* `id` - Indicates the message template ID.

* `name` - Indicates the message template name.

* `protocol` - Indicates the protocol supported by the template.

* `tag_names` - Indicates the variable list. The variable name will be quoted in braces ({}) in the template.
  When you use a template to send messages, you can replace the variable with any content.

* `created_at` - Indicates the create time.

* `updated_at` - Indicates the update time.
