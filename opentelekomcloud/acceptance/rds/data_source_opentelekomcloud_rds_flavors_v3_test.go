package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
)

const dataFlavorName = "data.opentelekomcloud_rds_flavors_v3.flavor"

func TestAccRdsFlavorV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsFlavorV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRdsFlavorV3DataSourceID(dataFlavorName),
				),
			},
		},
	})
}

func testAccCheckRdsFlavorV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find rds data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("rds data source ID not set ")
		}

		return nil
	}
}

const testAccRdsFlavorV3DataSourceBasic = `
data "opentelekomcloud_rds_flavors_v3" "flavor" {
  db_type       = "PostgreSQL"
  db_version    = "9.5"
  instance_mode = "rds.pg.s1.medium.ha"
}
`
