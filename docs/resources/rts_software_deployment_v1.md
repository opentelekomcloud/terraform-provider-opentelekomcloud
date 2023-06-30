---
subcategory: "Resource Template Service (RTS)"
---

Up-to-date reference of API arguments for RTS deployment you can get at
`https://docs.otc.t-systems.com/resource-template-service/api-ref/apis/software_configuration_management`.

# opentelekomcloud_rts_software_deployment_v1

Provides an RTS software deployment resource within OpenTelekomCloud.

## Example Usage

```hcl
variable "config_id" {}
variable "server_id" {}

resource "opentelekomcloud_rts_software_deployment_v1" "mydeployment" {
  config_id = var.config_id
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `config_id` - (Required) The id of the software configuration resource running on an instance.

* `server_id` - (Required) The id of the instance.

* `status` - (Optional) The current status of deployment resources.

* `action` - (Optional) The stack action that triggers this deployment resource.

* `input_values` - (Optional) The input data stored in the form of a key-value pair.

* `output_values` - (Optional) The output data stored in the form of a key-value pair.

* `status_reason` - (Optional) The cause of the current deployment resource status.

* `tenant_id` - (Optional) The id of the authenticated tenant who can perform operations on the deployment resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The id of the software deployment.

## Import

Software deployment can be imported using the `deployment id`, e.g.

```sh
terraform import opentelekomcloud_rts_software_deployment_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
