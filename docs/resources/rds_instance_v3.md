---
subcategory: "Relational Database Service (RDS)"
---

# opentelekomcloud_rds_instance_v3

Manages RDS instance v3 resource.

## Example Usage

### Create a single db instance

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "terraform_test_rds_instance"
  availability_zone = [var.availability_zone]

  db {
    password = "P@ssw0rd1!9851"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }

  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = var.subnet_id
  vpc_id            = var.vpc_id
  flavor            = "rds.pg.c2.medium"

  volume {
    type = "COMMON"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a primary/standby db instance

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "terraform_test_rds_instance"
  availability_zone = [var.availability_zone_1, var.availability_zone_2]

  db {
    password = "P@ssw0rd1!9851"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }
  security_group_id   = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id           = var.subnet_id
  vpc_id              = var.vpc_id
  flavor              = "rds.pg.s1.medium.ha"
  ha_replication_mode = "async"

  volume {
    type = "COMMON"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a db instance with public IP

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_compute_floatingip_v2" "ip" {}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  availability_zone = [
    var.availability_zone_1,
    var.availability_zone_2
  ]
  db {
    password = "Telekom!120521"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }
  name              = "terraform_test_rds_instance"
  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = var.subnet_id
  vpc_id            = var.vpc_id
  volume {
    type = "COMMON"
    size = 100
  }
  flavor              = "rds.pg.s1.medium.ha"
  ha_replication_mode = "async"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
  public_ips = [
    opentelekomcloud_compute_floatingip_v2.ip.address
  ]
  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Create a single db instance with encrypted volume

```hcl
resource "opentelekomcloud_kms_key_v1" "key" {
  key_alias       = "key_1"
  key_description = "first test key"
  is_enabled      = true
}

resource "opentelekomcloud_networking_secgroup_v2" "secgroup" {
  name        = "terraform_test_security_group"
  description = "terraform security group acceptance test"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "terraform_test_rds_instance"
  availability_zone = [var.availability_zone]

  security_group_id = opentelekomcloud_networking_secgroup_v2.secgroup.id
  subnet_id         = var.subnet_id
  vpc_id            = var.vpc_id
  flavor            = "rds.pg.c2.medium"

  db {
    password = "P@ssw0rd1!9851"
    type     = "PostgreSQL"
    version  = "9.5"
    port     = "8635"
  }

  volume {
    disk_encryption_id = opentelekomcloud_kms_key_v1.key.id
    type               = "COMMON"
    size               = 100
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

### Overriding parameters from template

```hcl
resource "opentelekomcloud_networking_secgroup_v2" "sg" {
  name = "sg-rds-test"
}

resource "opentelekomcloud_rds_parametergroup_v3" "pg" {
  name = "pg-rds-test"
  values = {
    autocommit = "OFF"
  }
  datastore {
    type    = "postgresql"
    version = "10"
  }
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = [var.availability_zone]

  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }

  security_group_id = opentelekomcloud_networking_secgroup_v2.sg.id
  subnet_id         = var.subnet_id
  vpc_id            = var.vpc_id
  flavor            = "rds.pg.c2.medium"
  volume {
    type = "COMMON"
    size = 40
  }

  parameters = {
    max_connections = "37",
  }
}
```

### Restore backup to a new instance

```hcl
data "opentelekomcloud_rds_backup_v3" "backup" {
  instance_id = var.rds_instance_id
  type        = "auto"
}

resource "opentelekomcloud_rds_instance_v3" "from_backup" {
  name              = "instance-restored"
  availability_zone = opentelekomcloud_rds_instance_v3.instance.availability_zone
  flavor            = "rds.pg.c2.medium"

  restore_point {
    instance_id = data.opentelekomcloud_rds_backup_v3.backup.instance_id
    backup_id   = data.opentelekomcloud_rds_backup_v3.backup.id
  }

  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = var.security_group_id
  subnet_id         = var.os_network_id
  vpc_id            = var.os_router_id
  volume {
    type = "COMMON"
    size = 40
  }
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) Specifies the AZ name. Changing this parameter will create a new resource.

* `db` - (Required) Specifies the database information. Structure is documented below. Changing this parameter will create a new resource.

* `flavor` - (Required) Specifies the specification code. Use data source [opentelekomcloud_rds_flavors_v3](../data-sources/rds_flavors_v3.md) to get a list of available flavor names. Examples could be `rds.pg.c2.medium` or `rds.pg.c2.medium.ha` for HA clusters.

* `name` - (Required) Specifies the DB instance name. The DB instance name of the same type
  must be unique for the same tenant. The value must be 4 to 64
  characters in length and start with a letter. It is case-sensitive
  and can contain only letters, digits, hyphens (-), and underscores
  (_).  Changing this parameter will create a new resource.

* `security_group_id` - (Required) Specifies the security group which the RDS DB instance belongs to.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required) Specifies the subnet id. Changing this parameter will create a new resource.

* `volume` - (Required) Specifies the volume information. Structure is documented below.

* `vpc_id` - (Required) Specifies the VPC ID. Changing this parameter will create a new resource.

* `backup_strategy` - (Optional) Specifies the advanced backup policy. Structure is documented below.

* `ha_replication_mode` - (Optional) Specifies the replication mode for the standby DB instance. For MySQL, the value
  is async or semisync. For PostgreSQL, the value is async or sync. For Microsoft SQL Server, the value is sync.

-> Async indicates the asynchronous replication mode. `semisync` indicates the
  semi-synchronous replication mode. sync indicates the synchronous
  replication mode.  Changing this parameter will create a new resource.

* `param_group_id` - (Optional) Specifies the parameter group ID.

* `parameters` - (Optional) Map of additional configuration parameters. Values should be strings. Parameters set here
  overrides values from configuration template (parameter group).

* `public_ips` - (Optional) Specifies floating IP to be assigned to the instance.
  This should be a list with single element only.

-> Setting public IP is done with assigning floating IP to internally
  created port. So RDS itself doesn't know about this assignment. This assignment
  won't show on the console.
  This argument will be ignored in future when RDSv3 API for EIP assignment will be implemented.

* `tag` - (Optional) Tags key/value pairs to associate with the instance. Deprecated, please use
  the `tags` instead.

* `tags` - (Optional) Tags key/value pairs to associate with the instance.

* `restore_point` - (Optional) Specifies the restoration information.

The `db` block supports:

* `password` - (Required) Specifies the database password. The value cannot be
  empty and should contain 8 to 32 characters, including uppercase
  and lowercase letters, digits, and the following special
  characters: ~!@#%^*-_=+? You are advised to enter a strong
  password to improve security, preventing security risks such as
  brute force cracking.  Changing this parameter will create a new resource.

* `port` - (Optional) Specifies the database port information. The MySQL database port
  ranges from 1024 to 65535 (excluding 12017 and 33071, which are
  occupied by the RDS system and cannot be used). The PostgreSQL
  database port ranges from 2100 to 9500. The Microsoft SQL Server
  database port can be 1433 or ranges from 2100 to 9500, excluding
  5355 and 5985. If this parameter is not set, the default value is
  as follows: For MySQL, the default value is 3306. For PostgreSQL,
  the default value is 5432. For Microsoft SQL Server, the default
  value is 1433.  Changing this parameter will create a new resource.

* `type` - (Required) Specifies the DB engine. Value: MySQL, PostgreSQL, SQLServer. Changing this parameter will create a new resource.

* `version` - (Required) Specifies the database version. MySQL databases support MySQL 5.6
  and 5.7. PostgreSQL databases support PostgreSQL 9.5 and 9.6. Microsoft SQL Server
  databases support 2014 SE, 2016 SE, and 2016 EE.
  Changing this parameter will create a new resource.

The `volume` block supports:

* `disk_encryption_id` - (Optional) Specifies the key ID for disk encryption. Changing this parameter will create a new resource.

* `size` - (Required) Specifies the volume size. Its value range is from 40 GB to 4000
  GB. The value must be a multiple of 10. Changing this resize the volume.

* `type` - (Required) Specifies the volume type. Its value can be any of the following
  and is case-sensitive: COMMON: indicates the SATA type.
  ULTRAHIGH: indicates the SSD type.  Changing this parameter will create a new resource.

The `backup_strategy` block supports:

* `keep_days` - (Optional) Specifies the retention days for specific backup files. The value
  range is from 0 to 732. If this parameter is not specified or set
  to 0, the automated backup policy is disabled. NOTICE:
  Primary/standby DB instances of Microsoft SQL Server do not
  support disabling the automated backup policy.

* `start_time` - (Required) Specifies the backup time window. Automated backups will be
  triggered during the backup time window. It must be a valid value in the &quot;hh:mm-HH:MM&quot;
  format. The current time is in the UTC format. The HH value must
  be 1 greater than the hh value. The values of mm and MM must be
  the same and must be set to any of the following: 00, 15, 30, or
  45. Example value: 08:15-09:15 23:00-00:00.

The `restore_point` block supports:

* `instance_id` - (Required) Specifies the original DB instance ID.

* `backup_id` - (Optional) Specifies the ID of the backup used to restore data.

* `restore_time` - (Optional) Specifies the time point of data restoration in the UNIX timestamp.
  The unit is millisecond and the time zone is UTC.

-> Exactly one of `backup_id` and `restore_time` needs to be set.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` - Indicates the creation time.

* `nodes` - Indicates the instance nodes information. Structure is documented below.

* `private_ips` - Indicates the private IP address list. It is a blank string until an
  ECS is created.

* `public_ips` - Indicates the public IP address list.

* `db` - See Argument Reference above. The `db` block additionally contains:

  * `user_name` - Indicates the default user name of database.

The `nodes` block contains:

* `availability_zone` - Indicates the AZ.

* `id` - Indicates the node ID.

* `name` - Indicates the node name.

* `role` - Indicates the node type. The value can be master or slave, indicating the primary node or standby node respectively.

* `status` - Indicates the node status.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 30 minute.

## Import

RDS instance can be imported using the `id`, e.g.

```sh
terraform import opentelekomcloud_rds_instance_v3.instance_1 7117d38e-4c8f-4624-a505-bd96b97d024c
```

## Notes

But due to some attributes missing from the API response, it's required to ignore changes as below.

```hcl
resource "opentelekomcloud_rds_instance_v3" "instance_1" {
  # ...

  lifecycle {
    ignore_changes = [
      "db",
    ]
  }
}
```
