package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOTCBMSServerV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckBMSServer(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOTCBMSServerV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSServerV2DataSourceID("data.opentelekomcloud_compute_bms_server_v2.server1"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_compute_bms_server_v2.server1", "id", OS_SERVER_ID),
					resource.TestCheckResourceAttr("data.opentelekomcloud_compute_bms_server_v2.server1", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckBMSServerV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find servers data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("server data source ID not set ")
		}

		return nil
	}
}

var testAccOTCBMSServerV2DataSource_basic = fmt.Sprintf(`
data "opentelekomcloud_compute_bms_server_v2" "server1" {
  id = "%s"  
}
`, OS_SERVER_ID)
