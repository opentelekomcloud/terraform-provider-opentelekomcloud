package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// PASS
func TestAccOTCBMSNicV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBMSNic(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenTelekomCloudBMSNicV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSNicV2DataSourceID("data.opentelekomcloud_compute_bms_nic_v2.nic_1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_compute_bms_nic_v2.nic_1", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckBMSNicV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nic data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("nic data source ID not set ")
		}

		return nil
	}
}

var testAccOpenTelekomCloudBMSNicV2DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_nic_v2" "nic_1" {
  server_id = "%s"
  id = "%s"
}
`, OS_SERVER_ID, OS_NIC_ID)
