---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_stack_v1"
sidebar_current: "docs-opentelekomcloud-resource-rts-stack-v1"
description: |-
  Manages an RTS Stack resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RTS stack you can get at
[documentation portal](https://docs.otc.t-systems.com/resource-template-service/api-ref/apis/stack_management)

# opentelekomcloud_rts_stack_v1

Provides an OpenTelekomCloud Stack.

## Example Usage

```hcl
variable "name" {}
variable "network_id" {}
variable "instance_type" {}
variable "image_id" {}

resource "opentelekomcloud_rts_stack_v1" "mystack" {
  name             = var.name
  disable_rollback = true
  timeout_mins     = 60
  parameters = {
    "network_id"    = var.network_id
    "instance_type" = var.instance_type
    "image_id"      = var.image_id
  }
  template_body = <<JSON
  {
    "heat_template_version": "2016-04-08",
    "description": "Simple template to deploy",
    "parameters": {
        "image_id": {
            "type": "string",
            "description": "Image to be used for compute instance",
            "label": "Image ID"
        },
        "network_id": {
            "type": "string",
            "description": "The Network to be used",
            "label": "Network UUID"
        },
        "instance_type": {
            "type": "string",
            "description": "Type of instance (Flavor) to be used",
            "label": "Instance Type"
        }
    },
    "resources": {
        "my_instance": {
            "type": "OS::Nova::Server",
            "properties": {
                "image": {
                    "get_param": "image_id"
                },
                "flavor": {
                    "get_param": "instance_type"
                },
                "networks": [{
                    "network": {
                        "get_param": "network_id"
                    }
                }]
            }
        }
    },
    "outputs": {
      "InstanceIP": {
        "description": "Instance IP",
        "value": { "get_attr": ["my_instance", "first_address"] }
      }
    }
  }
JSON
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the stack. The value must meet the regular expression rule (`^[a-zA-Z][a-zA-Z0-9_.-]{0,254}$`). Changing this creates a new stack.

* `template_body` - (Optional; Required if `template_url` is empty) Structure containing the template body. The template content must use the yaml syntax.

* `template_url` - (Optional; Required if `template_body` is empty) Location of a file containing the template body.

* `environment` - (Optional) Tthe environment information about the stack.

* `files` - (Optional) Files used in the environment.

* `parameters` - (Optional) A list of Parameter structures that specify input parameters for the stack.

* `disable_rollback` - (Optional) Set to true to disable rollback of the stack if stack creation failed.

* `timeout_mins` - (Optional) Specifies the timeout duration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `outputs` - A map of outputs from the stack.

* `capabilities` - List of stack capabilities for stack.

* `notification_topics` - List of notification topics for stack.

* `status` - Specifies the stack status.


## Import

RTS Stacks can be imported using the `name`, e.g.

```sh
terraform import opentelekomcloud_rts_stack_v1.mystack rts-stack
```
