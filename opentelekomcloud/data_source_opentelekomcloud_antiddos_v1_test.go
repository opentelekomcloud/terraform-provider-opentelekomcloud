package opentelekomcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAntiDdosV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiDdosV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAntiDdosV1DataSourceID("data.opentelekomcloud_antiddos_v1.antiddos"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_antiddos_v1.antiddos", "network_type", "EIP"),
					resource.TestCheckResourceAttr("data.opentelekomcloud_antiddos_v1.antiddos", "status", "normal"),
				),
			},
		},
	})
}

func testAccCheckAntiDdosV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find defense status of EIP data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Defense status of EIP data source ID not set")
		}

		return nil
	}
}

const testAccAntiDdosV1DataSource_basic = `
resource "opentelekomcloud_vpc_eip_v1" "eip_1" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name = "test"
    size = 8
    share_type = "PER"
    charge_mode = "traffic"
  }
}

resource "opentelekomcloud_antiddos_v1" "antiddos_1" {
  floating_ip_id = "${opentelekomcloud_vpc_eip_v1.eip_1.id}"
  enable_l7 = true
  traffic_pos_id = 1
  http_request_pos_id = 2
  cleaning_access_pos_id = 1
  app_type_id = 0
}

data "opentelekomcloud_antiddos_v1" "antiddos" {  
  floating_ip_id = "${opentelekomcloud_antiddos_v1.antiddos_1.id}"
}
`
