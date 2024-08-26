package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccResourceRDSV3Maintenance_basic(t *testing.T) {
	resourceName := "opentelekomcloud_rds_maintenance_v3.test"
	postfix := acctest.RandString(3)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsBackupV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRDSV3MaintenanceBasic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "start_time", "12:00"),
					resource.TestCheckResourceAttr(resourceName, "end_time", "16:00"),
				),
			},
		},
	})
}

func testAccResourceRDSV3MaintenanceBasic(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "10"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "COMMON"
    size = 40
  }
  flavor = "rds.pg.c2.large"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}

resource "opentelekomcloud_rds_maintenance_v3" "test" {
  instance_id = opentelekomcloud_rds_instance_v3.instance.id
  start_time  = "12:00"
  end_time    = "16:00"
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
