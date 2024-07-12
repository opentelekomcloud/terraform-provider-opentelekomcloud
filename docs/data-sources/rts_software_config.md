---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_software_config_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rts-software-config-v1"
description: |-
  Get details about a specific RTS Software Config from OpenTelekomCloud
---

# opentelekomcloud_rts_software_config_v1

Use this data source to get details about a specific RTS Software Config.

## Example Usage


```hcl
variable "config_name" {}

variable "server_id" {}

data "opentelekomcloud_rts_software_config_v1" "myconfig" {
  id = var.config_name
}

resource "opentelekomcloud_rts_software_deployment_v1" "mydeployment" {
  config_id = data.opentelekomcloud_rts_software_config_v1.myconfig.id
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the software configuration.

* `name` - (Optional) The name of the software configuration.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group` - The namespace that groups this software configuration by when it is delivered to a server.

* `inputs` -  A list of software configuration inputs.

* `outputs` - A list of software configuration outputs.

* `config` - The software configuration code.

* `options` - The software configuration options.
