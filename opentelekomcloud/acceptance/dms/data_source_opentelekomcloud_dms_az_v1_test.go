package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

const dataAzName = "data.opentelekomcloud_dms_az_v1.az1"

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { common.TestAccPreCheck(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID(dataAzName),
					resource.TestCheckResourceAttr(dataAzName, "name", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(dataAzName, "port", "8002"),
					resource.TestCheckResourceAttr(dataAzName, "code", env.OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSourceBasic = fmt.Sprintf(`
data "opentelekomcloud_dms_az_v1" "az1" {
  name = "%s"
  port = "8002"
}
`, env.OS_AVAILABILITY_ZONE)
