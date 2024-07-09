---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_software_config_v1"
sidebar_current: "docs-opentelekomcloud-resource-rts-software-config-v1"
description: |-
Manages an RTS Config resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for RTS config you can get at
[documentation portal](https://docs.otc.t-systems.com/resource-template-service/api-ref/apis/software_configuration_management)

# opentelekomcloud_rts_software_config_v1

Provides an RTS software config resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "config_name" {}

resource "opentelekomcloud_rts_software_config_v1" "myconfig" {
  name = var.config_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the software configuration.

* `group` - (Optional) The namespace that groups this software configuration by when it is delivered to a server.

* `input_values` - (Optional) A list of software configuration inputs.

* `output_values` - (Optional) A list of software configuration outputs.

* `config` - (Optional) The software configuration code.

* `options` - (Optional) The software configuration options.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the software config.

## Import

Software Config can be imported using the `config id`, e.g.

```sh
terraform import opentelekomcloud_rts_software_config_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
