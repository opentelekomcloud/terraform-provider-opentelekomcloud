package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataAzName = "data.opentelekomcloud_dcs_az_v1.az1"

func TestAccDcsAZV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsAZV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID(dataAzName),
					resource.TestCheckResourceAttr(dataAzName, "port", "8002"),
				),
			},
			{
				Config: testAccDcsAZV1DataSourceByName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDcsAZV1DataSourceID(dataAzName),
					resource.TestCheckResourceAttr(dataAzName, "port", "8002"),
					resource.TestCheckResourceAttr(dataAzName, "code", env.OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func testAccCheckDcsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find DCS AZ data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dcs AZ data source ID not set")
		}

		return nil
	}
}

var testAccDcsAZV1DataSourceBasic = fmt.Sprintf(`
data "opentelekomcloud_dcs_az_v1" "az1" {
  code = "%s"
  port = "8002"
}
`, env.OS_AVAILABILITY_ZONE)

var testAccDcsAZV1DataSourceByName = fmt.Sprintf(`
data "opentelekomcloud_dcs_az_v1" "az1" {
  name = "%s"
}
`, env.OS_AVAILABILITY_ZONE)
