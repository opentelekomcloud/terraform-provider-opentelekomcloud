package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccDmsAZV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDms(t) },
		ProviderFactories: common.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsAZV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDmsAZV1DataSourceID("data.opentelekomcloud_dms_az_v1.az1"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "name", env.OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "port", "8002"),
					resource.TestCheckResourceAttr(
						"data.opentelekomcloud_dms_az_v1.az1", "code", env.OS_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func testAccCheckDmsAZV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find Dms az data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("dms az data source ID not set")
		}

		return nil
	}
}

var testAccDmsAZV1DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_dms_az_v1" "az1" {
  name = "%s"
  port = "8002"
}
`, env.OS_AVAILABILITY_ZONE)
