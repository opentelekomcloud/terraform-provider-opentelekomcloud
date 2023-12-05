# Create a example for resources

This example should match to next pattern:
* Create folder for example, like `rds_instance_v3`

  **NOTE:** Please use only underscores in folder names
* Add this folder name in `runtime.yaml`, this file helps zuul job `otc-terraform-visualize-main` to process this example and run necessary checks.


`runtime.yaml` content:
```yaml
---
folders:
- rds_instance_v3
- elb_v3
- autoscaling_with_elb_v3
- autoscaling_with_alarm
- nginx_app_on_compute_instance_v2
- dns_zone_v2
---
```

* Create `main.tf` and `settings.tf` inside created folder

  `settings.tf` should contain main provider setup like:
```hcl
terraform {
  required_providers {
    opentelekomcloud = {
      source  = "opentelekomcloud/opentelekomcloud"
      version = ">= 1.35.9"
    }
  }
}

provider "opentelekomcloud" {
  cloud = "functest_cloud"
}
```
  and `main.tf` should contain the whole example

    **NOTE:** you can create separate file for variables and outputs.

`variables.tf`:
```hcl
variable "az" {
  default = "eu-de-01"
}
```

`outputs.tf`:
```hcl
output "db_id" {
  value = opentelekomcloud_rds_instance_v3.instance.id
}
```

`main.tf`
```hcl
module "network" {
  source = "../modules/network"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_1"
  availability_zone = [var.az]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = module.network.default_security_group_id
  subnet_id         = module.network.shared_subnet.network_id
  vpc_id            = module.network.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 0
  }
  tags = {
    muh = "value-create"
    kuh = "value-create"
  }
  lower_case_table_names = "0"
}
```

 * Now all is done

## Example Usage

```shell
terraform init
terraform plan
terraform apply
terraform destroy
```

## Requirements

| Name             | Version   |
|------------------|-----------|
| terraform        | >= 1.6.3  |
| opentelekomcloud | >= 1.35.9 |
