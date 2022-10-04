---
subcategory: "Elastic Cloud Server (ECS)"
---

# opentelekomcloud_compute_instance_v2

Get information on an ECS instance.

## Example Usage

```hcl
data "opentelekomcloud_compute_instance_v2" "instance" {
  # Search ecs by name
  name = "server_1"
}

data "opentelekomcloud_compute_instance_v2" "instance" {
  # Randomly generated UUID, for demonstration purposes
  id = "2ba26dc6-a12d-4889-8f25-794ea5bf4453"
}
```

## Argument Reference

* `id` - (Required) The UUID of the instance


## Attributes Reference

In addition to the above, the following attributes are exported:

* `name` - The name of the server.

* `image_id` - The image ID used to create the server.

* `image_name` - The image name used to create the server.

* `flavor_id` - The flavor ID used to create the server.

* `flavor_name` - The flavor name used to create the server.

* `user_data` - The user data added when the server was created.

* `security_groups` - An array of security group names associated with this server.

* `availability_zone` - The availability zone of this server.

* `network` - An array of maps, detailed below.

* `access_ip_v4` - The first IPv4 address assigned to this server.

* `access_ip_v6` - The first IPv6 address assigned to this server.

* `key_pair` - The name of the key pair assigned to this server.

* `tags` - A set of string tags assigned to this server.

* `metadata` - A set of key/value pairs made available to the server.

* `ssh_private_key_path` - (Optional) The path to the private key to use for SSH access. Required only if you want to
  get the password from the windows instance.


The `network` block is defined as:

* `uuid` - The UUID of the network

* `name` - The name of the network

* `fixed_ip_v4` - The IPv4 address assigned to this network port.

* `fixed_ip_v6` - The IPv6 address assigned to this network port.

* `port` - The port UUID for this network

* `mac` - The MAC address assigned to this network interface.
