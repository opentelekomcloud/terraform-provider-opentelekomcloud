---
subcategory: "Simple Message Notification (SMN)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_smn_message_template_v2"
sidebar_current: "docs-opentelekomcloud-resource-smn-message-template-v2"
description: |-
  Manages an SMN Message Template resource within OpenTelekomCloud.
---

# opentelekomcloud_smn_message_template_v2

Manages a SMN message template resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "name" {}
variable "protocol" {}

resource "opentelekomcloud_smn_message_template_v2" "test" {
  name     = var.name
  protocol = var.protocol
  content  = "this content contains {content1}, {content2}, {content3}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the message template name.
  A template name starts with a letter or digit, consists of 1 to 64 characters,
  and can contain only letters, digits,  hyphens (-), and underscores (_).

  Changing this parameter will create a new resource.

* `protocol` - (Required, String, ForceNew) Specifies the protocol supported by the template. Value options:
  + **default**: the default protocol
  + **email**: the email protocol
  + **sms**: the SMS protocol
  + **http**: the http protocol
  + **https**: the https protocol

  Changing this parameter will create a new resource.

* `content` - (Required, String) Specifies the template content, which supports plain text only.
  The template content cannot be left blank or larger than 256 KB.
  The fields within "{}" can be replaced based on the actual situation
  when you use the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tag_names` - Indicates the variable list. The variable name will be quoted in braces ({}) in the template.
  When you use a template to send messages, you can replace the variable with any content.


* `region` - The resource region.

## Import

The SMN message template can be imported using the `id`, e.g.

```bash
$ terraform import opentelekomcloud_smn_message_template_v2.test <message_template_id>
```
