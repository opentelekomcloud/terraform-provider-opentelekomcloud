---
subcategory: "Elastic Cloud Server (ECS)"
---

Up-to-date reference of API arguments for ECS server group management you can get at
`https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/openstack_nova_apis/ecs_group_management`.

# opentelekomcloud_compute_servergroup_v2

Manages a V2 Server Group resource within OpenTelekomCloud.

## Example Usage

```hcl
resource "opentelekomcloud_compute_servergroup_v2" "test-sg" {
  name     = "my-sg"
  policies = ["anti-affinity"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the server group. Changing this creates
  a new server group.

* `policies` - (Required) The set of policies for the server group. Only two
  two policies are available right now, and both are mutually exclusive. See
  the Policies section for more information. Changing this creates a new
  server group.

* `value_specs` - (Optional) Map of additional options.

## Policies

* `anti-affinity` - All instances/servers launched in this group will be
  hosted on different compute nodes.

## Attributes Reference

The following attributes are exported:

* `name` - See Argument Reference above.

* `policies` - See Argument Reference above.

* `members` - The instances that are part of this server group.

* `id` -  ID of the server group.

## Import

Server Groups can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_compute_servergroup_v2.test-sg 1bc30ee9-9d5b-4c30-bdd5-7f1e663f5edf
```
