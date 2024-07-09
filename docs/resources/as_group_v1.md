---
subcategory: "Autoscaling"
layout: "opentelekomcloud"
page_title: "OpenTelekomCloud: opentelekomcloud_as_group_v1"
sidebar_current: "docs-opentelekomcloud-resource-as-group-v1"
description: |-
Manages a AS Group resource within OpenTelekomCloud.
---

Up-to-date reference of API arguments for AS group you can get at
[documentation portal](https://docs.otc.t-systems.com/auto-scaling/api-ref/apis/as_groups)

# opentelekomcloud_as_group_v1

Manages a V1 Autoscaling Group resource within OpenTelekomCloud.

## Example Usage

### Basic Autoscaling Group

```hcl
resource "opentelekomcloud_as_group_v1" "as_group" {
  scaling_group_name       = "as_group"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }

  security_groups {
    id = "45e4c6de-6bf0-4843-8953-2babde3d4810"
  }

  vpc_id           = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip  = true
  delete_instances = "yes"

  tags = {
    muh = "kuh"
  }
}
```

### Autoscaling Group Only Remove Members When Scaling Down

```hcl
resource "opentelekomcloud_as_group_v1" "as_group_only_remove_members" {
  scaling_group_name       = "as_group_only_remove_members"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
  security_groups {
    id = "45e4c6de-6bf0-4843-8953-2babde3d4810"
  }

  vpc_id           = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip  = true
  delete_instances = "no"
}
```

### Autoscaling Group With ELB Listener

```hcl
resource "opentelekomcloud_as_group_v1" "as_group_with_elb" {
  scaling_group_name       = "as_group_with_elb"
  scaling_configuration_id = "37e310f5-db9d-446e-9135-c625f9c2bbfc"
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 10

  networks {
    id = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  }
  security_groups {
    id = "45e4c6de-6bf0-4843-8953-2babde3d4810"
  }

  vpc_id           = "1d8f7e7c-fe04-4cf5-85ac-08b478c290e9"
  delete_publicip  = true
  delete_instances = "yes"

  lbaas_listeners {
    pool_id       = opentelekomcloud_lb_pool_v2.pool_1.id
    protocol_port = opentelekomcloud_lb_listener_v2.as_listener.protocol_port
  }
}

resource "opentelekomcloud_lb_listener_v2" "as_listener" {
  name            = "as_listener"
  description     = "as test listener"
  protocol        = "TCP"
  protocol_port   = 80
  loadbalancer_id = "cba48790-baf5-4446-adb3-02069a916e97"
}

resource "opentelekomcloud_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = opentelekomcloud_lb_listener_v2.as_listener.id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_name` - (Required) The name of the scaling group. The name can contain letters,
  digits, underscores(_), and hyphens(-),and cannot exceed 64 characters.

* `scaling_configuration_id` - (Optional) The configuration ID which defines
  configurations of instances in the AS group.

* `desire_instance_number` - (Optional) The expected number of instances. The default
  value is the minimum number of instances. The value ranges from the minimum number of
  instances to the maximum number of instances.

* `min_instance_number` - (Optional) The minimum number of instances.
  The default value is 0.

* `max_instance_number` - (Optional) The maximum number of instances.
  The default value is 0.

* `cool_down_time` - (Optional) The cooling duration (in seconds). The value ranges
  from 0 to 86400, and is 900 by default.

* `lb_listener_id` **DEPRECATED** - (Optional) The Classic LB listener IDs. The system
  supports up to six Classic LB listeners, the IDs of which are separated using a comma (,).
  This parameter is alternative to `lbaas_listeners`.

* `lbaas_listeners` - (Optional) An array of one or more Enhanced Load Balancer.
  The system supports the binding of up to six Enhanced Load Balancers. The field is
  alternative to `lb_listener_id`. The `lbaas_listeners` object structure is
  documented below.

* `available_zones` - (Optional) Specifies the AZ information. The ECS
  associated with a scaling action will be created in a specified AZ.
  If you do not specify an AZ, the system automatically specifies one.

* `networks` - (Required) An array of one or more network IDs.
  The system supports up to five networks. The networks object structure
  is documented below.

* `security_groups` - (Optional) An array of security group IDs to associate with the group.
  A maximum of one security group can be selected. The `security_groups` object structure is
  documented below.

* `vpc_id` - (Required) The VPC ID. Changing this creates a new group.

* `health_periodic_audit_method` - (Optional) The health check method for instances
  in the AS group. The health check methods include `ELB_AUDIT` and `NOVA_AUDIT`.
  If load balancing is configured, the default value of this parameter is `ELB_AUDIT`.
  Otherwise, the default value is `NOVA_AUDIT`.

* `health_periodic_audit_time` - (Optional) The health check period for instances.
  The value can be 1, 5, 15, 60, or 180 in the unit of minutes. If this parameter
  is not specified, the default value is 5. If the value is set to 0, health check
  is performed every 10 seconds.

* `health_periodic_audit_grace_period` - (Optional) The grace period for instance health check.
  The unit is second and value range is 0-86400. The default value is 600. The health check grace
  period starts after an instance is added to an AS group and is enabled. The AS group will start
  checking the instance status only after the grace period ends. This parameter is valid only when
  the instance health check method of the AS group is ELB_AUDIT.

* `instance_terminate_policy` - (Optional) The instance removal policy. The policy has
  four options: `OLD_CONFIG_OLD_INSTANCE` (default), `OLD_CONFIG_NEW_INSTANCE`,
  `OLD_INSTANCE`, and `NEW_INSTANCE`.

* `notifications` - (Optional) The notification mode. The system only supports `EMAIL`
  mode which refers to notification by email.

* `delete_publicip` - (Required) Whether to delete the elastic IP address bound to the
  instances of AS group when deleting the instances. The options are `true` and `false`.

* `delete_instances` - (Required) Whether to delete the instances in the AS group
  when deleting the AS group. The options are `yes` and `no`.

The `networks` block supports:

* `id` - (Required) The network UUID.

The `security_groups` block supports:

* `id` - (Required) The UUID of the security group.

The `lbaas_listeners` block supports:

* `pool_id` - (Required) Specifies the backend ECS group ID.

* `protocol_port` - (Required) Specifies the backend protocol, which is the port on which
  a backend ECS listens for traffic. The number of the port ranges from 1 to 65535.

* `weight` - (Optional) Specifies the weight, which determines the portion of requests a
  backend ECS processes compared to other backend ECSs added to the same listener. The value
  of this parameter ranges from 0 to 100. The default value is 1.

* `tags` - (Optional) Tags key/value pairs to associate with the AutoScaling Group.

## Attributes Reference

The following attributes are exported:

* `scaling_group_name` - See Argument Reference above.

* `status` - Indicates the status of the AS group.

* `current_instance_number` - Indicates the number of current instances in the AS group.

* `desire_instance_number` - See Argument Reference above.

* `min_instance_number` - See Argument Reference above.

* `max_instance_number` - See Argument Reference above.

* `cool_down_time` - See Argument Reference above.

* `lb_listener_id` - See Argument Reference above.

* `health_periodic_audit_method` - See Argument Reference above.

* `health_periodic_audit_time` - See Argument Reference above.

* `instance_terminate_policy` - See Argument Reference above.

* `scaling_configuration_id` - See Argument Reference above.

* `delete_publicip` - See Argument Reference above.

* `notifications` - See Argument Reference above.

* `instances` - The instances IDs of the AS group.

* `tags` - See Argument Reference above.
