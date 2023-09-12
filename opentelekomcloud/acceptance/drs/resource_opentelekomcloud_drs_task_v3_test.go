package drs

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/drs/v3/public"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/common/cfg"
)

const taskV3ResourceName = "opentelekomcloud_drs_task_v3.test"

func TestAccDrsTaskV3MigrationBasic(t *testing.T) {
	var job public.QueryJobResp
	drsTaskName := tools.RandomString("drs_task-", 5)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDrsTaskV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDrsTaskV3Basic(drsTaskName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDrsTaskV3Exists(taskV3ResourceName, &job),
					resource.TestCheckResourceAttr(taskV3ResourceName, "name", drsTaskName),
					resource.TestCheckResourceAttr(taskV3ResourceName, "status", "FULL_TRANSFER_STARTED"),
					resource.TestCheckResourceAttr(taskV3ResourceName, "description", "TEST"),
					resource.TestCheckResourceAttr(taskV3ResourceName, "type", "migration"),
					resource.TestCheckResourceAttr(taskV3ResourceName, "direction", "down"),
				),
			},
		},
	})
}

func testAccCheckDrsTaskV3Exists(n string, jobQuery *public.QueryJobResp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		config := common.TestAccProvider.Meta().(*cfg.Config)
		client, err := config.DrsV3Client(env.OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating DRS v3 client, error: %s", err)
		}

		detailResp, err := public.BatchListTaskDetails(client, public.BatchQueryTaskOpts{Jobs: []string{rs.Primary.ID}})
		if err != nil {
			return err
		}
		found := detailResp.Results[0]

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("DRS Task not found")
		}

		// instant deletion of DRS task results in an error
		time.Sleep(20 * time.Second)
		jobQuery = &found

		return nil
	}
}

func testAccCheckDrsTaskV3Destroy(s *terraform.State) error {
	config := common.TestAccProvider.Meta().(*cfg.Config)
	client, err := config.DrsV3Client(env.OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating OpenTelekomCloud DRSv3 client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_drs_task_v3" {
			continue
		}

		_, err := public.BatchListTaskDetails(client, public.BatchQueryTaskOpts{Jobs: []string{rs.Primary.ID}})
		if err == nil {
			return fmt.Errorf("DRSv3 task still exists")
		}
	}

	return nil
}

func testAccDrsTaskV3Basic(drsName string) string {
	return fmt.Sprintf(`
%s
%s
%s

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "%s"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "down"
  net_type       = "eip"
  migration_type = "FULL_TRANS"
  description    = "TEST"
  force_destroy  = "true"

  source_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_1.address
    port        = "3306"
    user        = "root"
    password    = "MySql_120521"
    instance_id = opentelekomcloud_rds_instance_v3.mysql_1.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  }

  destination_db {
    engine_type = "mysql"
    ip          = opentelekomcloud_networking_floatingip_v2.fip_2.address
    port        = 3306
    user        = "root"
    password    = "MySql_120521"
    instance_id = opentelekomcloud_rds_instance_v3.mysql_2.id
    subnet_id   = data.opentelekomcloud_vpc_subnet_v1.shared_subnet.id
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet, RdsPreset, drsName)
}
