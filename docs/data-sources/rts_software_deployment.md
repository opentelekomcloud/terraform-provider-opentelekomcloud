---
subcategory: "Resource Template Service (RTS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_rts_software_deployment_v1"
sidebar_current: "docs-opentelekomcloud-datasource-rts-software-deployment-v1"
description: |-
  Get details about a specific RTS Software Deployment from OpenTelekomCloud
---

# opentelekomcloud_rts_software_deployment_v1

Use this data source to get details about a specific RTS Software Deployment.

## Example Usage

```hcl
variable "deployment_id" {}

data "opentelekomcloud_rts_software_deployment_v1" "mydeployment" {
  id = var.deployment_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The id of the software deployment.

* `config_id` - (Optional) The id of the software configuration resource running on an instance.

* `server_id` - (Optional) The id of the instance.

* `status` - (Optional) The current status of deployment resources.

* `action` - (Optional)  The stack action that triggers this deployment resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `input_values` - The input data stored in the form of a key-value pair.

* `output_values` - The output data stored in the form of a key-value pair.

* `status_reason` - The cause of the current deployment resource status.

