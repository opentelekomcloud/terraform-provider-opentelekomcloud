package acceptance

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/common"
	"github.com/opentelekomcloud/terraform-provider-opentelekomcloud/opentelekomcloud/acceptance/env"
)

func TestAccOTCBMSServerV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccBmsFlavorPreCheck(t) },
		Providers: common.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOTCBMSServerV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBMSServerV2DataSourceID("data.opentelekomcloud_compute_bms_server_v2.server1"),
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
resource "opentelekomcloud_compute_instance_v2" "instance_1" {
  name              = "BMSinstance_1"
  image_id          = "%s"
  security_groups   = ["default"]
  availability_zone = "%s"
  flavor_id         = "physical.o2.medium"
  flavor_name       = "physical.o2.medium"
  metadata          = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}

data "opentelekomcloud_compute_bms_server_v2" "server1" {
  id = opentelekomcloud_compute_instance_v2.instance_1.id
}
`, env.OS_IMAGE_ID, env.OS_AVAILABILITY_ZONE, env.OS_NETWORK_ID)
