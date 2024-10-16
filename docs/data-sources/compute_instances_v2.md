---
subcategory: "Elastic Cloud Server (ECS)"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_compute_instances_v2"
sidebar_current: "docs-opentelekomcloud-datasource-compute-instances-v2"
description: |-
  Get ECS instances details from OpenTelekomCloud
---

Up-to-date reference of API arguments for ECS instances you can get at
[documentation portal](https://docs.otc.t-systems.com/elastic-cloud-server/api-ref/native_openstack_nova_apis/lifecycle_management/querying_ecss.html#en-us-topic-0020212688)

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
  status = "ACTIVE"
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

* `description` - Server description.

* `instances` - List of ECS instance details. The object structure of each ECS instance is documented below.

The `instances` block supports:

* `id` - The instance ID in UUID format.

* `name` - The instance name.

* `image_id` - The image ID used to create the server.

* `image_name` - The image name used to create the server.

* `flavor_id` - The flavor ID used to create the server.

* `status` - The instance status.

* `key_pair` - The key pair that is used to authenticate the instance.

* `security_groups_ids` - An array of one or more security group Ids to associate with the instance.

* `availability_zone` - The availability zone of this server.

* `project_id` - The instance project ID.

* `network` - An array of maps, detailed below.

The `network` block is defined as:

* `uuid` - The UUID of the network

* `name` - The name of the network

* `fixed_ip_v4` - The IPv4 address assigned to this network port. Not supported.

* `fixed_ip_v6` - The IPv6 address assigned to this network port. Not supported.

* `port` - The port UUID for this network

* `mac` - The MAC address assigned to this network interface.
