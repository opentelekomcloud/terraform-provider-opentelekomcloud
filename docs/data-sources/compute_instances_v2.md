---
subcategory: "Elastic Cloud Server (ECS)"
---

# opentelekomcloud_compute_instances_v2

Get information on an ECS instances.

## Example Usage

```hcl
variable "name_regex" {}

data "opentelekomcloud_compute_instances_v2" "test" {
  name = var.name_regex
}
```

```hcl

data "opentelekomcloud_compute_instances_v2" "test" {
  status =  "ACTIVE"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the instance name, which can be queried with a regular expression.
  The instance name supports fuzzy matching query too.

* `instance_id` - (Optional, String) Specifies the ECS ID.

* `flavor_name` - (Optional, String) Specifies the flavor name of the instance.

* `status` - (Optional, String) Specifies the status of the instance. The valid values are as follows:
    + **ACTIVE**: The instance is running properly.
    + **SHUTOFF**: The instance has been properly stopped.
    + **ERROR**: An error has occurred on the instance.

* `image_id` - (Optional, String) Specifies the image ID of the instance.

* `flavor_id` - (Optional, String) Specifies the flavor ID.

* `key_pair` - (Optional, String) Specifies the key pair that is used to authenticate the instance.

* `project_id` - (Optional, String) Specifies the project where instance hosted.

* `limit` - (Optional, Integer) Specifies the number of instances to be queried. The value is an integer and is 100 by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - List of ECS instance details. The object structure of each ECS instance is documented below.

The `instances` block supports:

* `id` - The instance ID in UUID format.

* `name` - The instance name.

* `image_id` - The image ID of the instance.

* `flavor_id` - The flavor ID.

* `status` - The instance status.

* `key_pair` - The key pair that is used to authenticate the instance.

* `security_group_ids` - An array of one or more security group IDs to associate with the instance.

* `project_id` - The instance project ID.
