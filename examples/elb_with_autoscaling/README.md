# ELB&AS Configuration

This example will show you a configuration which AS will be a backend of ELB.
This is a very common demo for this scenario. For more detailed information, please refer to
[doc](https://www.terraform.io/docs/providers/opentelekomcloud/index.html).

The ```main.tf``` contains the major resource scripts.

```hcl
resource "opentelekomcloud_elb_loadbalancer" "lb_example" {
  name           = "lb_example"
  type           = "External"
  description    = "This is an example configuration for LB"
  vpc_id         = var.vpc_id
  admin_state_up = "true"
  bandwidth      = 5
}

resource "opentelekomcloud_elb_listener" "listener_example" {
  name             = "listener_example"
  description      = "This is a listener example"
  protocol         = "TCP"
  backend_protocol = "TCP"
  protocol_port    = 12345
  backend_port     = 8080
  lb_algorithm     = "roundrobin"
  loadbalancer_id  = opentelekomcloud_elb_loadbalancer.lb_example.id
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_elb_healthcheck" "ht_example" {
  listener_id          = opentelekomcloud_elb_listener.listener_example.id
  healthcheck_protocol = "HTTP"
  healthy_threshold    = 5
  healthcheck_timeout  = 25
  healthcheck_interval = 3
  healthcheck_uri      = "/"
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}

resource "opentelekomcloud_as_group_v1" "group_example" {
  scaling_group_name       = "group_example"
  scaling_configuration_id = opentelekomcloud_as_configuration_v1.config_example.id
  desire_instance_number   = 2
  min_instance_number      = 0
  max_instance_number      = 3
  networks {
    id = var.subnet_id
  }
  security_groups {
    id = var.security_group_id
  }
  vpc_id           = var.vpc_id
  lb_listener_id   = opentelekomcloud_elb_listener.listener_example.id
  delete_publicip  = true
  delete_instances = "yes"
}

resource "opentelekomcloud_as_policy_v1" "policy_example" {
  scaling_policy_name = "policy_example"
  scaling_group_id    = opentelekomcloud_as_group_v1.group_example.id
  scaling_policy_type = "RECURRENCE"
  scaling_policy_action {
    operation = "ADD"
  }
  scheduled_policy {
    launch_time     = "07:00"
    recurrence_type = "Daily"
    end_time        = "2017-12-30T12:00Z"
  }
}

resource "opentelekomcloud_as_configuration_v1" "config_example" {
  scaling_configuration_name = "config_example"
  instance_config {
    flavor = var.flavor
    image  = var.image_id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }

    key_name  = var.keyname
    user_data = file("userdata.txt")
  }
}
```


Note: Before you run these scripts, please do not forget to replace the
<YOUR_XXX> with your actual values.
