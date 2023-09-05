package drs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const taskV3ResourceName = "opentelekomcloud_drs_task_v3.test"

func TestAccDrsTaskV3Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsTaskV3Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(taskV3ResourceName, "name", "test"),
				),
			},
		},
	})
}

func testAccDrsTaskV3Basic() string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_networking_floatingip_v2" "fip_1" {}
resource "opentelekomcloud_networking_floatingip_v2" "fip_2" {}

resource "opentelekomcloud_rds_instance_v3" "mysql_1" {
  name                = "RDS-ha"
  flavor              = "rds.mysql.s1.large.ha"
  ha_replication_mode = "semisync"
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone   = ["eu-de-01,eu-de-03"]

  db {
    type     = "MySQL"
    version  = "5.7"
    password = "MySql!120521"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 7
  }

	public_ips = [opentelekomcloud_networking_floatingip_v2.fip_1.address]
}

resource "opentelekomcloud_rds_instance_v3" "mysql_2" {
  name                = "RDS-single"
  flavor              = "rds.mysql.m1.large"
  vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  availability_zone   = ["eu-de-01"]

  db {
    type     = "MySQL"
    version  = "5.7"
    password = "MySql!120521"
  }

  volume {
    type = "ULTRAHIGH"
    size = 100
  }

  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 7
  }

	public_ips = [opentelekomcloud_networking_floatingip_v2.fip_2.address]
}

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "drs-test"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "down"
  net_type       = "vpc"
  migration_type = "FULL_TRANS"
  description    = "TEST"
  force_destroy  = "true"

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_2.address
	instance_id = opentelekomcloud_rds_instance_v3.mysql_2.id
    port        = "3306"
    user        = "root"
    password    = "MySql!120521"
  }

  destination_db {
    region      = "eu-de-01"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "MySql!120521"
    instance_id = opentelekomcloud_rds_instance_v3.mysql_1.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)
}
