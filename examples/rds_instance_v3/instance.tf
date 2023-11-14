module "network" {
  source = "../modules/network"
}

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
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
