---
subcategory: "Direct Connect (DCaaS)"
---
# opentelekomcloud_dc_endpoint_group_v2 (Resource)

Up-to-date reference of API arguments for Direct Connect Endpoint Group you can get at
[Official Docs Portal](https://docs.otc.t-systems.com/direct-connect/api-ref/apis/direct_connect_endpoint_group/index.html).

~>
opentelekomcloud_dc_endpoint_group_v2 is no longer provided. Impossible to update assigned endpoint group.
Please use `opentelekomcloud_dc_virtual_gateway_v2` with `local_ep_group` block instead.

## Example Usage

```hcl
data "opentelekomcloud_identity_project_v3" "project_1" {
  name = "eu-de_project_1"
}

resource "opentelekomcloud_dc_endpoint_group_v2" "dc_endpoint_group" {
  name        = "ep-1"
  type        = "cidr"
  endpoints   = ["10.2.0.0/24", "10.3.0.0/24"]
  description = "first"
  project_id  = data.opentelekomcloud_identity_project_v3.project_1.id
}
```

## Argument Reference

The following arguments are supported:

* `name` (String, Optional, ForceNew) - Specifies the name of the Direct Connect endpoint group.
* `project_id` (String, Required, ForceNew) - Specifies the project ID.
* `description` (String, Optional, ForceNew) - Provides supplementary information about the Direct Connect endpoint group.
* `endpoints` (List, required, ForceNew) - Specifies the list of the endpoints in a Direct Connect endpoint group.
* `type` (String, Required, ForceNew) - Specifies the type of the Direct Connect endpoints. The value can only be cidr.

## Attributes Reference

The following attributes are exported:

* `id` -  ID of the endpoint group.

## Import

Direct Connect Endpoint Group can be imported using `id`, e.g.

```sh
$ terraform import opentelekomcloud_dc_endpoint_group_v2.eg <id>
```
