package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccRdsMajorVersionsDataSource_basic(t *testing.T) {
	dataSourceName := "data.opentelekomcloud_rds_major_versions_v3.test"
	postfix := acctest.RandString(3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckRdsInstanceV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsMajorVersionsDataSource_basic(postfix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "available_versions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "instance_id"),
				),
			},
		},
	})
}

func testAccRdsMajorVersionsDataSource_basic(postfix string) string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_rds_instance_v3" "instance" {
  name              = "tf_rds_instance_%s"
  availability_zone = ["%s"]
  db {
    password = "Postgres!120521"
    type     = "PostgreSQL"
    version  = "13"
    port     = "8635"
  }
  security_group_id = data.opentelekomcloud_networking_secgroup_v2.default_secgroup.id
  subnet_id         = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.network_id
  vpc_id            = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.vpc_id
  volume {
    type = "CLOUDSSD"
    size = 100
  }
  flavor = "rds.pg.n1.large.4"
  backup_strategy {
    start_time = "08:00-09:00"
    keep_days  = 1
  }
}

data "opentelekomcloud_rds_major_versions_v3" "test" {
  depends_on = [
    opentelekomcloud_rds_instance_v3.instance,
  ]
  instance_id = opentelekomcloud_rds_instance_v3.instance.id
}
`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, postfix, env.OS_AVAILABILITY_ZONE)
}
