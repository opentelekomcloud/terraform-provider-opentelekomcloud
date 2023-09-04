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

variable "param_group_id" {
	type    = string
	default = "qwert"
}

variable "password" {
	type    = string
	default = "12345"
}

locals {
	rg = csvdecode(file("mysql.csv"))
}

resource "opentelekomcloud_rds_instance_v3" "mysql" {
	for_each            = { for rg in local.rg : rg.name => rg }
	name                = each.value.name
	flavor              = each.value.flavor
	ha_replication_mode = each.value.ha
	vpc_id              = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
	subnet_id           = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
	security_group_id   = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.secgroup_id
	param_group_id      = var.param_group_id
	availability_zone   = [each.value.az]

	db {
		type     = "MySQL"
		version  = "5.7"
		password = var.password
	}

	volume {
		type = "CLOUDSSD"
		size = each.value.size
	}

	backup_strategy {
		start_time = "08:00-09:00"
		keep_days  = 7
	}
}

resource "opentelekomcloud_drs_task_v3" "test" {
	for_each       = { for rg in local.rg : rg.name => rg }
		name           = each.value.name1
		type           = "migration"
		engine_type    = "mysql"
		direction      = "down"
		net_type       = "vpc"
		migration_type = "FULL_TRANS"
		description    = "TEST"
		force_destroy  = "true"

	source_db {
		engine_type = "mysql"
		ip          = each.value.sip
		port        = "3306"
		user        = each.value.suser
		password    = each.value.spass
	}

	destination_db {
		region      = "cn-east-3"
		ip          = opentelekomcloud_rds_instance_v3.mysql[each.value.name].fixed_ip
		port        = 3306
		engine_type = "mysql"
		user        = "root"
		password    = var.password
		instance_id = opentelekomcloud_rds_instance_v3.mysql[each.value.name].id
		subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
	}
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)
}
