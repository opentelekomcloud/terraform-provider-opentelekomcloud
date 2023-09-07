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
					resource.TestCheckResourceAttr(taskV3ResourceName, "name", "drs-test"),
				),
			},
		},
	})
}

func testAccDrsTaskV3Basic() string {
	return fmt.Sprintf(`
%s
%s

resource "opentelekomcloud_drs_task_v3" "test" {
  name           = "drs-test"
  type           = "migration"
  engine_type    = "mysql"
  direction      = "down"
  net_type       = "eip"
  migration_type = "FULL_TRANS"
  description    = "TEST"
  force_destroy  = "true"

  source_db {
    engine_type = "mysql"
    ip          = "80.158.44.40"
    port        = "3306"
    user        = "root"
    password    = "MySql_120521"
	instance_id = "7f4aad8f27384ac9aea2abf6a6ea2ef8in01"
    subnet_id   = "7aa19a7c-d0ea-4e02-a06b-30cfe2c8fe1f"
  }

  destination_db {
    region      = "eu-de"
    ip          = "80.158.62.80"
    port        = 3306
    engine_type = "mysql"
    user        = "root"
    password    = "MySql_120521"
    instance_id = "2c576c3c8864478bb43ade4a8d32cf84in01"
    subnet_id   = "7aa19a7c-d0ea-4e02-a06b-30cfe2c8fe1f"
  }
}`, common.DataSourceSecGroupDefault, common.DataSourceSubnet)
}
